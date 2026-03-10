package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alidonyaie/decision-llm-engine/internal/model"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

const decisionSchemaResource = "decision_schema.json"

// ValidateDecision checks that the decision matches the provided schema definition.
func ValidateDecision(decision model.Decision, schema string) error {
	trimmedSchema := strings.TrimSpace(schema)
	if trimmedSchema == "" {
		return fmt.Errorf("decision validation failed: schema is empty")
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource(decisionSchemaResource, strings.NewReader(trimmedSchema)); err != nil {
		return fmt.Errorf("decision validation failed: compile schema resource: %w", err)
	}

	compiledSchema, err := compiler.Compile(decisionSchemaResource)
	if err != nil {
		return fmt.Errorf("decision validation failed: compile schema: %w", err)
	}

	decisionJSON, err := json.Marshal(decision)
	if err != nil {
		return fmt.Errorf("decision validation failed: encode decision: %w", err)
	}

	var decisionValue any
	if err := json.Unmarshal(decisionJSON, &decisionValue); err != nil {
		return fmt.Errorf("decision validation failed: decode decision: %w", err)
	}

	if err := compiledSchema.Validate(decisionValue); err != nil {
		return fmt.Errorf("decision validation failed: %w", err)
	}

	return nil
}
