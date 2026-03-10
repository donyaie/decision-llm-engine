package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/alidonyaie/decision-llm-engine/internal/llm"
	"github.com/alidonyaie/decision-llm-engine/internal/model"
	"github.com/alidonyaie/decision-llm-engine/internal/reliability"
)

// DecisionEngine coordinates prompt construction, LLM generation, repair, and validation.
type DecisionEngine struct {
	promptBuilder *PromptBuilder
	client        llm.Client
	maxRetries    int
	retryDelay    time.Duration
}

// NewDecisionEngine builds an engine with production-minded defaults.
func NewDecisionEngine(promptBuilder *PromptBuilder, client llm.Client) *DecisionEngine {
	return &DecisionEngine{
		promptBuilder: promptBuilder,
		client:        client,
		maxRetries:    3,
		retryDelay:    250 * time.Millisecond,
	}
}

// Analyze converts a user question into a validated decision object.
func (e *DecisionEngine) Analyze(ctx context.Context, question string) (model.Decision, error) {
	var zero model.Decision

	prompts, err := e.buildPrompts(question)
	if err != nil {
		return zero, err
	}

	raw, err := e.generate(ctx, prompts)
	if err != nil {
		return zero, err
	}

	decision, err := e.parseAndRepair(ctx, raw, prompts.System)
	if err != nil {
		return zero, err
	}

	decision, err = e.finalizeDecision(decision)
	if err != nil {
		return zero, err
	}

	return decision, nil
}

func (e *DecisionEngine) buildPrompts(question string) (PromptSet, error) {
	prompts, err := e.promptBuilder.Build(question)
	if err != nil {
		return PromptSet{}, fmt.Errorf("build prompt: %w", err)
	}

	return prompts, nil
}

func (e *DecisionEngine) generate(ctx context.Context, prompts PromptSet) (string, error) {
	raw, err := reliability.DoWithResult(ctx, e.maxRetries, e.retryDelay, func() (string, error) {
		return e.client.Generate(ctx, prompts.User, prompts.System)
	})
	if err != nil {
		return "", fmt.Errorf("generate decision: %w", err)
	}

	return raw, nil
}

func (e *DecisionEngine) parseAndRepair(ctx context.Context, raw, systemPrompt string) (model.Decision, error) {
	decision, err := ParseDecision(raw)
	if err == nil {
		return decision, nil
	}

	repairedRaw, repairErr := e.repairResponse(ctx, raw, systemPrompt)
	if repairErr != nil {
		return model.Decision{}, fmt.Errorf("parse decision response: %w", err)
	}

	decision, err = ParseDecision(repairedRaw)
	if err != nil {
		return model.Decision{}, fmt.Errorf("parse repaired decision response: %w", err)
	}

	return decision, nil
}

func (e *DecisionEngine) finalizeDecision(decision model.Decision) (model.Decision, error) {
	decision = normalizeDecision(decision)
	if err := ValidateDecision(decision, e.promptBuilder.schema); err != nil {
		return model.Decision{}, err
	}

	return decision, nil
}

func (e *DecisionEngine) repairResponse(ctx context.Context, raw, systemPrompt string) (string, error) {
	if repaired := basicJSONRepair(raw); repaired != "" {
		if _, err := ParseDecision(repaired); err == nil {
			return repaired, nil
		}
	}

	repairPrompt := fmt.Sprintf(`You fix malformed JSON.
Return ONLY valid JSON that matches the original schema.
Do not add commentary.

Original response:
%s`, sanitizeJSON(raw))

	repaired, err := e.client.Generate(ctx, repairPrompt, systemPrompt)
	if err != nil {
		return "", fmt.Errorf("repair malformed JSON: %w", err)
	}

	return repaired, nil
}
