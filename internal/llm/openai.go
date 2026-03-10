package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alidonyaie/decision-llm-engine/internal/config"
)

// OpenAIClient talks to an OpenAI-compatible chat completion endpoint.
type OpenAIClient struct {
	apiKey     string
	model      string
	baseURL    string
	httpClient *http.Client
}

type chatCompletionRequest struct {
	Model          string                `json:"model"`
	Messages       []chatCompletionInput `json:"messages"`
	Temperature    float64               `json:"temperature,omitempty"`
	ResponseFormat map[string]string     `json:"response_format,omitempty"`
}

type chatCompletionInput struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func newOpenAIClient(cfg config.LLMConfig) LLMClient {
	apiKey := strings.TrimSpace(cfg.OpenAIAPIKey)
	if apiKey == "" {
		return MockLLMClient{}
	}

	baseURL := strings.TrimSpace(cfg.OpenAIBaseURL)
	model := strings.TrimSpace(cfg.OpenAIModel)

	return &OpenAIClient{
		apiKey:  apiKey,
		model:   model,
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}


// Generate sends the prompt to the provider and returns raw text.
func (c *OpenAIClient) Generate(ctx context.Context, prompt, systemPrompt string) (string, error) {
	payload := chatCompletionRequest{
		Model: c.model,
		Messages: []chatCompletionInput{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: prompt},
		},
		Temperature:    0.2,
		ResponseFormat: map[string]string{"type": "json_object"},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal llm request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build llm request: %w", err)
	}

	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("send llm request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read llm response: %w", err)
	}

	var parsed chatCompletionResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return "", fmt.Errorf("decode llm response: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		if parsed.Error != nil && parsed.Error.Message != "" {
			return "", fmt.Errorf("llm provider error: %s", parsed.Error.Message)
		}
		return "", fmt.Errorf("llm provider returned status %d", response.StatusCode)
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("llm provider returned no choices")
	}

	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if content == "" {
		return "", fmt.Errorf("llm provider returned empty content")
	}

	return content, nil
}
