// Package llm provides LLM integration for the node planner.
package llm

// CompletionRequest represents a request to the LLM.
type CompletionRequest struct {
	// SystemPrompt is the system message
	SystemPrompt string

	// UserPrompt is the user message
	UserPrompt string

	// MaxTokens is the maximum tokens to generate
	MaxTokens int

	// Temperature controls randomness (0.0-1.0)
	Temperature float64

	// StopSequences are sequences that stop generation
	StopSequences []string
}

// CompletionResponse represents a response from the LLM.
type CompletionResponse struct {
	// Content is the generated text
	Content string

	// Model is the model that generated the response
	Model string

	// TokensUsed is the total tokens consumed
	TokensUsed int

	// FinishReason indicates why generation stopped
	FinishReason string

	// Error contains error details if the request failed
	Error error
}

// UsageStats tracks token usage across multiple LLM calls.
type UsageStats struct {
	// TotalTokens is the total tokens used
	TotalTokens int

	// TotalCalls is the number of LLM calls made
	TotalCalls int

	// SuccessfulCalls is the number of successful calls
	SuccessfulCalls int

	// FailedCalls is the number of failed calls
	FailedCalls int
}

// AddCall updates usage statistics with a new call.
func (u *UsageStats) AddCall(tokensUsed int, success bool) {
	u.TotalCalls++
	u.TotalTokens += tokensUsed

	if success {
		u.SuccessfulCalls++
	} else {
		u.FailedCalls++
	}
}
