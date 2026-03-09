package llm

import (
	"context"
	"strings"

	"github.com/alidonyaie/decision-llm-engine/internal/config"
)

const (
	ProviderMock   = "mock"
	ProviderOpenAI = "openai"
	ProviderOllama = "ollama"
)

const (
	providerEnvVar = "LLM_PROVIDER"
)

// Client abstracts LLM generation.
type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// LLMClient is the provider-specific implementation behind the shared client.
type LLMClient interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

type client struct {
	provider LLMClient
}

// NewClient wraps a provider client with the shared client interface.
func NewClient(provider LLMClient) Client {
	if provider == nil {
		provider = MockLLMClient{}
	}

	return &client{provider: provider}
}

// NewClientFromEnv returns an LLM client based on environment configuration.
func NewClientFromEnv() Client {
	return NewClientFromConfig(config.LoadFromEnv().LLM)
}

// NewClientFromConfig returns an LLM client using the provided config.
func NewClientFromConfig(cfg config.LLMConfig) Client {
	return NewClient(newLLMClientFromConfig(cfg))
}

func (c *client) Generate(ctx context.Context, prompt string) (string, error) {
	return c.provider.Generate(ctx, prompt)
}

func newLLMClientFromConfig(cfg config.LLMConfig) LLMClient {
	provider := normalizeProvider(cfg.Provider)
	provider = defaultProvider(provider)

	switch provider {
	case ProviderOpenAI:
		return newOpenAIClient(cfg)
	case ProviderOllama:
		return newOllamaClient(cfg)
	default:
		return MockLLMClient{}
	}
}

func providerFromEnv() string {
	return strings.TrimSpace(config.LoadFromEnv().LLM.Provider)
}

func normalizeProvider(provider string) string {
	normalized := strings.ToLower(strings.TrimSpace(provider))
	if normalized == "openapi" {
		return ProviderOpenAI
	}
	return normalized
}

func defaultProvider(provider string) string {
	if provider != "" {
		return provider
	}

	if strings.TrimSpace(config.LoadFromEnv().LLM.OpenAIAPIKey) != "" {
		return ProviderOpenAI
	}

	return ProviderMock
}
