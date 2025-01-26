package headliner

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ChronamPage struct {
	URL     string `json:"url"`
	RawText string `json:"raw_text"`
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

		slog.Info("received ChronAm URL", "url", page)

		err := GetRawText(&page)
		if err != nil {
			slog.Error("unable to download OCR text", "error", err, "url", page.URL)
		}

	}
}

func GetRawText(page *ChronamPage) error {

	downloadURL := page.URL + "ocr.txt"
	slog.Info("downloading OCR text", "url", downloadURL)

	resp, err := http.Get(downloadURL)
	if err != nil {
		slog.Error("failed to download URL", "error", err)
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

	slog.Info("download completed successfully", "bytes_read", len(body))

	page.RawText = string(body)

	return nil

}
