package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/alidonyaie/decision-llm-engine/internal/engine"
	"github.com/alidonyaie/decision-llm-engine/internal/model"
	"github.com/alidonyaie/decision-llm-engine/internal/validation"
)

type stubClient struct {
	responses []string
	calls     int
}

func (s *stubClient) Generate(_ context.Context, _ string) (string, error) {
	index := s.calls
	if index >= len(s.responses) {
		index = len(s.responses) - 1
	}
	s.calls++
	return s.responses[index], nil
}

func TestPromptBuilder(t *testing.T) {
	builder := engine.NewPromptBuilder("Schema\nUser Question:\n{{user_input}}")

	prompt, err := builder.Build("Should I switch from mobile development to AI/LLM engineering?")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(prompt, "Should I switch from mobile development to AI/LLM engineering?") {
		t.Fatalf("expected prompt to include question, got %q", prompt)
	}
}

func TestDecisionValidation(t *testing.T) {
	decision := model.Decision{
		ProblemDefinition: "Switching from mobile development to AI/LLM engineering",
		DecisionType:      "career",
		Options:           []string{"Transition into AI/LLM engineering", "Stay in mobile development"},
		KeyFactors:        []string{"market demand", "learning curve"},
	}

	if err := validation.ValidateDecision(decision); err != nil {
		t.Fatalf("expected decision to be valid, got %v", err)
	}
}

func TestJSONParsing(t *testing.T) {
	raw := "```json\n{\n  \"problem_definition\": \"Switching from mobile development to AI/LLM engineering\",\n  \"decision_type\": \"career\",\n  \"options\": [\"Transition into AI/LLM engineering\", \"Stay in mobile development\"],\n  \"key_factors\": [\"market demand\", \"learning curve\"],\n  \"risks\": [\"temporary productivity dip\"],\n  \"unknowns\": [\"time needed to become job-ready\"],\n  \"recommended_next_questions\": [\"What AI/LLM skills are most required for target roles?\"]\n}\n```"

	decision, err := engine.ParseDecision(raw)
	if err != nil {
		t.Fatalf("expected parsable json, got %v", err)
	}

	if decision.ProblemDefinition != "Switching from mobile development to AI/LLM engineering" {
		t.Fatalf("unexpected problem definition: %q", decision.ProblemDefinition)
	}
}

func TestEngineDecisionFlow(t *testing.T) {
	client := &stubClient{responses: []string{
		`{
			"problem_definition": "Switching from mobile development to AI/LLM engineering",
			"decision_type": "career",
			"options": ["Transition into AI/LLM engineering", "Stay in mobile development"],
			"key_factors": ["market demand", "learning curve"],
			"risks": ["temporary productivity dip"],
			"unknowns": ["time needed to become job-ready"],
			"recommended_next_questions": ["What AI/LLM skills are most required for target roles?"]
		}`,
	}}

	decisionEngine := engine.NewDecisionEngine(engine.NewPromptBuilder("User Question:\n{{user_input}}"), client)
	decision, err := decisionEngine.Analyze(context.Background(), "Should I switch from mobile development to AI/LLM engineering?")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(decision.Options) != 2 {
		t.Fatalf("expected 2 options, got %d", len(decision.Options))
	}
}

func TestEngineRepairsMalformedJSON(t *testing.T) {
	client := &stubClient{responses: []string{
		`{
			"problem_definition": "Building a SaaS startup versus staying employed",
			"decision_type": "business/career",
			"options": ["Build a SaaS startup", "Stay employed"],
			"key_factors": ["runway", "market validation"],
			"risks": ["startup failure"],
			"unknowns": ["customer demand"],
			"recommended_next_questions": ["How many months of runway are available?"],
		}`,
	}}

	decisionEngine := engine.NewDecisionEngine(engine.NewPromptBuilder("User Question:\n{{user_input}}"), client)
	decision, err := decisionEngine.Analyze(context.Background(), "Should I build a SaaS startup or stay employed?")
	if err != nil {
		t.Fatalf("expected repair to succeed, got %v", err)
	}

	if decision.DecisionType != "business/career" {
		t.Fatalf("unexpected decision type: %q", decision.DecisionType)
	}
}
