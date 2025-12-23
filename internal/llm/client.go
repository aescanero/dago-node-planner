package llm

import (
	"context"
	"fmt"

	"github.com/aescanero/dago-libs/pkg/ports"
	"github.com/aescanero/dago-node-planner/internal/config"
	"go.uber.org/zap"
)

// Client wraps an LLM client with planner-specific functionality.
type Client struct {
	llmClient ports.LLMClient
	config    *config.LLMConfig
	retrier   *Retrier
	logger    *zap.Logger
	stats     *UsageStats
}

// NewClient creates a new LLM client for the planner.
func NewClient(llmClient ports.LLMClient, cfg *config.LLMConfig, logger *zap.Logger) *Client {
	return &Client{
		llmClient: llmClient,
		config:    cfg,
		retrier:   NewRetrier(cfg.RetryConfig),
		logger:    logger,
		stats:     &UsageStats{},
	}
}

// Complete sends a completion request to the LLM with retry logic.
func (c *Client) Complete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	c.logger.Debug("sending LLM completion request",
		zap.Int("max_tokens", req.MaxTokens),
		zap.Float64("temperature", req.Temperature),
	)

	var resp *CompletionResponse
	var err error

	// Execute with retry
	err = c.retrier.Do(ctx, func() error {
		resp, err = c.doComplete(ctx, req)
		return err
	})

	if err != nil {
		c.stats.AddCall(0, false)
		return nil, err
	}

	c.stats.AddCall(resp.TokensUsed, true)

	c.logger.Debug("received LLM response",
		zap.Int("tokens_used", resp.TokensUsed),
		zap.String("finish_reason", resp.FinishReason),
	)

	return resp, nil
}

// doComplete performs a single completion request without retry.
func (c *Client) doComplete(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	// Build messages
	messages := []ports.Message{}

	// Add system prompt as first message if provided
	if req.SystemPrompt != "" {
		messages = append(messages, ports.Message{
			Role:    "system",
			Content: req.SystemPrompt,
		})
	}

	// Add user prompt
	messages = append(messages, ports.Message{
		Role:    "user",
		Content: req.UserPrompt,
	})

	// Create LLM request
	llmReq := ports.CompletionRequest{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		Stop:        req.StopSequences,
	}

	// Call LLM
	llmResp, err := c.llmClient.Complete(ctx, llmReq)
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}

	// Extract response
	resp := &CompletionResponse{
		Content:      llmResp.Message.Content,
		Model:        llmResp.Model,
		TokensUsed:   llmResp.Usage.TotalTokens,
		FinishReason: llmResp.FinishReason,
	}

	return resp, nil
}

// GetStats returns current usage statistics.
func (c *Client) GetStats() *UsageStats {
	return c.stats
}

// ResetStats resets usage statistics.
func (c *Client) ResetStats() {
	c.stats = &UsageStats{}
}
