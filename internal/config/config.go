// Package config provides configuration management for the node planner service.
package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the complete configuration for the node planner service.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	LLM      LLMConfig      `yaml:"llm"`
	Planning PlanningConfig `yaml:"planning"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// ServerConfig contains HTTP server configuration.
type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// LLMConfig contains LLM provider configuration.
type LLMConfig struct {
	Provider    string        `yaml:"provider"` // anthropic, openai
	APIKey      string        `yaml:"api_key"`
	Model       string        `yaml:"model"`
	MaxTokens   int           `yaml:"max_tokens"`
	Temperature float64       `yaml:"temperature"`
	Timeout     time.Duration `yaml:"timeout"`
	RetryConfig RetryConfig   `yaml:"retry"`
}

// RetryConfig contains retry configuration for LLM calls.
type RetryConfig struct {
	MaxAttempts  int           `yaml:"max_attempts"`
	InitialDelay time.Duration `yaml:"initial_delay"`
	MaxDelay     time.Duration `yaml:"max_delay"`
	Multiplier   float64       `yaml:"multiplier"`
}

// PlanningConfig contains planning-specific configuration.
type PlanningConfig struct {
	MaxIterations       int     `yaml:"max_iterations"`
	MaxNodes            int     `yaml:"max_nodes"`
	PromptPath          string  `yaml:"prompt_path"`
	EnableValidation    bool    `yaml:"enable_validation"`
	EnableAnalysis      bool    `yaml:"enable_analysis"`
	ConfidenceThreshold float64 `yaml:"confidence_threshold"`
}

// LoggingConfig contains logging configuration.
type LoggingConfig struct {
	Level      string `yaml:"level"` // debug, info, warn, error
	Format     string `yaml:"format"` // json, console
	OutputPath string `yaml:"output_path"`
}

// Load loads configuration from a YAML file and environment variables.
// Environment variables override YAML values.
func Load(configPath string) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host:         "0.0.0.0",
			Port:         8080,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		LLM: LLMConfig{
			Provider:    "anthropic",
			Model:       "claude-3-5-sonnet-20241022",
			MaxTokens:   4096,
			Temperature: 0.0,
			Timeout:     60 * time.Second,
			RetryConfig: RetryConfig{
				MaxAttempts:  3,
				InitialDelay: 1 * time.Second,
				MaxDelay:     10 * time.Second,
				Multiplier:   2.0,
			},
		},
		Planning: PlanningConfig{
			MaxIterations:       3,
			MaxNodes:            50,
			PromptPath:          "./prompts",
			EnableValidation:    true,
			EnableAnalysis:      true,
			ConfidenceThreshold: 0.8,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "json",
			OutputPath: "stdout",
		},
	}

	// Load from YAML file if provided
	if configPath != "" {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(data, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	cfg.applyEnvVars()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// applyEnvVars overrides configuration with environment variables.
func (c *Config) applyEnvVars() {
	if v := os.Getenv("PLANNER_SERVER_HOST"); v != "" {
		c.Server.Host = v
	}
	if v := os.Getenv("PLANNER_SERVER_PORT"); v != "" {
		_, _ = fmt.Sscanf(v, "%d", &c.Server.Port)
	}

	if v := os.Getenv("PLANNER_LLM_PROVIDER"); v != "" {
		c.LLM.Provider = v
	}
	if v := os.Getenv("PLANNER_LLM_API_KEY"); v != "" {
		c.LLM.APIKey = v
	}
	if v := os.Getenv("PLANNER_LLM_MODEL"); v != "" {
		c.LLM.Model = v
	}
	if v := os.Getenv("PLANNER_LLM_MAX_TOKENS"); v != "" {
		_, _ = fmt.Sscanf(v, "%d", &c.LLM.MaxTokens)
	}
	if v := os.Getenv("PLANNER_LLM_TEMPERATURE"); v != "" {
		_, _ = fmt.Sscanf(v, "%f", &c.LLM.Temperature)
	}

	if v := os.Getenv("PLANNER_MAX_ITERATIONS"); v != "" {
		_, _ = fmt.Sscanf(v, "%d", &c.Planning.MaxIterations)
	}
	if v := os.Getenv("PLANNER_MAX_NODES"); v != "" {
		_, _ = fmt.Sscanf(v, "%d", &c.Planning.MaxNodes)
	}
	if v := os.Getenv("PLANNER_PROMPT_PATH"); v != "" {
		c.Planning.PromptPath = v
	}

	if v := os.Getenv("PLANNER_LOG_LEVEL"); v != "" {
		c.Logging.Level = v
	}
	if v := os.Getenv("PLANNER_LOG_FORMAT"); v != "" {
		c.Logging.Format = v
	}
}

// Validate checks if the configuration is valid.
func (c *Config) Validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.LLM.Provider == "" {
		return fmt.Errorf("LLM provider is required")
	}
	if c.LLM.APIKey == "" {
		return fmt.Errorf("LLM API key is required")
	}
	if c.LLM.Model == "" {
		return fmt.Errorf("LLM model is required")
	}

	if c.Planning.MaxIterations <= 0 {
		return fmt.Errorf("max iterations must be positive")
	}
	if c.Planning.MaxNodes <= 0 {
		return fmt.Errorf("max nodes must be positive")
	}

	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}
