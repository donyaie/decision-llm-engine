package model

// Decision is the normalized contract returned by the LLM decision engine.
type Decision struct {
	ProblemDefinition string   `json:"problem_definition"`
	DecisionType      string   `json:"decision_type"`
	Options           []string `json:"options"`
	KeyFactors        []string `json:"key_factors"`
	Risks             []string `json:"risks"`
	Unknowns          []string `json:"unknowns"`
	NextQuestions     []string `json:"recommended_next_questions"`
}

// AnalyzeRequest is the API input payload.
type AnalyzeRequest struct {
	Question string `json:"question"`
}

// AnalyzeResponse is the API output payload.
type AnalyzeResponse struct {
	Decision Decision `json:"decision"`
}
