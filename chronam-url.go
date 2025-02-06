package headliner

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type ChronamPage struct {
	URL      string          `json:"url"`
	RawText  string          `json:"raw_text"`
	Articles ChronamArticles `json:"articles"`
}

type ChronamArticles []struct {
	Headline string `json:"headline"`
	Body     string `json:"body"`
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

		err := GetRawText(&page)
		if err != nil {
			slog.Error("unable to download OCR text", "error", err, "url", page.URL)
		}

		// req, err := a.AIClient.ConstructAIRequest(&page)
		// if err != nil {
		// 	slog.Error("error contructing request for this page", "url", page.URL, "error", err)
		// 	return
		// }

		// err = a.AIClient.RunPrompt(&page, req)
		// if err != nil {
		// 	slog.Error("error making request to model", "url", page.URL, "error", err)
		// 	return
		// }

		// for i, art := range page.Articles {
		// 	fmt.Println("Count: ", i)
		// 	fmt.Println("Headline: ", art.Headline)
		// 	// fmt.Println("Body: ", art.Body)
		// }

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
