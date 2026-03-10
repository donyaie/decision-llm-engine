package config

import (
	"os"
	"strings"
)

const (
	defaultPort             = "8080"
	defaultSchemaPath       = "prompts/decision_schema.json"
	defaultSystemPromptPath = "prompts/system_prompt.txt"
	defaultOpenAIBaseURL    = "https://api.openai.com"
	defaultOpenAIModel      = "gpt-4.1-mini"
	defaultOllamaBaseURL    = "http://localhost:11434"
	defaultOllamaModel      = "llama3.2"
)

// Config represents the application's environment-driven configuration.
type Config struct {
	Server ServerConfig
	LLM    LLMConfig
}

// ServerConfig holds HTTP server and prompt settings.
type ServerConfig struct {
	Port string
}

// LLMConfig holds provider-specific LLM settings.
type LLMConfig struct {
	Provider         string
	SchemaPath       string
	SystemPromptPath string
	OpenAIAPIKey     string
	OpenAIBaseURL    string
	OpenAIModel      string
	OllamaBaseURL    string
	OllamaModel      string
}

// LoadFromEnv builds a Config from process environment variables.
func LoadFromEnv() Config {
	cfg := Config{
		Server: ServerConfig{
			Port: envOrDefault("PORT", defaultPort),
		},
		LLM: LLMConfig{
			Provider:         strings.TrimSpace(os.Getenv("LLM_PROVIDER")),
			SchemaPath:       envOrDefault("SCHEMA_PATH", defaultSchemaPath),
			SystemPromptPath: envOrDefault("SYSTEM_PROMPT_PATH", defaultSystemPromptPath),
			OpenAIAPIKey:     strings.TrimSpace(os.Getenv("OPENAI_API_KEY")),
			OpenAIBaseURL:    envOrDefault("OPENAI_BASE_URL", defaultOpenAIBaseURL),
			OpenAIModel:      envOrDefault("OPENAI_MODEL", defaultOpenAIModel),
			OllamaBaseURL:    envOrDefault("OLLAMA_BASE_URL", defaultOllamaBaseURL),
			OllamaModel:      envOrDefault("OLLAMA_MODEL", defaultOllamaModel),
		},
	}

	return cfg
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}
