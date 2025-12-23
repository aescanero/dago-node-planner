#!/bin/bash
# Run dago-node-planner locally for development

set -e

# Default configuration
export PLANNER_SERVER_HOST="${PLANNER_SERVER_HOST:-0.0.0.0}"
export PLANNER_SERVER_PORT="${PLANNER_SERVER_PORT:-8080}"
export PLANNER_LLM_PROVIDER="${PLANNER_LLM_PROVIDER:-anthropic}"
export PLANNER_LLM_MODEL="${PLANNER_LLM_MODEL:-claude-3-5-sonnet-20241022}"
export PLANNER_MAX_ITERATIONS="${PLANNER_MAX_ITERATIONS:-3}"
export PLANNER_MAX_NODES="${PLANNER_MAX_NODES:-50}"
export PLANNER_PROMPT_PATH="${PLANNER_PROMPT_PATH:-./prompts}"
export PLANNER_LOG_LEVEL="${PLANNER_LOG_LEVEL:-info}"
export PLANNER_LOG_FORMAT="${PLANNER_LOG_FORMAT:-console}"

# Check for API key
if [ -z "$PLANNER_LLM_API_KEY" ]; then
    echo "ERROR: PLANNER_LLM_API_KEY environment variable is required"
    echo ""
    echo "Set it with:"
    echo "  export PLANNER_LLM_API_KEY=your-api-key"
    echo ""
    exit 1
fi

echo "Starting dago-node-planner..."
echo "Server: ${PLANNER_SERVER_HOST}:${PLANNER_SERVER_PORT}"
echo "LLM Provider: ${PLANNER_LLM_PROVIDER}"
echo "LLM Model: ${PLANNER_LLM_MODEL}"
echo ""

# Build if binary doesn't exist
if [ ! -f "./bin/node-planner" ]; then
    echo "Binary not found. Building..."
    ./scripts/build.sh
    echo ""
fi

# Run the binary
exec ./bin/node-planner "$@"
