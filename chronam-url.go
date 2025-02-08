package headliner

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ChronamPage struct {
	URL      string   `json:"url"`
	RawText  string   `json:"raw_text"`
	Articles Articles `json:"articles"`
}

func (a *App) ChronamUrlHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var page ChronamPage

		if err := json.NewDecoder(r.Body).Decode(&page); err != nil {
			slog.Error("failed to decode request body",
				"error", err,
				"path", r.URL.Path,
				"method", r.Method)

			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if page.URL == "" {
			slog.Warn("empty URL provided",
				"path", r.URL.Path,
				"method", r.Method)

			http.Error(w, "URL cannot be empty", http.StatusBadRequest)
			return
		}

		slog.Debug("received ChronAm URL", "url", page.URL)

		p, ok := a.Store.Load(page.URL)
		if !ok {
			slog.Debug("unable to find URL in cache", "url", page.URL)

			err := GetRawText(&page)
			if err != nil {
				slog.Error("unable to download OCR text", "error", err, "url", page.URL)
			}

			err = a.RunPrompt(&page)
			if err != nil {
				slog.Error("error running prompt with Claude", "error", err)
				http.Error(w, "Unable to process that page", http.StatusInternalServerError)
				return
			}
			// Prompt was successful so store results in cache
			a.Store.Store(page.URL, page)
		} else {
			slog.Debug("found URL in cache", "url", page.URL)
			page = p.(ChronamPage) // Update our working page object with data from cache
		}

		// Return the JSON in response to the POST request
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(page)
	}
}

func GetRawText(page *ChronamPage) error {
	downloadURL := page.URL + "ocr.txt"
	slog.Debug("downloading OCR text", "url", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		slog.Error("failed to download URL", "url", downloadURL, "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		slog.Error("HTTP request failed", "error", err, "status_code", resp.StatusCode)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read response body", "error", err)
		return fmt.Errorf("failed to read response body: %w", err)
	}

	slog.Debug("download completed successfully", "url", downloadURL, "bytes_read", len(body))

	page.RawText = string(body)

	return nil
}
