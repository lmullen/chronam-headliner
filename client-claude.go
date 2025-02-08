package headliner

import (
	"fmt"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
)

type ClaudeResponse []int

func (a *App) RunPrompt() (ClaudeResponse, error) {
	var res ClaudeResponse

	slog.Debug("sending prompt to Claude")
	message, err := a.AIClient.Messages.New(a.ShutdownCtx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
		MaxTokens: anthropic.F(int64(1024)),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("Give me a JSON array of the integers 1, 2, 3, 4, 5. Return only the JSON.")),
		}),
	})
	if err != nil {
		slog.Error("failed sending prompt to Claude", "error", err)
		return nil, err
	}
	fmt.Printf("%+v\n", message.Content)

	return res, nil
}
