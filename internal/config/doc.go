// Package config provides configuration management for the dago-node-planner service.
//
// Configuration can be loaded from YAML files and overridden with environment variables.
// All environment variables use the PLANNER_ prefix.
//
// Example configuration file (config.yaml):
//
//	server:
//	  host: "0.0.0.0"
//	  port: 8080
//	  read_timeout: 30s
//	  write_timeout: 30s
//
//	llm:
//	  provider: "anthropic"
//	  api_key: "your-api-key"
//	  model: "claude-3-5-sonnet-20241022"
//	  max_tokens: 4096
//	  temperature: 0.0
//	  timeout: 60s
//	  retry:
//	    max_attempts: 3
//	    initial_delay: 1s
//	    max_delay: 10s
//	    multiplier: 2.0
//
//	planning:
//	  max_iterations: 3
//	  max_nodes: 50
//	  prompt_path: "./prompts"
//	  enable_validation: true
//	  enable_analysis: true
//	  confidence_threshold: 0.8
//
//	logging:
//	  level: "info"
//	  format: "json"
//	  output_path: "stdout"
//
// Environment variables:
//   - PLANNER_SERVER_HOST: Server host address
//   - PLANNER_SERVER_PORT: Server port number
//   - PLANNER_LLM_PROVIDER: LLM provider (anthropic, openai)
//   - PLANNER_LLM_API_KEY: LLM API key
//   - PLANNER_LLM_MODEL: LLM model name
//   - PLANNER_LLM_MAX_TOKENS: Maximum tokens for LLM responses
//   - PLANNER_LLM_TEMPERATURE: LLM temperature (0.0-1.0)
//   - PLANNER_MAX_ITERATIONS: Maximum planning iterations
//   - PLANNER_MAX_NODES: Maximum nodes per graph
//   - PLANNER_PROMPT_PATH: Path to prompt template files
//   - PLANNER_LOG_LEVEL: Logging level (debug, info, warn, error)
//   - PLANNER_LOG_FORMAT: Logging format (json, console)
//
// Note: JSON schemas are embedded in dago-libs and loaded automatically.
package config
