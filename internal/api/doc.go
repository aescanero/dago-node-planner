// Package api provides HTTP REST API for the dago-node-planner service.
//
// The API server exposes endpoints for graph planning and validation.
//
// Endpoints:
//
//   GET /health - Health check endpoint
//   GET /ready  - Readiness check endpoint
//
//   POST /api/v1/plan     - Generate a graph from a task
//   POST /api/v1/validate - Validate a graph JSON
//
// Example plan request:
//
//	POST /api/v1/plan
//	{
//	  "task": "Analyze user sentiment and route to appropriate handler",
//	  "context": {
//	    "user_id": "12345"
//	  },
//	  "constraints": {
//	    "max_nodes": 10,
//	    "preferred_modes": ["agent", "llm"]
//	  }
//	}
//
// Example plan response:
//
//	{
//	  "plan_id": "uuid",
//	  "graph": { ... },
//	  "graph_json": "{ ... }",
//	  "reasoning": "...",
//	  "iterations": 2,
//	  "metadata": {
//	    "llm_provider": "anthropic",
//	    "llm_model": "claude-3-5-sonnet-20241022",
//	    "tokens_used": 1234,
//	    "duration": "2.5s",
//	    "success": true
//	  },
//	  "created_at": "2024-01-01T12:00:00Z"
//	}
//
// Example validate request:
//
//	POST /api/v1/validate
//	{
//	  "graph_json": "{ ... }"
//	}
//
// Example validate response:
//
//	{
//	  "valid": true
//	}
//
// Or if validation fails:
//
//	{
//	  "valid": false,
//	  "errors": [
//	    "nodes is required",
//	    "edges is required"
//	  ]
//	}
package api
