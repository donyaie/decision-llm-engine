# LLM Decision Engine

A production-oriented Go service that converts ambiguous natural-language decision questions into structured decision models that downstream systems can validate, analyze, and extend.

## Overview

Users often ask emotionally loaded and unstructured questions such as:

> Should I move to Dubai for a software job?

This project turns that type of input into a structured decision object with clear options, factors, risks, unknowns, and recommended follow-up questions.

## Example Decision Object

```json
{
  "problem_definition": "Moving to Dubai for a software job",
  "decision_type": "life/career",
  "options": ["Move to Dubai", "Stay in current country"],
  "key_factors": ["salary", "cost of living", "career growth", "visa stability", "quality of life"],
  "risks": ["job instability", "high living costs", "cultural adaptation"],
  "unknowns": ["exact salary offer", "long term visa policy", "housing costs near the workplace"],
  "recommended_next_questions": ["What salary is offered?", "What is the visa duration?", "What is the rent cost?"]
}
```

## Architecture

```text
User Question
      |
      v
POST /v1/decision/analyze
      |
      v
Prompt Builder
      |
      v
LLM Client
      |
      v
Structured JSON Output
      |
      v
Schema Validation
      |
      v
Decision Object
      |
      v
API Response
```

## Repository Structure

```text
decision-llm-engine
│
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── api
│   │   ├── docs
│   │   │   └── openapi.yaml
│   │   ├── handler.go
│   │   └── swagger.go
│   ├── engine
│   │   ├── decision_engine.go
│   │   └── prompt_builder.go
│   ├── config
│   │   └── config.go
│   ├── llm
│   │   ├── client.go
│   │   ├── mock_client.go
│   │   ├── ollama.go
│   │   └── openai.go
│   ├── model
│   │   └── decision.go
│   ├── reliability
│   │   └── retry.go
│   └── validation
│       └── schema_validator.go
├── prompts
│   └── decision_prompt.txt
├── tests
│   ├── engine_test.go
│   └── real_llm_test.go
├── README.md
└── go.mod
```

## Decision Schema

```go
type Decision struct {
    ProblemDefinition string   `json:"problem_definition"`
    DecisionType      string   `json:"decision_type"`
    Options           []string `json:"options"`
    KeyFactors        []string `json:"key_factors"`
    Risks             []string `json:"risks"`
    Unknowns          []string `json:"unknowns"`
    NextQuestions     []string `json:"recommended_next_questions"`
}
```

## API

### Endpoint

`POST /v1/decision/analyze`

### Health Endpoint

`GET /health`

### Swagger

- Swagger UI: [http://localhost:8080/swagger/](http://localhost:8080/swagger/)
- OpenAPI spec: [http://localhost:8080/swagger/openapi.yaml](http://localhost:8080/swagger/openapi.yaml)

### Request

```json
{
  "question": "Should I move to Dubai for a software job?"
}
```

### Response

```json
{
  "decision": {
    "problem_definition": "Moving to Dubai for a software job",
    "decision_type": "life/career",
    "options": ["Move to Dubai", "Stay in current country"],
    "key_factors": ["salary", "cost of living", "career growth", "visa stability", "quality of life"],
    "risks": ["job instability", "high living costs", "cultural adaptation"],
    "unknowns": ["exact salary offer", "long term visa policy", "housing costs near the workplace"],
    "recommended_next_questions": ["What salary is offered?", "What is the visa duration?", "What is the rent cost?"]
  }
}
```

## Prompt Engineering Strategy

The prompt layer instructs the model to:

- extract the core decision
- identify realistic options
- list key decision factors
- surface risks and unknowns
- return only valid JSON matching the schema

The default template lives in [prompts/decision_prompt.txt](prompts/decision_prompt.txt).

## Reliability Strategy

This service is built to demonstrate production-minded LLM engineering:

- retry model calls up to 3 times
- enforce provider request timeouts
- parse fenced JSON and embedded JSON objects
- attempt lightweight JSON repair for malformed output
- fall back to a secondary repair prompt if needed
- validate required schema fields before responding

## LLM Client Behavior

The service supports three modes:

1. **OpenAI mode** for hosted API calls.
2. **Ollama mode** for local models.
3. **Mock mode** for deterministic local development and demos.

Environment variables are loaded into a typed config object in [internal/config/config.go](internal/config/config.go), then shared by the server and LLM client factory.

### Environment Variables

- `LLM_PROVIDER` - `openai`, `ollama`, or `mock`; default is `openai` when `OPENAI_API_KEY` is set, otherwise `mock`
- `PORT` - HTTP port, default `8080`
- `PROMPT_PATH` - prompt template path, default `prompts/decision_prompt.txt`
- `OPENAI_API_KEY` - enables live provider mode
- `OPENAI_BASE_URL` - optional override for OpenAI-compatible endpoints
- `OPENAI_MODEL` - optional model name override
- `OLLAMA_BASE_URL` - Ollama host, default `http://localhost:11434`
- `OLLAMA_MODEL` - Ollama model name, default `llama3.2`

The server reads these values from a local `.env` file automatically on startup.

`openapi` is also accepted as an alias for `openai` in `LLM_PROVIDER`.

## Running Locally

### Create a local env file

Copy [.env.example](.env.example) to `.env` and fill in your values.

### Start the API

```bash
go run ./cmd/server
```

Then open [http://localhost:8080/swagger/](http://localhost:8080/swagger/) to explore the API.

### Example Request

```bash
curl -X POST http://localhost:8080/v1/decision/analyze \
  -H "Content-Type: application/json" \
  -d '{"question":"Should I move to Dubai for a software job?"}'
```

### Run Tests

```bash
go test ./...
```

### Run Real LLM Integration Tests

Set `LLM_PROVIDER=openai` or `LLM_PROVIDER=ollama` in `.env`. Then run:

```bash
go test ./tests -run TestRealLLMDecisionFlow -v
```

## Demonstration Scenarios

### Career

- Should I move to Dubai for a software job?

### Education

- Should I pursue a master's degree in AI?

### Business

- Should I build a SaaS startup or stay employed?

## Testing Strategy

### Unit Tests

- `TestPromptBuilder`
- `TestDecisionValidation`
- `TestJSONParsing`

### Integration-Style Flow Tests

- `TestEngineDecisionFlow`
- `TestEngineRepairsMalformedJSON`
- `TestRealLLMDecisionFlow` (skipped unless enabled with env)

## Future Improvements

- provider abstraction for OpenAI and Anthropic
- decision scoring and ranking
- persistent decision history
- streaming partial reasoning metadata
- human-in-the-loop editing workflows
- decision graph generation with nodes and edges
