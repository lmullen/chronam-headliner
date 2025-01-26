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

//go:embed ai-response-schema.json
var embeddedSchema []byte

var modelURL = "http://localhost:8080/v1/chat/completions"

type AIClient struct {
	client *http.Client
	schema ResponseSchema
}

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIRequestBody struct {
	Model    string      `json:"model"`
	Messages []AIMessage `json:"messages"`
	// ResponseFormat ResponseFormat `json:"response_format"`
}

type ResponseSchema struct {
	Schema string `json:"$schema"`
	Type   string `json:"type"`
	Items  []struct {
		Properties struct {
			Headline struct {
				Type string `json:"type"`
			} `json:"headline"`
			Body struct {
				Type string `json:"type"`
			} `json:"body"`
		} `json:"properties"`
		Required []string `json:"required"`
	} `json:"items"`
}

type ResponseFormat struct {
	Type       string         `json:"type"`
	JSONSchema ResponseSchema `json:"json_schema"`
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

	// var schema ResponseSchema
	// err := json.Unmarshal(embeddedSchema, &schema)
	// if err != nil {
	// 	slog.Error("failed to decode JSON Schema", "error", err)
	// 	return nil
	// }
	// fmt.Println(string(embeddedSchema))
	// fmt.Println(schema)

	var schema ResponseSchema

	if err := json.Unmarshal(embeddedSchema, &schema); err != nil {
		fmt.Printf("Error unmarshalling: %v\n", err)
		return nil
	}

	c.schema = schema

	return &c
}

func (c *AIClient) ConstructAIRequest(page *ChronamPage) (*http.Request, error) {

	prompt := embeddedPrompt + "```\n" + page.RawText + "\n```"

	reqBody := AIRequestBody{
		Model: "LLaMA_CPP",
		// ResponseFormat: ResponseFormat{
		// 	Type:       "json_schema",
		// 	JSONSchema: c.schema,
		// },
		Messages: []AIMessage{
			{
				Role:    "user",
				Content: prompt,
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
	slog.Info("prompting the model for a response")
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
	fmt.Println(response)
	fmt.Println(response.Choices[0].Message.Content)

	return nil
}
