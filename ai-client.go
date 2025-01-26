package headliner

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

//go:embed ai-prompt.txt
var embeddedPrompt string

var modelURL = "http://localhost:8080/v1/chat/completions"

type AIClient struct {
	client *http.Client
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIRequestBody struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
}

type ModelResponse struct {
	Choices []struct {
		FinishReason string      `json:"finish_reason"`
		Index        int         `json:"index"`
		LogProbs     interface{} `json:"logprobs"`
		Message      struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
	} `json:"choices"`
	Created           int64  `json:"created"`
	ID                string `json:"id"`
	Model             string `json:"model"`
	Object            string `json:"object"`
	SystemFingerprint string `json:"system_fingerprint"`
	Usage             struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewAIClient(ctx context.Context) *AIClient {
	c := AIClient{}

	c.client = &http.Client{}

	return &c

}

func (c *AIClient) ConstructAIRequest() (*http.Request, error) {
	reqBody := AIRequestBody{
		Model: "LLaMA_CPP",
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: embeddedPrompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		slog.Error("failed to marshal JSON for AI request", "error", err)
		return nil, fmt.Errorf("marshal JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, modelURL, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("failed to create request", "error", err)
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *AIClient) RunPrompt(req *http.Request) error {
	resp, err := c.client.Do(req)
	if err != nil {
		slog.Error("failed to send request", "error", err)
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Error("unexpected status code", "status", resp.StatusCode)
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Need to decode the response correctly
	var response ModelResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("failed to decode response", "error", err)
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}
