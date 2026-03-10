package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// MockLLMClient provides deterministic demo-friendly responses when no provider is configured.
type MockLLMClient struct{}

// Generate synthesizes realistic JSON for local development and tests.
func (MockLLMClient) Generate(_ context.Context, prompt, _ string) (string, error) {
	if strings.Contains(strings.ToLower(prompt), "fix this json") {
		return repairPromptResponse(prompt), nil
	}

	question := extractQuestion(prompt)
	return generateMockDecision(question)
}

func extractQuestion(prompt string) string {
	marker := "User Question:"
	index := strings.LastIndex(prompt, marker)
	if index == -1 {
		return strings.TrimSpace(prompt)
	}
	return strings.TrimSpace(prompt[index+len(marker):])
}

func generateMockDecision(question string) (string, error) {
	lower := strings.ToLower(strings.TrimSpace(question))

	type decisionPayload struct {
		ProblemDefinition string   `json:"problem_definition"`
		DecisionType      string   `json:"decision_type"`
		Options           []string `json:"options"`
		KeyFactors        []string `json:"key_factors"`
		Risks             []string `json:"risks"`
		Unknowns          []string `json:"unknowns"`
		NextQuestions     []string `json:"recommended_next_questions"`
	}

	payload := decisionPayload{
		ProblemDefinition: strings.TrimSuffix(strings.TrimSpace(question), "?"),
		DecisionType:      classifyDecisionType(lower),
		Options:           []string{"Proceed with the change", "Keep the current path"},
		KeyFactors:        []string{"financial impact", "time horizon", "upside potential", "downside risk"},
		Risks:             []string{"insufficient information", "overconfidence", "unexpected trade-offs"},
		Unknowns:          []string{"timeline details", "full cost", "best alternative"},
		NextQuestions:     []string{"What is the expected upside?", "What is the downside if it fails?", "What information is still missing?"},
	}

	switch {
	case strings.Contains(lower, "mobile") && (strings.Contains(lower, "ai") || strings.Contains(lower, "llm")):
		payload = decisionPayload{
			ProblemDefinition: "Should I stop being a mobile developer and move into AI?",
			DecisionType:      "Career Change",
			Options:           []string{"Stay as a Mobile Developer", "Move into AI"},
			KeyFactors:        []string{"Skill Set", "Financial Gain", "Personal Interest", "Job Security"},
			Risks:             []string{"Loss of Current Income", "Steep Learning Curve", "Uncertainty in Future Demand"},
			Unknowns:          []string{"Current Demand for AI Talent", "Time and Effort Required to Adapt"},
			NextQuestions:     []string{"What are the current demand and salary ranges for AI developers in my location?", "How much time and effort will it take to adapt my skill set to AI?", "What are the potential career growth opportunities in AI?"},
		}
	case strings.Contains(lower, "master") && strings.Contains(lower, "ai"):
		payload = decisionPayload{
			ProblemDefinition: "Pursuing a master's degree in AI",
			DecisionType:      "education/career",
			Options:           []string{"Start a master's degree in AI", "Keep working and upskill independently"},
			KeyFactors:        []string{"tuition cost", "career acceleration", "opportunity cost", "program quality", "industry relevance"},
			Risks:             []string{"debt burden", "weak return on investment", "curriculum mismatch"},
			Unknowns:          []string{"admission chances", "scholarship availability", "expected salary uplift"},
			NextQuestions:     []string{"What is the total cost of the degree?", "How strong is the program's placement record?", "Can the same skills be gained through work experience?"},
		}
	case strings.Contains(lower, "saas") || strings.Contains(lower, "stay employed"):
		payload = decisionPayload{
			ProblemDefinition: "Building a SaaS startup versus staying employed",
			DecisionType:      "business/career",
			Options:           []string{"Build a SaaS startup", "Stay employed"},
			KeyFactors:        []string{"runway", "market validation", "income stability", "risk tolerance", "career leverage"},
			Risks:             []string{"startup failure", "burnout", "loss of stable income"},
			Unknowns:          []string{"customer demand", "time to revenue", "personal runway"},
			NextQuestions:     []string{"How many months of runway are available?", "Is there evidence of customer demand?", "Can the startup be tested part time first?"},
		}
	}

	body, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal mock payload: %w", err)
	}

	return string(body), nil
}

func classifyDecisionType(lower string) string {
	switch {
	case strings.Contains(lower, "job") || strings.Contains(lower, "career") || strings.Contains(lower, "employed"):
		return "career"
	case strings.Contains(lower, "master") || strings.Contains(lower, "degree") || strings.Contains(lower, "study"):
		return "education"
	case strings.Contains(lower, "startup") || strings.Contains(lower, "business") || strings.Contains(lower, "saas"):
		return "business"
	default:
		return "general"
	}
}

func repairPromptResponse(prompt string) string {
	start := strings.Index(prompt, "Original response:")
	if start == -1 {
		return "{}"
	}

	original := strings.TrimSpace(prompt[start+len("Original response:"):])
	original = strings.Trim(original, "`")
	original = strings.ReplaceAll(original, ",\n}", "\n}")
	original = strings.ReplaceAll(original, ",\n]", "\n]")
	original = strings.ReplaceAll(original, ",}", "}")
	original = strings.ReplaceAll(original, ",]", "]")
	return original
}
