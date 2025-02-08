package headliner

import (
	"encoding/json"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
)

// Generate the type for the expected response from Claude
//go:generate go-jsonschema -p headliner -t --tags json ai-response-schema.json -o response.go

func (a *App) RunPrompt(page *ChronamPage) error {
	slog.Debug("sending prompt to Claude", "chronam", page.URL)
	prompt, err := a.MakePrompt(page.RawText)
	if err != nil {
		return err
	}

	message, err := a.AIClient.Messages.New(a.ShutdownCtx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(8192)),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(promptSystem),
		}),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
	})
	if err != nil {
		slog.Debug("error response from claude", "claude", message)
		return err
	}
	slog.Debug("claude usage", "usage", message.Usage)
	slog.Debug("estimated cost", "cost", calculateCost(message.Usage))

	for _, v := range message.Content {
		if v.Type == "text" {
			err = json.Unmarshal([]byte(v.Text), page)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Price in dollars per token for Claude 3.5 Sonnet
const inputPrice float64 = 3 / 1e6
const outputPrice float64 = 15 / 1e6

func calculateCost(u anthropic.Usage) float64 {
	return float64(u.InputTokens)*inputPrice + float64(u.OutputTokens)*outputPrice
}
