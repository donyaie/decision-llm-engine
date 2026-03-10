package engine

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/alidonyaie/decision-llm-engine/internal/model"
)

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
