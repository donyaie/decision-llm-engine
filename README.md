# LLM Decision Engine

A production-oriented Go service that converts ambiguous natural-language decision questions into structured decision models that downstream systems can validate, analyze, and extend.

## Overview

Users often ask emotionally loaded and unstructured questions such as:

> Should I move to Dubai for a software job?

This project turns that type of input into a structured decision object with clear options, factors, risks, unknowns, and recommended follow-up questions.

## Core Idea

The engine receives a free-form question, sends a carefully built prompt to an LLM, validates the returned JSON, and responds with a normalized decision model.

Example decision object:

```json
{
  "problem_definition": "Moving to Dubai for a software job",
  "decision_type": "life/career",
  "options": ["Move to Dubai", "Stay in current country"],
  "key_factors": ["salary", "cost of living", "career growth", "visa stability", "quality of life"],
  "risks": ["job instability", "high living costs", "cultural adaptation"],
  "unknowns": ["exact salary offer", "long term visa policy"],
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

## Planned Repository Structure

```text
decision-llm-engine
│
├── cmd
│   └── server
│       └── main.go
│
├── internal
│   ├── api
│   │   └── handler.go
│   ├── engine
│   │   ├── decision_engine.go
│   │   └── prompt_builder.go
│   ├── llm
│   │   └── client.go
│   ├── model
│   │   └── decision.go
│   ├── validation
│   │   └── schema_validator.go
│   └── reliability
│       └── retry.go
│
├── prompts
│   └── decision_prompt.txt
├── tests
│   └── engine_test.go
├── README.md
└── go.mod
```

## API Design

### Endpoint

`POST /v1/decision/analyze`

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
    "unknowns": ["exact salary offer", "long term visa policy"],
    "recommended_next_questions": ["What salary is offered?", "What is the visa duration?", "What is the rent cost?"]
  }
}
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

## Prompt Engineering Strategy

The prompt layer will instruct the model to:

- extract the core decision
- identify realistic options
- list key decision factors
- surface risks and unknowns
- return only valid JSON matching the schema

Example template:

```text
You are a decision analysis assistant.

Your job is to convert a user question into a structured decision model.

Rules:
- Extract the core decision
- Identify options
- Identify key factors
- Identify risks
- Identify unknown information

Return ONLY valid JSON.
```

## System Components

### 1. API Layer

Receives requests, validates input, and returns the structured decision response.

### 2. Engine Layer

Coordinates prompt building, LLM calls, parsing, repair, and validation.

### 3. LLM Client Layer

Handles provider communication, authentication, timeouts, and retries.

### 4. Validation Layer

Ensures the returned object matches the expected contract and contains required fields.

### 5. Reliability Layer

Adds retry logic, timeout control, and JSON repair for malformed model output.

## Reliability Strategy

This project is intended to demonstrate production-minded LLM engineering.

### Retry

- retry transient LLM failures up to 3 times
- apply bounded backoff between attempts

### Timeout

- default LLM timeout target: 15 seconds
- cancel requests that exceed the response budget

### JSON Repair

- detect malformed JSON
- send repair prompt when possible
- re-validate repaired output before returning it

### Schema Validation

- reject empty `problem_definition`
- reject outputs with no extracted `options`
- return structured errors for invalid model output

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

### Integration Tests

- `TestEngineDecisionFlow`

## Future Improvements

- provider abstraction for OpenAI and Anthropic
- configurable prompt templates
- scoring and ranking of options
- persistent decision history
- streaming responses
- human feedback loop
- decision graph generation with nodes and edges

## Why This Project Matters

This repository is designed to demonstrate:

- AI engineering
- prompt engineering
- LLM orchestration
- API design
- system reliability
- schema validation
- Go backend engineering
