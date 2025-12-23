// Package models provides data structures for the dago-node-planner service.
//
// This package defines the core data models used throughout the planner,
// including task representations, planning requests/responses, and metadata.
//
// Key types:
//   - Task: Represents a natural language task to be converted into a graph
//   - PlanRequest: Request to generate a graph from a task
//   - PlanResponse: Response containing the generated graph and metadata
//   - TaskAnalysis: Results of task analysis before planning
//   - PlanMetadata: Metadata about the planning process
//   - ValidationResult: Result of graph schema validation
//   - IterationLog: Log of a single planning iteration
//
// The models in this package are designed to be serializable to/from JSON
// for use in REST APIs and other interfaces.
package models
