package engine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/alidonyaie/decision-llm-engine/internal/model"
)

// ValidateDecision checks that the decision matches the provided schema definition.
func ValidateDecision(decision model.Decision, schema string) error {
	trimmedSchema := strings.TrimSpace(schema)
	if trimmedSchema == "" {
		return fmt.Errorf("decision validation failed: schema is empty")
	}

	var schemaMap map[string]any
	if err := json.Unmarshal([]byte(trimmedSchema), &schemaMap); err != nil {
		return fmt.Errorf("decision validation failed: decode schema: %w", err)
	}

	decisionJSON, err := json.Marshal(decision)
	if err != nil {
		return fmt.Errorf("decision validation failed: encode decision: %w", err)
	}

	var decisionMap map[string]any
	if err := json.Unmarshal(decisionJSON, &decisionMap); err != nil {
		return fmt.Errorf("decision validation failed: decode decision: %w", err)
	}

	var validationErrors []error

	if isObjectSchema(schemaMap) {
		validationErrors = append(validationErrors, validateObjectSchema(decisionMap, schemaMap)...)
	} else {
		for key, expected := range schemaMap {
			actual, ok := decisionMap[key]
			if !ok {
				validationErrors = append(validationErrors, fmt.Errorf("missing required field %q", key))
				continue
			}

			if err := validateValueAgainstSchema(key, actual, expected); err != nil {
				validationErrors = append(validationErrors, err)
			}
		}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	return fmt.Errorf("decision validation failed: %w", errors.Join(validationErrors...))
}

func isObjectSchema(schema map[string]any) bool {
	_, hasProperties := schema["properties"]
	if !hasProperties {
		return false
	}

	typeValue, _ := schema["type"].(string)
	return strings.TrimSpace(typeValue) == "" || strings.TrimSpace(typeValue) == "object"
}

func validateObjectSchema(decisionMap map[string]any, schema map[string]any) []error {
	var validationErrors []error

	for _, key := range requiredKeys(schema["required"]) {
		if _, ok := decisionMap[key]; !ok {
			validationErrors = append(validationErrors, fmt.Errorf("missing required field %q", key))
		}
	}

	properties, _ := schema["properties"].(map[string]any)
	for key, expected := range properties {
		actual, ok := decisionMap[key]
		if !ok {
			continue
		}

		if err := validateValueAgainstSchema(key, actual, expected); err != nil {
			validationErrors = append(validationErrors, err)
		}
	}

	return validationErrors
}

func requiredKeys(raw any) []string {
	values, ok := raw.([]any)
	if !ok {
		return nil
	}

	required := make([]string, 0, len(values))
	for _, value := range values {
		key, ok := value.(string)
		if ok && strings.TrimSpace(key) != "" {
			required = append(required, key)
		}
	}

	return required
}

func validateValueAgainstSchema(key string, actual, expected any) error {
	switch expectedValue := expected.(type) {
	case string:
		return validatePrimitiveSchema(key, actual, expectedValue)
	case []any:
		return validateArraySchema(key, actual, expectedValue)
	case map[string]any:
		return validateTypedSchema(key, actual, expectedValue)
	default:
		return fmt.Errorf("unsupported schema type for field %q: %s", key, reflect.TypeOf(expected))
	}
}

func validateArrayItem(key string, actual, expected any) error {
	switch expectedValue := expected.(type) {
	case string:
		return validatePrimitiveSchema(key, actual, expectedValue)
	case map[string]any:
		return validateTypedSchema(key, actual, expectedValue)
	default:
		return fmt.Errorf("unsupported schema array item type for field %q", key)
	}
}

func validateStringValue(key string, actual any) error {
	actualString, ok := actual.(string)
	if !ok {
		return fmt.Errorf("field %q should be string", key)
	}
	if strings.TrimSpace(actualString) == "" {
		return fmt.Errorf("field %q is empty", key)
	}
	return nil
}

func validatePrimitiveSchema(key string, actual any, expectedType string) error {
	switch strings.TrimSpace(expectedType) {
	case "string":
		return validateStringValue(key, actual)
	default:
		return fmt.Errorf("unsupported schema type for field %q: %q", key, expectedType)
	}
}

func validateTypedSchema(key string, actual any, schema map[string]any) error {
	typeValue, _ := schema["type"].(string)
	switch strings.TrimSpace(typeValue) {
	case "string":
		return validateStringValue(key, actual)
	case "array":
		itemSchema, _ := schema["items"]
		return validateArrayWithItems(key, actual, itemSchema, minItems(schema))
	default:
		return fmt.Errorf("unsupported schema type for field %q: %q", key, typeValue)
	}
}

func validateArraySchema(key string, actual any, expectedItems []any) error {
	var itemSchema any
	if len(expectedItems) > 0 {
		itemSchema = expectedItems[0]
	}

	return validateArrayWithItems(key, actual, itemSchema, 1)
}

func validateArrayWithItems(key string, actual any, itemSchema any, minItems int) error {
	actualSlice, ok := actual.([]any)
	if !ok {
		return fmt.Errorf("field %q should be array", key)
	}
	if len(actualSlice) < minItems {
		if minItems > 1 {
			return fmt.Errorf("field %q must contain at least %d items", key, minItems)
		}
		return fmt.Errorf("field %q is empty", key)
	}

	if itemSchema == nil {
		return nil
	}

	for _, item := range actualSlice {
		if err := validateArrayItem(key, item, itemSchema); err != nil {
			return err
		}
	}

	return nil
}

func minItems(schema map[string]any) int {
	value, ok := schema["minItems"]
	if !ok {
		return 1
	}

	number, ok := value.(float64)
	if !ok || number < 1 {
		return 1
	}

	return int(number)
}
