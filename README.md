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

> Should I switch from mobile development to AI/LLM engineering?

This project turns that type of input into a structured decision object with clear options, factors, risks, unknowns, and recommended follow-up questions.

## Features

- natural-language decision analysis over HTTP
- structured JSON output with validation
- OpenAI and Ollama provider support
- mock mode for local development
- Swagger UI and embedded OpenAPI spec
- retry and malformed JSON repair logic
- real LLM integration testing support
- modular Go project layout for extension

## Example Decision Object

```json
{
  "problem_definition": "Switching from mobile development to AI/LLM engineering",
  "decision_type": "career",
  "options": ["Transition into AI/LLM engineering", "Stay in mobile development"],
  "key_factors": [
    "market demand",
    "learning curve",
    "portfolio readiness",
    "income stability",
    "long-term career growth"
  ],
  "risks": ["temporary productivity dip", "shallow domain knowledge", "slower-than-expected job transition"],
  "unknowns": ["time needed to become job-ready", "availability of AI/LLM roles", "salary impact during transition"],
  "recommended_next_questions": [
    "What AI/LLM skills are most required for target roles?",
    "How long will it take to build a credible portfolio?",
    "Can the transition start inside the current role?"
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
  "question": "Should I switch from mobile development to AI/LLM engineering?"
}
```

### Response

```json
{
  "decision": {
    "problem_definition": "Switching from mobile development to AI/LLM engineering",
    "decision_type": "career",
    "options": ["Transition into AI/LLM engineering", "Stay in mobile development"],
    "key_factors": [
      "market demand",
      "learning curve",
      "portfolio readiness",
      "income stability",
      "long-term career growth"
    ],
    "risks": ["temporary productivity dip", "shallow domain knowledge", "slower-than-expected job transition"],
    "unknowns": ["time needed to become job-ready", "availability of AI/LLM roles", "salary impact during transition"],
    "recommended_next_questions": [
      "What AI/LLM skills are most required for target roles?",
      "How long will it take to build a credible portfolio?",
      "Can the transition start inside the current role?"
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
      -d '{"question":"Should I switch from mobile development to AI/LLM engineering?"}'
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

- Should I switch from mobile development to AI/LLM engineering?

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
