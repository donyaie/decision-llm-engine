package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/alidonyaie/decision-llm-engine/internal/llm"
	"github.com/alidonyaie/decision-llm-engine/internal/model"
	"github.com/alidonyaie/decision-llm-engine/internal/reliability"
	"github.com/alidonyaie/decision-llm-engine/internal/validation"
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

	prompt, err := e.promptBuilder.Build(question)
	if err != nil {
		return zero, fmt.Errorf("build prompt: %w", err)
	}

	raw, err := reliability.DoWithResult(ctx, e.maxRetries, e.retryDelay, func() (string, error) {
		return e.client.Generate(ctx, prompt)
	})
	if err != nil {
		return zero, fmt.Errorf("generate decision: %w", err)
	}

	decision, err := ParseDecision(raw)
	if err != nil {
		repairedRaw, repairErr := e.repairResponse(ctx, raw)
		if repairErr != nil {
			return zero, fmt.Errorf("parse decision response: %w", err)
		}

		decision, err = ParseDecision(repairedRaw)
		if err != nil {
			return zero, fmt.Errorf("parse repaired decision response: %w", err)
		}
	}

	decision = normalizeDecision(decision)
	if err := validation.ValidateDecision(decision); err != nil {
		return zero, err
	}

	return decision, nil
}

// ParseDecision decodes a decision object from an LLM response.
func ParseDecision(raw string) (model.Decision, error) {
	var decision model.Decision

	candidates := []string{
		strings.TrimSpace(raw),
		extractJSONObject(raw),
	}

	for _, candidate := range candidates {
		candidate = sanitizeJSON(candidate)
		if candidate == "" {
			continue
		}
		if err := json.Unmarshal([]byte(candidate), &decision); err == nil {
			return decision, nil
		}
	}

	return model.Decision{}, fmt.Errorf("response is not valid decision JSON")
}

func (e *DecisionEngine) repairResponse(ctx context.Context, raw string) (string, error) {
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

	repaired, err := e.client.Generate(ctx, repairPrompt)
	if err != nil {
		return "", fmt.Errorf("repair malformed JSON: %w", err)
	}

	return repaired, nil
}

func sanitizeJSON(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSuffix(trimmed, "```")
	return strings.TrimSpace(trimmed)
}

func extractJSONObject(raw string) string {
	trimmed := sanitizeJSON(raw)
	start := strings.Index(trimmed, "{")
	end := strings.LastIndex(trimmed, "}")
	if start == -1 || end == -1 || start >= end {
		return ""
	}
	return trimmed[start : end+1]
}

func basicJSONRepair(raw string) string {
	repaired := extractJSONObject(raw)
	if repaired == "" {
		repaired = sanitizeJSON(raw)
	}
	if repaired == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		"“", `"`,
		"”", `"`,
	)

	previous := repaired
	trailingCommaPattern := regexp.MustCompile(`,(\s*[}\]])`)
	for {
		repaired = replacer.Replace(previous)
		repaired = trailingCommaPattern.ReplaceAllString(repaired, `$1`)
		if repaired == previous {
			break
		}
		previous = repaired
	}

	return repaired
}

func normalizeDecision(decision model.Decision) model.Decision {
	decision.ProblemDefinition = strings.TrimSpace(decision.ProblemDefinition)
	decision.DecisionType = strings.TrimSpace(decision.DecisionType)
	decision.Options = normalizeList(decision.Options)
	decision.KeyFactors = normalizeList(decision.KeyFactors)
	decision.Risks = normalizeList(decision.Risks)
	decision.Unknowns = normalizeList(decision.Unknowns)
	decision.NextQuestions = normalizeList(decision.NextQuestions)
	return decision
}

func normalizeList(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	clean := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		key := strings.ToLower(trimmed)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		clean = append(clean, trimmed)
	}
	return clean
}
