# LLM Decision Engine

![Go Version](https://img.shields.io/badge/Go-1.22-00ADD8?logo=go)
![Gin](https://img.shields.io/badge/API-Gin-008ECF)
![OpenAPI](https://img.shields.io/badge/OpenAPI-3.0-6BA539?logo=openapiinitiative)
![Ollama Ready](https://img.shields.io/badge/Ollama-ready-111111)
![Status](https://img.shields.io/badge/status-experimental-orange)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen)

A production-oriented Go service that converts ambiguous natural-language decision questions into structured decision models that downstream systems can validate, analyze, and extend.

## Why this project exists

Decision questions are usually messy, emotional, and incomplete. This project demonstrates how to turn that ambiguity into a clean JSON contract that an API, workflow engine, agent system, or UI can safely consume.

It is designed as a practical AI systems engineering example for:

- structured prompting
- backend orchestration
- provider abstraction
- schema validation
- API design
- reliability patterns

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Repository Structure](#repository-structure)
- [API](#api)
- [Running Locally](#running-locally)
- [Testing Strategy](#testing-strategy)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Security](#security)
- [Project Status](#project-status)

## Overview

Users often ask emotionally loaded and unstructured questions such as:

> Should I stop being a mobile developer and move into AI?

This project turns that type of input into a structured decision object with clear options, factors, risks, unknowns, and recommended follow-up questions.

## Features

- natural-language decision analysis over HTTP
- structured JSON output with JSON Schema validation
- OpenAI and Ollama provider support
- mock mode for local development
- Swagger UI and embedded OpenAPI spec
- retry and malformed JSON repair logic
- real LLM integration testing support
- modular Go project layout for extension

## Example Decision Object

```json
{
  "problem_definition": "Should I stop being a mobile developer and move into AI?",
  "decision_type": "Career Change",
  "options": ["Stay as a Mobile Developer", "Move into AI"],
  "key_factors": ["Skill Set", "Financial Gain", "Personal Interest", "Job Security"],
  "risks": ["Loss of Current Income", "Steep Learning Curve", "Uncertainty in Future Demand"],
  "unknowns": ["Current Demand for AI Talent", "Time and Effort Required to Adapt"],
  "recommended_next_questions": [
    "What are the current demand and salary ranges for AI developers in my location?",
    "How much time and effort will it take to adapt my skill set to AI?",
    "What are the potential career growth opportunities in AI?"
  ]
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
Schema Validation from decision_schema.json
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
│   │   ├── decision_parser.go
│   │   ├── prompt_builder.go
│   │   └── schema_validator.go
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
├── prompts
│   ├── decision_schema.json
│   └── system_prompt.txt
├── tests
│   ├── engine_test.go
│   └── real_llm_test.go
├── README.md
└── go.mod
```

## Decision Schema

Runtime validation uses [prompts/decision_schema.json](prompts/decision_schema.json) through `github.com/santhosh-tekuri/jsonschema/v5`.

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
  "question": "Should I stop being a mobile developer and move into AI?"
}
```

### Response

```json
{
  "decision": {
    "problem_definition": "Should I stop being a mobile developer and move into AI?",
    "decision_type": "Career Change",
    "options": ["Stay as a Mobile Developer", "Move into AI"],
    "key_factors": ["Skill Set", "Financial Gain", "Personal Interest", "Job Security"],
    "risks": ["Loss of Current Income", "Steep Learning Curve", "Uncertainty in Future Demand"],
    "unknowns": ["Current Demand for AI Talent", "Time and Effort Required to Adapt"],
    "recommended_next_questions": [
      "What are the current demand and salary ranges for AI developers in my location?",
      "How much time and effort will it take to adapt my skill set to AI?",
      "What are the potential career growth opportunities in AI?"
    ]
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

The user prompt is built directly from [prompts/decision_schema.json](prompts/decision_schema.json) and the incoming question, while shared behavior instructions live in [prompts/system_prompt.txt](prompts/system_prompt.txt).

## Reliability Strategy

This service is built to demonstrate production-minded LLM engineering:

- retry model calls up to 3 times
- enforce provider request timeouts
- parse fenced JSON and embedded JSON objects
- attempt lightweight JSON repair for malformed output
- fall back to a secondary repair prompt if needed
- validate responses against [prompts/decision_schema.json](prompts/decision_schema.json) with `santhosh-tekuri/jsonschema/v5` before responding

## LLM Client Behavior

The service supports three modes:

1. **OpenAI mode** for hosted API calls.
2. **Ollama mode** for local models.
3. **Mock mode** for deterministic local development and demos.

Environment variables are loaded into a typed config object in [internal/config/config.go](internal/config/config.go), then shared by the server and LLM client factory.

### Environment Variables

- `LLM_PROVIDER` - `openai`, `ollama`, or `mock`; default is `openai` when `OPENAI_API_KEY` is set, otherwise `mock`
- `PORT` - HTTP port, default `8080`
- `SCHEMA_PATH` - decision schema path, default `prompts/decision_schema.json`
- `SYSTEM_PROMPT_PATH` - shared system prompt for all LLM providers, default `prompts/system_prompt.txt`
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
      -d '{"question":"Should I stop being a mobile developer and move into AI?"}'
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

- Should I stop being a mobile developer and move into AI?

### Education

- Should I pursue a master's degree in AI?

### Business

- Should I build a SaaS startup or stay employed?

## Testing Strategy

### Unit Tests

- `TestPromptBuilder`
- `TestPromptBuilderBuildsSystemPrompt`
- `TestDecisionValidation`
- `TestJSONParsing`

### Integration-Style Flow Tests

- `TestEngineDecisionFlow`
- `TestEngineRepairsMalformedJSON`
- `TestRealLLMDecisionFlow` (skipped unless enabled with env)

## Roadmap

- provider abstraction for OpenAI and Anthropic
- decision scoring and ranking
- persistent decision history
- streaming partial reasoning metadata
- human-in-the-loop editing workflows
- decision graph generation with nodes and edges

## Contributing

Contributions are welcome.

Useful contribution areas:

- additional LLM providers
- richer validation and scoring
- authentication and rate limiting
- persistence and history support
- Docker and deployment tooling
- CI workflows and release automation

For pull requests, prefer:

- focused changes
- tests for new behavior
- updates to [README.md](README.md) or [internal/api/docs/openapi.yaml](internal/api/docs/openapi.yaml) when behavior changes

## Security

- do not commit real API keys or `.env` files
- prefer `.env.example` for documenting configuration
- treat LLM output as untrusted until validated

If you plan to publish this repository publicly, add a dedicated security policy and license file.

## Project Status

This project is currently in an experimental, portfolio-ready state. The core service works, but it is still a good candidate for hardening before production deployment.
