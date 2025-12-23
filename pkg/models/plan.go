package models

import "time"

// PlanRequest represents a request to generate a graph plan.
type PlanRequest struct {
	// Task is the natural language task description
	Task string `json:"task" binding:"required"`

	// Context provides additional context for planning
	Context map[string]any `json:"context,omitempty"`

	// Constraints specifies planning constraints
	Constraints *Constraints `json:"constraints,omitempty"`

	// SkipAnalysis skips the task analysis phase
	SkipAnalysis bool `json:"skip_analysis,omitempty"`
}

// PlanResponse represents the result of graph planning.
type PlanResponse struct {
	// PlanID is a unique identifier for this plan
	PlanID string `json:"plan_id"`

	// Graph is the generated graph definition
	Graph any `json:"graph"` // Using any to avoid circular dependency with dago-libs

	// GraphJSON is the raw JSON representation of the graph
	GraphJSON string `json:"graph_json,omitempty"`

	// Reasoning explains the LLM's reasoning for the graph design
	Reasoning string `json:"reasoning"`

	// Analysis contains the task analysis (if performed)
	Analysis *TaskAnalysis `json:"analysis,omitempty"`

	// Iterations is the number of refinement iterations performed
	Iterations int `json:"iterations"`

	// ValidationLogs contains validation messages from each iteration
	ValidationLogs []string `json:"validation_logs,omitempty"`

	// Metadata contains metadata about the planning process
	Metadata *PlanMetadata `json:"metadata"`

	// CreatedAt is the timestamp when the plan was created
	CreatedAt time.Time `json:"created_at"`
}

// PlanMetadata contains metadata about the planning process.
type PlanMetadata struct {
	// LLMProvider is the LLM provider used (e.g., "anthropic")
	LLMProvider string `json:"llm_provider"`

	// LLMModel is the specific model used
	LLMModel string `json:"llm_model"`

	// TokensUsed is the total tokens consumed
	TokensUsed int `json:"tokens_used"`

	// Duration is the total planning duration
	Duration time.Duration `json:"duration"`

	// ConfidenceScore is an optional confidence score (0.0-1.0)
	ConfidenceScore float64 `json:"confidence_score,omitempty"`

	// Success indicates if planning succeeded
	Success bool `json:"success"`

	// ErrorMessage contains error details if planning failed
	ErrorMessage string `json:"error_message,omitempty"`
}

// ValidationResult represents the result of graph validation.
type ValidationResult struct {
	// Valid indicates if the graph is valid
	Valid bool `json:"valid"`

	// Errors contains validation error messages
	Errors []string `json:"errors,omitempty"`

	// Warnings contains validation warnings
	Warnings []string `json:"warnings,omitempty"`

	// SchemaVersion is the version of the schema used for validation
	SchemaVersion string `json:"schema_version,omitempty"`
}

// IterationLog represents a log entry for a planning iteration.
type IterationLog struct {
	// Iteration is the iteration number (1-based)
	Iteration int `json:"iteration"`

	// Prompt is the prompt sent to the LLM
	Prompt string `json:"prompt,omitempty"`

	// Response is the LLM's response
	Response string `json:"response,omitempty"`

	// GraphJSON is the extracted graph JSON
	GraphJSON string `json:"graph_json,omitempty"`

	// Validation is the validation result
	Validation *ValidationResult `json:"validation,omitempty"`

	// TokensUsed is tokens consumed in this iteration
	TokensUsed int `json:"tokens_used"`

	// Duration is the iteration duration
	Duration time.Duration `json:"duration"`

	// Timestamp is when this iteration occurred
	Timestamp time.Time `json:"timestamp"`
}
