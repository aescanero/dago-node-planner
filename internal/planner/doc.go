// Package planner provides the core graph planning service for dago-node-planner.
//
// This package implements the main planning algorithm that converts natural
// language task descriptions into executable graph definitions.
//
// Architecture:
//
//   Task → Analyzer → Generator → Validator → Graph
//             ↓           ↓
//         Analysis    Prompter
//                     Extractor
//                     Iterator
//
// Components:
//
//   - Service: Main entry point orchestrating the planning process
//   - Analyzer: Analyzes tasks to understand requirements
//   - Generator: Generates graphs with iterative refinement
//   - Prompter: Builds LLM prompts from templates
//   - Extractor: Extracts JSON from LLM responses
//   - Iterator: Manages the iterative refinement loop
//
// The planning process:
//
// 1. Task Analysis (optional):
//    - Understand task complexity
//    - Identify if tools/routing needed
//    - Extract key entities and intent
//
// 2. Graph Generation:
//    - Build planning prompt with schemas
//    - Call LLM to generate graph
//    - Extract JSON from response
//    - Validate against schemas
//    - If validation fails, iterate with error feedback
//    - Return valid graph or error after max iterations
//
// 3. Response Assembly:
//    - Package graph with metadata
//    - Include reasoning and analysis
//    - Track tokens and timing
//
// Example usage:
//
//	service := planner.NewService(llmClient, schemaRepo, cfg, logger)
//
//	resp, err := service.Plan(ctx, &models.PlanRequest{
//	    Task: "Analyze user sentiment and route to appropriate handler",
//	    Constraints: &models.Constraints{
//	        MaxNodes: 10,
//	        PreferredModes: []string{"agent", "llm"},
//	    },
//	})
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Printf("Generated graph with %d iterations\n", resp.Iterations)
//	fmt.Printf("Graph: %s\n", resp.GraphJSON)
package planner
