// Package llm provides LLM integration for the dago-node-planner.
//
// This package wraps the LLM client from dago-adapters and provides:
//   - Retry logic with exponential backoff
//   - Request/response models specific to planning
//   - Usage statistics tracking
//   - Planner-specific error handling
//
// The Client type is the main entry point for making LLM requests
// with automatic retries and error handling.
//
// Example usage:
//
//	llmClient, err := anthropic.NewClient(apiKey, "claude-3-5-sonnet-20241022")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	client := llm.NewClient(llmClient, cfg.LLM, logger)
//
//	resp, err := client.Complete(ctx, &llm.CompletionRequest{
//	    SystemPrompt: "You are a graph planner...",
//	    UserPrompt: "Generate a graph for: " + task,
//	    MaxTokens: 4096,
//	    Temperature: 0.0,
//	})
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	fmt.Println("Generated:", resp.Content)
//	fmt.Println("Tokens used:", resp.TokensUsed)
package llm
