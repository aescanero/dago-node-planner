# DA Node Planner (dago-node-planner)

[![CI](https://github.com/aescanero/dago-node-planner/actions/workflows/ci.yml/badge.svg)](https://github.com/aescanero/dago-node-planner/actions/workflows/ci.yml)
[![Docker](https://github.com/aescanero/dago-node-planner/actions/workflows/docker.yml/badge.svg)](https://github.com/aescanero/dago-node-planner/actions/workflows/docker.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aescanero/dago-node-planner)](https://goreportcard.com/report/github.com/aescanero/dago-node-planner)
[![License](https://img.shields.io/github/license/aescanero/dago-node-planner)](LICENSE)

**LLM-powered graph planner for the DA Orchestrator workflow system.**

The DA Node Planner converts natural language task descriptions into executable workflow graphs for the DA Orchestrator. It uses Large Language Models (LLMs) to understand tasks and automatically generate optimal node configurations with iterative validation.

## Features

- **Natural Language Input**: Describe tasks in plain English
- **Intelligent Graph Generation**: LLM-powered graph design with schema awareness
- **Iterative Refinement**: Automatic error correction through validation feedback
- **Schema Validation**: Ensures generated graphs conform to DA Orchestrator schemas
- **Task Analysis**: Pre-analyzes tasks for complexity and requirements
- **REST API**: HTTP endpoints for integration
- **Modular Design**: Can be used standalone or embedded in orchestrator

## Quick Start

### Prerequisites

- Go 1.25.5 or later
- An Anthropic API key (or OpenAI key if configured)

### Installation

```bash
# Clone the repository
git clone https://github.com/aescanero/dago-node-planner.git
cd dago-node-planner

# Install dependencies
make deps

# Build
make build
```

**Note**: JSON schemas are embedded in `dago-libs` and loaded automatically. No additional configuration needed.

### Running Locally

```bash
# Set your API key
export PLANNER_LLM_API_KEY=your-anthropic-api-key

# Run the server
./scripts/run-local.sh

# Or use make
make run
```

The server will start on `http://localhost:8080`.

### Using Docker

```bash
docker run -p 8080:8080 \
  -e PLANNER_LLM_API_KEY=your-api-key \
  aescanero/dago-node-planner:latest
```

## Usage

### REST API

Generate a graph from a task:

```bash
curl -X POST http://localhost:8080/api/v1/plan \
  -H "Content-Type: application/json" \
  -d '{
    "task": "Analyze user sentiment and route to appropriate handler based on sentiment score",
    "constraints": {
      "max_nodes": 10,
      "preferred_modes": ["agent", "llm"]
    }
  }'
```

Response:

```json
{
  "plan_id": "uuid",
  "graph": {
    "nodes": [...],
    "edges": [...],
    "entry_point": "sentiment_analyzer"
  },
  "graph_json": "...",
  "reasoning": "This graph uses an executor node in agent mode to analyze sentiment, followed by a router node to direct the flow based on sentiment score...",
  "iterations": 1,
  "metadata": {
    "llm_provider": "anthropic",
    "llm_model": "claude-3-5-sonnet-20241022",
    "tokens_used": 1234,
    "duration": "2.5s",
    "success": true
  }
}
```

Validate a graph:

```bash
curl -X POST http://localhost:8080/api/v1/validate \
  -H "Content-Type: application/json" \
  -d '{
    "graph_json": "{\"nodes\": [...], \"edges\": [...]}"
  }'
```

## Configuration

Configuration can be provided via YAML file or environment variables.

### Environment Variables

- `PLANNER_LLM_API_KEY`: LLM provider API key (required)
- `PLANNER_LLM_PROVIDER`: LLM provider (default: "anthropic")
- `PLANNER_LLM_MODEL`: LLM model name (default: "claude-3-5-sonnet-20241022")
- `PLANNER_SERVER_PORT`: HTTP server port (default: 8080)
- `PLANNER_MAX_ITERATIONS`: Max refinement iterations (default: 3)
- `PLANNER_MAX_NODES`: Max nodes per graph (default: 50)
- `PLANNER_PROMPT_PATH`: Path to prompt templates (default: "./prompts")
- `PLANNER_LOG_LEVEL`: Log level (default: "info")

**Note**: JSON schemas are embedded in dago-libs and loaded automatically.

See [Configuration Documentation](docs/README.md#configuration) for full details.

## Architecture

```
┌─────────────────────────────────────────┐
│         Natural Language Task           │
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│          Task Analyzer                  │
│   - Complexity estimation               │
│   - Tool/routing detection              │
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│        Graph Generator                  │
│   ┌──────────────────────────────────┐  │
│   │  1. Build planning prompt        │  │
│   │  2. Call LLM with schemas        │  │
│   │  3. Extract graph JSON           │  │
│   │  4. Validate against schemas     │  │
│   │  5. If invalid → iterate         │  │
│   └──────────────────────────────────┘  │
└─────────────┬───────────────────────────┘
              │
              ▼
┌─────────────────────────────────────────┐
│          Valid Graph JSON               │
└─────────────────────────────────────────┘
```

See [Architecture Documentation](docs/README.md#architecture) for detailed design.

## Development

### Running Tests

```bash
make test
```

### Running Linter

```bash
make lint
```

### Building Docker Image

```bash
make docker-build
```

### Live Reload Development

Requires [air](https://github.com/cosmtrek/air):

```bash
make dev
```

## Documentation

- [Architecture Overview](docs/README.md)
- [Planning Algorithm](docs/PLANNING.md)
- [Prompt Engineering](docs/PROMPTS.md)
- [Examples](docs/EXAMPLES.md)

## Dependencies

- [dago-libs](https://github.com/aescanero/dago-libs) v0.2.0: Domain models and ports
- [dago-adapters](https://github.com/aescanero/dago-adapters) v0.2.0: LLM adapters

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Authors

- **Antonio Escanero** - [aescanero](https://github.com/aescanero)

## Acknowledgments

- DA Orchestrator project
- Anthropic for Claude API
- The Go community
