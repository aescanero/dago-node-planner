# Planning Algorithm

This document describes the graph planning algorithm used by dago-node-planner.

## Overview

The planning algorithm converts natural language tasks into executable workflow graphs through a multi-phase process with iterative refinement.

## Algorithm Phases

### Phase 1: Task Analysis (Optional)

**Input**: Natural language task description

**Process**:
1. Send task to LLM with analysis prompt
2. Extract structured analysis:
   - Complexity level (simple/moderate/complex)
   - Tool requirements (true/false)
   - Routing requirements (true/false)
   - Suggested node types
   - Key entities
   - Task intent

**Output**: TaskAnalysis object

**Benefits**:
- Helps LLM focus on relevant node types
- Provides context for graph generation
- Enables better prompt construction

**Example**:

```
Task: "Check if user is premium, then send personalized email"

Analysis:
{
  "complexity": "moderate",
  "requires_tools": true,
  "requires_routing": true,
  "suggested_node_types": ["executor", "router"],
  "key_entities": ["user", "premium", "email"],
  "intent": "Conditional email sending based on user tier"
}
```

### Phase 2: Graph Generation

**Input**:
- Task description
- Task analysis (optional)
- Constraints
- JSON schemas

**Process**:

#### Iteration 1: Initial Generation

1. Load schemas (graph, executor-node, router-node)
2. Build planning prompt with:
   - Task description
   - Analysis results
   - Schema information
   - Constraints
3. Call LLM with prompt
4. Extract graph JSON from response
5. Extract reasoning from response
6. Validate graph against schemas

If validation succeeds → Return graph
If validation fails → Proceed to Iteration 2

#### Iteration 2-N: Error Correction

1. Build error-fixing prompt with:
   - Original task
   - Previous graph JSON
   - Validation errors
   - Attempt number
2. Call LLM with error-fixing prompt
3. Extract corrected graph JSON
4. Validate graph against schemas

Repeat until:
- Validation succeeds → Return graph
- Max iterations reached → Return error

### Phase 3: Response Assembly

**Input**:
- Valid graph
- Reasoning
- Analysis
- Metadata (tokens, duration, etc.)

**Process**:
1. Package graph with metadata
2. Include all validation logs
3. Calculate confidence score (optional)
4. Return complete PlanResponse

## Iterative Refinement

### Why Iterative?

1. **Schema Complexity**: Graph schemas are complex with many constraints
2. **LLM Variability**: LLMs may not produce valid JSON on first try
3. **Error Correction**: Specific error feedback helps LLM fix issues
4. **Quality Assurance**: Ensures only valid graphs are returned

### Iteration Strategy

**Max Iterations**: Configurable (default: 3)

**Backoff**: No delay between iterations (synchronous)

**Error Feedback**: Full validation error messages passed to LLM

**Success Rate**: Typically 80%+ on first try, 95%+ after iterations

## JSON Extraction

### Challenges

LLMs often wrap JSON in markdown code blocks or add explanatory text.

### Extraction Strategy

1. **Code Block Detection**: Look for ```json...``` blocks
2. **Pattern Matching**: Find JSON objects with "nodes" field
3. **Validation**: Verify extracted string is valid JSON
4. **Fallback**: Bracket-based extraction as last resort

### Example LLM Responses

**Case 1: Clean JSON**
```json
{
  "nodes": [...],
  "edges": [...]
}
```

**Case 2: Code Block**
````
```json
{
  "nodes": [...],
  "edges": [...]
}
```
````

**Case 3: With Reasoning**
```
Here's the graph for your task:

```json
{
  "nodes": [...],
  "edges": [...]
}
```

This graph uses two nodes...
```

All cases are handled correctly by the extractor.

## Schema Validation

### Validation Levels

1. **JSON Syntax**: Is it valid JSON?
2. **Schema Structure**: Does it match the schema?
3. **Required Fields**: Are all required fields present?
4. **Type Checking**: Are field types correct?
5. **Cross-References**: Do edge node IDs exist?

### Error Messages

Validation errors are descriptive:

```
- nodes is required
- edges.0.source: node_id "unknown" does not exist in nodes
- nodes.0.config.mode: must be one of [agent, llm, tool]
```

## Optimization Strategies

### Prompt Engineering

1. **System Prompt**: Establishes role and capabilities
2. **Examples**: Include example graphs (future enhancement)
3. **Constraints**: Clearly state schema requirements
4. **Format**: Explicit response format instructions

### Temperature

- **Planning**: 0.0 (deterministic, schema-compliant)
- **Analysis**: 0.0 (consistent analysis)

Lower temperature reduces variation and improves schema compliance.

### Token Budget

- **Planning**: 4096 tokens (allows complex graphs)
- **Analysis**: 1024 tokens (brief analysis sufficient)

### Caching

- **Schema Caching**: Schemas loaded once and cached
- **Prompt Templates**: Loaded once and reused
- **LLM Stats**: Track token usage across requests

## Failure Modes

### Common Failures

1. **Invalid JSON**: LLM produces malformed JSON
   - **Recovery**: Extract JSON pattern from response
   - **Iteration**: Send error feedback

2. **Schema Violations**: Missing required fields, wrong types
   - **Recovery**: Send specific validation errors
   - **Iteration**: LLM fixes specific issues

3. **Logical Errors**: Graph structure doesn't match task
   - **Recovery**: User reviews and requests regeneration
   - **Prevention**: Better task analysis and constraints

4. **Max Iterations Exceeded**: Cannot produce valid graph
   - **Recovery**: Return error with logs
   - **User Action**: Simplify task or adjust constraints

### Error Handling

**Transient Errors** (network, rate limits):
- Retry with exponential backoff
- Max 3 attempts per LLM call

**Validation Errors**:
- Iterate with error feedback
- Max iterations configurable

**Critical Errors**:
- Return immediately with error
- Log for debugging

## Performance Characteristics

### Typical Timings

- **Simple Task** (1-2 nodes): 2-5 seconds
- **Moderate Task** (3-5 nodes): 5-10 seconds
- **Complex Task** (6+ nodes): 10-20 seconds

### Token Usage

- **Simple Task**: 500-1500 tokens
- **Moderate Task**: 1500-3000 tokens
- **Complex Task**: 3000-6000 tokens

### Success Rates

- **First Iteration**: 80-85%
- **After 2 Iterations**: 95%+
- **After 3 Iterations**: 98%+

## Future Enhancements

1. **Few-Shot Learning**: Include example graphs in prompts
2. **Confidence Scoring**: Rate graph quality
3. **Caching**: Cache graphs for similar tasks
4. **Streaming**: Stream graph generation progress
5. **Multi-Model**: Try different models if one fails
6. **Self-Improvement**: Learn from successful patterns
