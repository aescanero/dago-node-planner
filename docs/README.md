# DA Node Planner Documentation

Welcome to the DA Node Planner documentation.

## Overview

The DA Node Planner is an LLM-powered service that converts natural language task descriptions into executable workflow graphs for the DA Orchestrator. It combines intelligent task analysis, schema-aware generation, and iterative validation to produce reliable graph definitions.

## Table of Contents

- [Architecture](#architecture)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Planning Algorithm](#planning-algorithm)
- [Prompt Engineering](#prompt-engineering)
- [Examples](#examples)

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                    dago-node-planner                        │
│                                                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────┐  │
│  │   REST API   │  │  WebSocket   │  │   gRPC (future)  │  │
│  └──────┬───────┘  └──────┬───────┘  └──────┬───────────┘  │
│         │                 │                  │              │
│         └─────────────────┴──────────────────┘              │
│                           │                                 │
│                           ▼                                 │
│                  ┌─────────────────┐                        │
│                  │ Planning Service│                        │
│                  └────────┬────────┘                        │
│                           │                                 │
│         ┌─────────────────┼─────────────────┐               │
│         │                 │                 │               │
│         ▼                 ▼                 ▼               │
│  ┌────────────┐   ┌──────────────┐  ┌────────────┐         │
│  │  Analyzer  │   │  Generator   │  │ Validator  │         │
│  └────────────┘   └──────────────┘  └────────────┘         │
│         │                 │                 │               │
│         └─────────────────┴─────────────────┘               │
│                           │                                 │
│         ┌─────────────────┼─────────────────┐               │
│         │                 │                 │               │
│         ▼                 ▼                 ▼               │
│  ┌────────────┐   ┌──────────────┐  ┌────────────┐         │
│  │LLM Client  │   │Schema Repo   │  │ Prompter   │         │
│  └────────────┘   └──────────────┘  └────────────┘         │
└─────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

**Planning Service**: Main orchestrator that coordinates the planning process

**Analyzer**: Pre-analyzes tasks to understand complexity, tools needed, and routing requirements

**Generator**: Generates graphs using LLM with iterative refinement

**Validator**: Validates graphs against JSON schemas

**LLM Client**: Wraps LLM providers (Anthropic, OpenAI) with retry logic

**Schema Repository**: Loads and caches JSON schemas for validation

**Prompter**: Builds LLM prompts from templates

### Data Flow

1. Client sends task description via REST API
2. Planning Service receives request
3. (Optional) Analyzer analyzes the task
4. Generator builds planning prompt with schemas
5. Generator calls LLM to generate graph
6. Generator extracts JSON from LLM response
7. Validator validates graph against schemas
8. If validation fails, Generator iterates with error feedback
9. Planning Service returns validated graph to client

## Configuration

### Configuration File

Create a `config.yaml` file:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s

llm:
  provider: "anthropic"
  api_key: "your-api-key"  # Or use env var
  model: "claude-3-5-sonnet-20241022"
  max_tokens: 4096
  temperature: 0.0
  timeout: 60s
  retry:
    max_attempts: 3
    initial_delay: 1s
    max_delay: 10s
    multiplier: 2.0

planning:
  max_iterations: 3
  max_nodes: 50
  prompt_path: "./prompts"
  # Note: JSON schemas are embedded in dago-libs
  enable_validation: true
  enable_analysis: true
  confidence_threshold: 0.8

logging:
  level: "info"
  format: "json"
  output_path: "stdout"
```

### Environment Variables

All configuration can be overridden with environment variables:

```bash
# Server
export PLANNER_SERVER_HOST=0.0.0.0
export PLANNER_SERVER_PORT=8080

# LLM
export PLANNER_LLM_PROVIDER=anthropic
export PLANNER_LLM_API_KEY=your-api-key
export PLANNER_LLM_MODEL=claude-3-5-sonnet-20241022
export PLANNER_LLM_MAX_TOKENS=4096
export PLANNER_LLM_TEMPERATURE=0.0

# Planning
export PLANNER_MAX_ITERATIONS=3
export PLANNER_MAX_NODES=50
export PLANNER_PROMPT_PATH=./prompts

# Note: JSON schemas are embedded in dago-libs

# Logging
export PLANNER_LOG_LEVEL=info
export PLANNER_LOG_FORMAT=json
```

## API Reference

### POST /api/v1/plan

Generate a graph from a task description.

**Request:**

```json
{
  "task": "Task description in natural language",
  "context": {
    "key": "value"
  },
  "constraints": {
    "max_nodes": 10,
    "preferred_modes": ["agent", "llm"],
    "available_tools": ["tool1", "tool2"],
    "max_iterations": 3
  },
  "skip_analysis": false
}
```

**Response:**

```json
{
  "plan_id": "uuid",
  "graph": { ... },
  "graph_json": "...",
  "reasoning": "...",
  "analysis": { ... },
  "iterations": 2,
  "validation_logs": [...],
  "metadata": {
    "llm_provider": "anthropic",
    "llm_model": "claude-3-5-sonnet-20241022",
    "tokens_used": 1234,
    "duration": "2.5s",
    "success": true
  },
  "created_at": "2024-01-01T12:00:00Z"
}
```

### POST /api/v1/validate

Validate a graph JSON.

**Request:**

```json
{
  "graph_json": "{ ... }"
}
```

**Response:**

```json
{
  "valid": true,
  "errors": [],
  "warnings": []
}
```

### GET /health

Health check endpoint.

**Response:**

```json
{
  "status": "healthy",
  "time": "2024-01-01T12:00:00Z"
}
```

## Planning Algorithm

See [PLANNING.md](PLANNING.md) for detailed algorithm documentation.

## Prompt Engineering

See [PROMPTS.md](PROMPTS.md) for prompt engineering guidelines.

## Examples

See [EXAMPLES.md](EXAMPLES.md) for example tasks and generated graphs.
