package engine

import (
	"fmt"
	"os"
	"strings"
)

const promptPlaceholder = "{{user_input}}"

// PromptBuilder loads and renders the prompt template used for decision analysis.
type PromptBuilder struct {
	template string
}

// NewPromptBuilder creates a prompt builder from a raw template string.
func NewPromptBuilder(template string) *PromptBuilder {
	return &PromptBuilder{template: template}
}

// NewPromptBuilderFromFile loads a prompt template from disk.
func NewPromptBuilderFromFile(path string) (*PromptBuilder, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read prompt template: %w", err)
	}

	return NewPromptBuilder(string(content)), nil
}

// Build injects the user's question into the prompt template.
func (b *PromptBuilder) Build(question string) (string, error) {
	if b == nil || strings.TrimSpace(b.template) == "" {
		return "", fmt.Errorf("prompt template is empty")
	}

	trimmedQuestion := strings.TrimSpace(question)
	if trimmedQuestion == "" {
		return "", fmt.Errorf("question is required")
	}

	if strings.Contains(b.template, promptPlaceholder) {
		return strings.ReplaceAll(b.template, promptPlaceholder, trimmedQuestion), nil
	}

	return strings.TrimSpace(b.template) + "\n\nUser Question:\n" + trimmedQuestion, nil
}
