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

// OllamaClient talks to Ollama's generate endpoint.
type OllamaClient struct {
	model      string
	baseURL    string
	httpClient *http.Client
}

type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Format string `json:"format,omitempty"`
}

type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Error    string `json:"error,omitempty"`
}

func newOllamaClient(cfg config.LLMConfig) LLMClient {
	baseURL := strings.TrimSpace(cfg.OllamaBaseURL)
	model := strings.TrimSpace(cfg.OllamaModel)

	return &OllamaClient{
		model:   model,
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Generate sends the prompt to Ollama and returns raw text.
func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	payload := ollamaGenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Format: "json",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal ollama request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("build ollama request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("send ollama request: %w", err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read ollama response: %w", err)
	}

	var parsed ollamaGenerateResponse
	if err := json.Unmarshal(responseBody, &parsed); err != nil {
		return "", fmt.Errorf("decode ollama response: %w", err)
	}

	if response.StatusCode >= http.StatusBadRequest {
		if parsed.Error != "" {
			return "", fmt.Errorf("ollama error: %s", parsed.Error)
		}
		return "", fmt.Errorf("ollama returned status %d", response.StatusCode)
	}

	content := strings.TrimSpace(parsed.Response)
	if content == "" {
		return "", fmt.Errorf("ollama returned empty content")
	}

	return content, nil
}
