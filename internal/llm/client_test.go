package llm

import (
	"testing"

	"github.com/alidonyaie/decision-llm-engine/internal/config"
)

func TestNewClientUsesMockLLMClientByDefault(t *testing.T) {
	generatedClient, ok := NewClientFromConfig(config.LLMConfig{}).(*client)
	if !ok {
		t.Fatal("expected concrete client implementation")
	}

	if _, ok := generatedClient.provider.(MockLLMClient); !ok {
		t.Fatalf("expected mock llm client, got %T", generatedClient.provider)
	}
}

func TestNewClientUsesOpenAIClient(t *testing.T) {
	generatedClient, ok := NewClientFromConfig(config.LLMConfig{
		Provider:      ProviderOpenAI,
		OpenAIAPIKey:  "test-key",
		OpenAIModel:   "gpt-4.1-mini",
		OpenAIBaseURL: "https://api.openai.com",
	}).(*client)
	if !ok {
		t.Fatal("expected concrete client implementation")
	}

	if _, ok := generatedClient.provider.(*OpenAIClient); !ok {
		t.Fatalf("expected OpenAI client, got %T", generatedClient.provider)
	}
}

func TestNewClientUsesOllamaClient(t *testing.T) {
	generatedClient, ok := NewClientFromConfig(config.LLMConfig{
		Provider:      ProviderOllama,
		OllamaModel:   "llama3.2",
		OllamaBaseURL: "http://localhost:11434",
	}).(*client)
	if !ok {
		t.Fatal("expected concrete client implementation")
	}

	if _, ok := generatedClient.provider.(*OllamaClient); !ok {
		t.Fatalf("expected Ollama client, got %T", generatedClient.provider)
	}
}

func TestNormalizeProviderAlias(t *testing.T) {
	if got := normalizeProvider("openapi"); got != ProviderOpenAI {
		t.Fatalf("expected alias to normalize to %q, got %q", ProviderOpenAI, got)
	}
}

func TestProviderFromEnvPrefersNewKey(t *testing.T) {
	t.Setenv("LLM_PROVIDER", "ollama")

	if got := providerFromEnv(); got != "ollama" {
		t.Fatalf("expected providerFromEnv to prefer LLM_PROVIDER, got %q", got)
	}
}
