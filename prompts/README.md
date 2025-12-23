# LLM Prompt Templates

This directory contains prompt templates used by the dago-node-planner for graph generation.

## Files

- **system-prompt.txt**: The system prompt that establishes the LLM's role and capabilities
- **task-planning.txt**: Template for the initial graph planning request
- **error-fixing.txt**: Template for iterative error correction

## Placeholders

Templates support the following placeholders:

### task-planning.txt
- `{{TASK}}`: The natural language task description
- `{{ANALYSIS}}`: Task analysis results (optional)
- `{{CONSTRAINTS}}`: Planning constraints (optional)
- `{{SCHEMAS}}`: JSON schema information

### error-fixing.txt
- `{{TASK}}`: The original task description
- `{{PREVIOUS_GRAPH}}`: The graph JSON from the previous attempt
- `{{VALIDATION_ERRORS}}`: List of validation errors
- `{{ATTEMPT}}`: Current attempt number

## Customization

You can customize these prompts by:

1. Editing the files directly
2. Setting `PLANNER_PROMPT_PATH` environment variable to a custom directory
3. Providing a custom path in the configuration file

The planner will fall back to built-in default prompts if custom files are not found.

## Tips for Effective Prompts

1. **Be Specific**: Clearly define what constitutes a valid graph
2. **Provide Examples**: Include example graph structures in the system prompt
3. **Set Expectations**: Specify the response format explicitly
4. **Include Context**: Provide schema information and node type details
5. **Iterate**: Test prompts with various task types and refine based on results

## Testing Prompts

To test prompt changes:

```bash
# Run the planner with your custom prompts
PLANNER_PROMPT_PATH=./prompts ./node-planner

# Or specify in config.yaml
planning:
  prompt_path: "./prompts"
```

Monitor the generated graphs and adjust prompts to improve:
- Graph correctness
- Efficiency (fewer nodes when appropriate)
- Validation success rate
- Reasoning quality
