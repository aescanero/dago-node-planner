# Prompt Engineering Guide

This document provides guidance on prompt engineering for the dago-node-planner.

## Overview

The planner uses carefully crafted prompts to guide the LLM in generating valid workflow graphs. Effective prompts are critical for:

- Schema compliance
- Graph quality
- Validation success rate
- Token efficiency

## Prompt Structure

### System Prompt

The system prompt establishes the LLM's role and capabilities:

```
You are an expert graph planning assistant for the DA Orchestrator workflow system.

Your role is to convert natural language task descriptions into valid execution graphs.
```

**Key elements:**
- Role definition
- Domain context
- Response format requirements
- Schema awareness

### Planning Prompt

The planning prompt provides:
- Task description
- Task analysis (if available)
- Constraints
- Schema information
- Instructions

**Template variables:**
- `{{TASK}}`: The natural language task
- `{{ANALYSIS}}`: Task analysis results
- `{{CONSTRAINTS}}`: Planning constraints
- `{{SCHEMAS}}`: Schema information

### Error-Fixing Prompt

When validation fails, the error-fixing prompt includes:
- Original task
- Previous graph attempt
- Specific validation errors
- Iteration number

## Best Practices

### 1. Be Explicit About Schema

✅ Good:
```
The graph must have:
- At least one node
- Valid node types: "executor" or "router"
- All edge source/target IDs must reference existing nodes
```

❌ Bad:
```
Generate a valid graph
```

### 2. Provide Examples (Future Enhancement)

Include example graphs for common patterns:
```
Example simple graph:
{
  "nodes": [{"id": "task1", "type": "executor", ...}],
  "edges": [],
  "entry_point": "task1"
}
```

### 3. Use Structured Output Format

✅ Good:
```
Respond with a JSON object:
{
  "reasoning": "...",
  "graph": {...}
}
```

❌ Bad:
```
Provide a graph and explain your reasoning
```

### 4. Be Specific About Constraints

✅ Good:
```
Constraints:
- Maximum 10 nodes
- Prefer agent mode for complex tasks
- Available tools: send_email, create_ticket
```

❌ Bad:
```
Keep it simple
```

### 5. Leverage Task Analysis

Use analysis to focus the LLM:
```
Based on the analysis showing this is a complex task requiring tools and routing...
```

## Temperature Settings

| Prompt Type | Temperature | Rationale |
|-------------|-------------|-----------|
| Planning | 0.0 | Deterministic, schema-compliant |
| Analysis | 0.0 | Consistent analysis |
| Error Fixing | 0.0 | Precise corrections |

Lower temperature reduces variation and improves schema compliance.

## Token Budget

| Prompt Type | Max Tokens | Rationale |
|-------------|------------|-----------|
| Planning | 4096 | Complex graphs need space |
| Analysis | 1024 | Brief analysis sufficient |
| Error Fixing | 4096 | May need full graph regeneration |

## Error Feedback Strategy

### Specific vs Generic Feedback

✅ Good:
```
Validation errors:
1. nodes.0.config.mode: must be one of [agent, llm, tool]
2. edges.0.target: node "unknown_node" does not exist
```

❌ Bad:
```
The graph has validation errors. Please fix them.
```

### Iteration Context

Include iteration number to help LLM understand urgency:
```
This is attempt 3 of 3. Please carefully review all validation errors.
```

## Common Failure Patterns

### Pattern: Invalid JSON

**Symptom**: LLM wraps JSON in markdown or adds text

**Solution**: Explicitly request clean JSON
```
Respond with ONLY the JSON object, no markdown code blocks.
```

### Pattern: Missing Required Fields

**Symptom**: Schema validation fails for missing fields

**Solution**: List all required fields
```
Every executor node must have:
- id (string)
- type ("executor")
- config.mode (one of: agent, llm, tool)
```

### Pattern: Invalid Node References

**Symptom**: Edges reference non-existent nodes

**Solution**: Emphasize validation
```
Ensure all edge source/target IDs reference nodes that exist in the nodes array.
```

## Advanced Techniques

### Few-Shot Learning (Future)

Include 2-3 example task→graph pairs:
```
Example 1:
Task: "Send welcome email"
Graph: {...}

Example 2:
Task: "Analyze sentiment and route"
Graph: {...}

Now generate a graph for: {{TASK}}
```

### Chain-of-Thought

Encourage reasoning before generation:
```
First, analyze what the task requires:
1. What actions are needed?
2. What tools are required?
3. Is routing needed?

Then, generate the graph.
```

### Confidence Scoring (Future)

Ask LLM to rate its confidence:
```
After generating the graph, rate your confidence (0-1) that it will pass validation.
```

## Testing Prompts

### A/B Testing

Test prompt variations:
```bash
# Baseline
PLANNER_PROMPT_PATH=./prompts/baseline make test

# Variation A
PLANNER_PROMPT_PATH=./prompts/variation-a make test

# Compare results
```

### Metrics to Track

- Validation success rate (first iteration)
- Average iterations needed
- Token usage
- Graph quality (manual review)

## Customization

Customize prompts for specific domains:

```
# E-commerce domain
You are a workflow planner for e-commerce automation...
Common tasks include: order processing, inventory management, customer notifications...

# Support ticket domain
You are a workflow planner for support ticket routing...
Common patterns include: triage, categorization, assignment...
```

## Resources

- [OpenAI Prompt Engineering Guide](https://platform.openai.com/docs/guides/prompt-engineering)
- [Anthropic Prompt Design](https://docs.anthropic.com/claude/docs/prompt-design)
- [Few-Shot Learning](https://arxiv.org/abs/2005.14165)
