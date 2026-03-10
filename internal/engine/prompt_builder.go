package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultUserPromptTemplate = "Schema:\n{{decision_schema}}\n\nUser Question:\n{{user_input}}"
)

// PromptBuilder loads and renders the prompt template used for decision analysis.
type PromptBuilder struct {
	system string
	schema string
}

// PromptSet contains the rendered prompts for a single request.
type PromptSet struct {
	User   string
	System string
	Schema string
}

// NewPromptBuilder creates a prompt builder from raw templates.
func NewPromptBuilder(system, schema string) *PromptBuilder {
	return &PromptBuilder{system: system, schema: schema}
}

// NewPromptBuilderFromFiles loads prompt templates from disk.
func NewPromptBuilderFromFiles(systemPath, schemaPath string) (*PromptBuilder, error) {
	systemTemplate, err := readPromptTemplate(systemPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt template: %w", err)
	}

	schemaTemplate, err := readPromptTemplate(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("load schema template: %w", err)
	}

	return NewPromptBuilder(systemTemplate, schemaTemplate), nil
}

func readPromptTemplate(path string) (string, error) {
	trimmedPath := strings.TrimSpace(path)
	if trimmedPath == "" {
		return "", nil
	}

	candidatePaths := []string{trimmedPath}
	if !filepath.IsAbs(trimmedPath) {
		candidatePaths = append(candidatePaths,
			filepath.Join("..", trimmedPath),
			filepath.Join("..", "..", trimmedPath),
		)
	}

	var content []byte
	var err error
	for _, candidatePath := range candidatePaths {
		content, err = os.ReadFile(candidatePath)
		if err == nil {
			return string(content), nil
		}
	}

	return "", fmt.Errorf("read prompt template: %w", err)
}

// Build renders the user and system prompts.
func (b *PromptBuilder) Build(question string) (PromptSet, error) {
	if b == nil {
		return PromptSet{}, fmt.Errorf("prompt builder is empty")
	}

	userPrompt, err := buildUserPrompt(b.schema, question)
	if err != nil {
		return PromptSet{}, err
	}

	systemPrompt, err := buildOptionalPrompt(b.system)
	if err != nil {
		return PromptSet{}, fmt.Errorf("build system prompt: %w", err)
	}

	return PromptSet{
		User:   userPrompt,
		System: systemPrompt,
		Schema: strings.TrimSpace(b.schema),
	}, nil
}

func buildUserPrompt(schema, question string) (string, error) {
	trimmedSchema := strings.TrimSpace(schema)
	if trimmedSchema == "" {
		return "", fmt.Errorf("decision schema is required")
	}

	trimmedQuestion := strings.TrimSpace(question)
	if trimmedQuestion == "" {
		return "", fmt.Errorf("question is required")
	}

	userPrompt := strings.ReplaceAll(defaultUserPromptTemplate, "{{decision_schema}}", trimmedSchema)
	userPrompt = strings.ReplaceAll(userPrompt, "{{user_input}}", trimmedQuestion)
	return userPrompt, nil
}

func buildOptionalPrompt(template string) (string, error) {
	if strings.TrimSpace(template) == "" {
		return "", nil
	}

	return strings.TrimSpace(template), nil
}
