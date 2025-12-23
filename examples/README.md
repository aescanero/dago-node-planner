# Examples

This directory contains example tasks and their expected graph outputs.

## Files

- **simple-task.json**: Basic single-node graph for sending an email
- **complex-task.json**: Multi-node graph with sentiment analysis, categorization, and routing

## Structure

Each example file contains:

```json
{
  "request": {
    "task": "Natural language task description",
    "context": { ... },
    "constraints": { ... }
  },
  "expected_graph": {
    "nodes": [ ... ],
    "edges": [ ... ],
    "entry_point": "node_id"
  },
  "description": "Explanation of the graph design"
}
```

## Using Examples

### Test the Planner

Send an example task to the planner:

```bash
curl -X POST http://localhost:8080/api/v1/plan \
  -H "Content-Type: application/json" \
  -d @examples/simple-task.json
```

### Validate Expected Output

Validate the expected graph:

```bash
curl -X POST http://localhost:8080/api/v1/validate \
  -H "Content-Type: application/json" \
  -d "$(jq -c '{graph_json: (.expected_graph | tostring)}' examples/simple-task.json)"
```

## Example Categories

### Simple Tasks

Single executor node, no routing:
- Send email
- Call API
- Update database

### Moderate Tasks

Multiple executor nodes with sequential flow:
- Fetch data, transform, store
- Analyze text, extract entities, summarize

### Complex Tasks

Multi-node graphs with routing and branching:
- Sentiment analysis with conditional routing
- Multi-step workflows with error handling
- Parallel processing with aggregation

## Creating New Examples

1. Write a clear task description
2. Define any context data
3. Specify constraints (optional)
4. Design the expected graph manually
5. Test with the planner
6. Compare planner output to expected output
7. Document any differences

Example template:

```json
{
  "request": {
    "task": "Your task here",
    "context": {},
    "constraints": {
      "max_nodes": 10
    }
  },
  "expected_graph": {
    "nodes": [],
    "edges": [],
    "entry_point": "start"
  },
  "description": "What this graph does"
}
```

## Testing

Use these examples for:

- **Unit Tests**: Validate planner logic
- **Integration Tests**: End-to-end planning workflow
- **Regression Tests**: Ensure consistent output
- **Documentation**: Show what the planner can do
- **Benchmarking**: Measure performance on different task types

## Contributing Examples

When adding new examples:

1. Ensure the expected graph is valid (passes schema validation)
2. Test with the actual planner
3. Document any interesting design decisions
4. Keep tasks realistic and useful
5. Include context data when relevant
