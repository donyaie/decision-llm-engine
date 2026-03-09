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

	prompt, err := builder.Build("Should I move to Dubai for a software job?")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !strings.Contains(prompt, "Should I move to Dubai for a software job?") {
		t.Fatalf("expected prompt to include question, got %q", prompt)
	}
}

func TestDecisionValidation(t *testing.T) {
	decision := model.Decision{
		ProblemDefinition: "Moving to Dubai for a software job",
		DecisionType:      "life/career",
		Options:           []string{"Move to Dubai", "Stay in current country"},
		KeyFactors:        []string{"salary", "cost of living"},
	}

	if err := validation.ValidateDecision(decision); err != nil {
		t.Fatalf("expected decision to be valid, got %v", err)
	}
}

func TestJSONParsing(t *testing.T) {
	raw := "```json\n{\n  \"problem_definition\": \"Moving to Dubai for a software job\",\n  \"decision_type\": \"life/career\",\n  \"options\": [\"Move to Dubai\", \"Stay in current country\"],\n  \"key_factors\": [\"salary\", \"cost of living\"],\n  \"risks\": [\"job instability\"],\n  \"unknowns\": [\"exact salary offer\"],\n  \"recommended_next_questions\": [\"What salary is offered?\"]\n}\n```"

	decision, err := engine.ParseDecision(raw)
	if err != nil {
		t.Fatalf("expected parsable json, got %v", err)
	}

	if decision.ProblemDefinition != "Moving to Dubai for a software job" {
		t.Fatalf("unexpected problem definition: %q", decision.ProblemDefinition)
	}
}

func TestEngineDecisionFlow(t *testing.T) {
	client := &stubClient{responses: []string{
		`{
			"problem_definition": "Moving to Dubai for a software job",
			"decision_type": "life/career",
			"options": ["Move to Dubai", "Stay in current country"],
			"key_factors": ["salary", "cost of living"],
			"risks": ["job instability"],
			"unknowns": ["exact salary offer"],
			"recommended_next_questions": ["What salary is offered?"]
		}`,
	}}

	decisionEngine := engine.NewDecisionEngine(engine.NewPromptBuilder("User Question:\n{{user_input}}"), client)
	decision, err := decisionEngine.Analyze(context.Background(), "Should I move to Dubai for a software job?")
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
