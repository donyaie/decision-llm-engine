package validation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/alidonyaie/decision-llm-engine/internal/model"
)

// ValidateDecision checks that required fields are present and non-empty.
func ValidateDecision(decision model.Decision) error {
	var validationErrors []error

	if strings.TrimSpace(decision.ProblemDefinition) == "" {
		validationErrors = append(validationErrors, errors.New("problem definition missing"))
	}

	if strings.TrimSpace(decision.DecisionType) == "" {
		validationErrors = append(validationErrors, errors.New("decision type missing"))
	}

	if len(nonEmpty(decision.Options)) == 0 {
		validationErrors = append(validationErrors, errors.New("no options extracted"))
	}

	if len(nonEmpty(decision.KeyFactors)) == 0 {
		validationErrors = append(validationErrors, errors.New("no key factors extracted"))
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return fmt.Errorf("decision validation failed: %w", errors.Join(validationErrors...))
}

func nonEmpty(values []string) []string {
	clean := make([]string, 0, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			clean = append(clean, trimmed)
		}
	}
	return clean
}
