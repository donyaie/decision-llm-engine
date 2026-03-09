package tests

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/alidonyaie/decision-llm-engine/internal/config"
	"github.com/alidonyaie/decision-llm-engine/internal/engine"
	"github.com/alidonyaie/decision-llm-engine/internal/llm"
)

func TestRealLLMDecisionFlow(t *testing.T) {
	_ = godotenv.Load("../.env", ".env")

	llmConfig := config.LoadFromEnv().LLM
	provider := strings.ToLower(strings.TrimSpace(llmConfig.Provider))
	if provider == "" {
		if strings.TrimSpace(llmConfig.OpenAIAPIKey) != "" {
			provider = "openai"
		} else {
			provider = "mock"
		}
	}

	switch provider {
	case "openai", "openapi":
		if strings.TrimSpace(llmConfig.OpenAIAPIKey) == "" {
			t.Skip("OPENAI_API_KEY is required for OpenAI real LLM tests")
		}
	case "ollama":
		if strings.TrimSpace(llmConfig.OllamaModel) == "" {
			t.Skip("OLLAMA_MODEL is required for Ollama real LLM tests")
		}
	case "mock":
		t.Skip("real LLM test disabled; set LLM_PROVIDER=openai or LLM_PROVIDER=ollama to enable")
	default:
		t.Skip("real LLM test requires LLM_PROVIDER=openai or LLM_PROVIDER=ollama")
	}

	promptPath := filepath.Join("..", "prompts", "decision_prompt.txt")
	promptBuilder, err := engine.NewPromptBuilderFromFile(promptPath)
	if err != nil {
		t.Fatalf("load prompt template: %v", err)
	}

	decisionEngine := engine.NewDecisionEngine(promptBuilder, llm.NewClientFromEnv())

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	decision, err := decisionEngine.Analyze(ctx, "Should I switch from mobile development to AI/LLM engineering?")
	if err != nil {
		t.Fatalf("real LLM decision analysis failed: %v", err)
	}

	if strings.TrimSpace(decision.ProblemDefinition) == "" {
		t.Fatal("expected problem definition to be populated")
	}

	if len(decision.Options) == 0 {
		t.Fatal("expected at least one option from real LLM")
	}

	if len(decision.KeyFactors) == 0 {
		t.Fatal("expected at least one key factor from real LLM")
	}
}
