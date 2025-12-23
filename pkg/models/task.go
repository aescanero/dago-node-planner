// Package models defines data structures for the node planner service.
package models

import "time"

// Task represents a natural language task to be converted into a graph.
type Task struct {
	// ID is a unique identifier for the task
	ID string `json:"id"`

	// Description is the natural language description of the task
	Description string `json:"description"`

	// Context provides additional context for the task
	Context map[string]any `json:"context,omitempty"`

	// Constraints specifies constraints for graph generation
	Constraints *Constraints `json:"constraints,omitempty"`

	// CreatedAt is the timestamp when the task was created
	CreatedAt time.Time `json:"created_at"`
}

// Constraints defines constraints for graph generation.
type Constraints struct {
	// MaxNodes limits the maximum number of nodes in the graph
	MaxNodes int `json:"max_nodes,omitempty"`

	// PreferredModes specifies preferred execution modes (agent, llm, tool)
	PreferredModes []string `json:"preferred_modes,omitempty"`

	// AvailableTools lists tools available for use in the graph
	AvailableTools []string `json:"available_tools,omitempty"`

	// MaxIterations limits planning refinement iterations
	MaxIterations int `json:"max_iterations,omitempty"`

	// RequireValidation ensures the graph is validated before returning
	RequireValidation bool `json:"require_validation,omitempty"`
}

// TaskAnalysis contains the results of task analysis.
type TaskAnalysis struct {
	// TaskID is the ID of the analyzed task
	TaskID string `json:"task_id"`

	// Complexity estimates the task complexity
	Complexity ComplexityLevel `json:"complexity"`

	// RequiresTools indicates if tools are needed
	RequiresTools bool `json:"requires_tools"`

	// RequiresRouting indicates if routing logic is needed
	RequiresRouting bool `json:"requires_routing"`

	// SuggestedNodeTypes lists suggested node types
	SuggestedNodeTypes []string `json:"suggested_node_types"`

	// KeyEntities extracts key entities from the task
	KeyEntities []string `json:"key_entities"`

	// Intent describes the task intent
	Intent string `json:"intent"`

	// Reasoning provides the LLM's reasoning
	Reasoning string `json:"reasoning,omitempty"`

	// AnalyzedAt is the timestamp of analysis
	AnalyzedAt time.Time `json:"analyzed_at"`
}

// ComplexityLevel represents task complexity.
type ComplexityLevel string

const (
	// ComplexitySimple indicates a simple task
	ComplexitySimple ComplexityLevel = "simple"

	// ComplexityModerate indicates a moderate complexity task
	ComplexityModerate ComplexityLevel = "moderate"

	// ComplexityComplex indicates a complex task
	ComplexityComplex ComplexityLevel = "complex"
)
