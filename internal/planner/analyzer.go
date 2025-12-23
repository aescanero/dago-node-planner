package planner

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aescanero/dago-node-planner/internal/llm"
	"github.com/aescanero/dago-node-planner/pkg/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Analyzer analyzes tasks before graph generation.
type Analyzer struct {
	llmClient *llm.Client
	logger    *zap.Logger
}

// NewAnalyzer creates a new task analyzer.
func NewAnalyzer(llmClient *llm.Client, logger *zap.Logger) *Analyzer {
	return &Analyzer{
		llmClient: llmClient,
		logger:    logger,
	}
}

// Analyze analyzes a task to understand its requirements.
func (a *Analyzer) Analyze(ctx context.Context, task string) (*models.TaskAnalysis, error) {
	a.logger.Debug("analyzing task", zap.String("task", task))

	// Build analysis prompt
	systemPrompt := a.buildSystemPrompt()
	userPrompt := a.buildUserPrompt(task)

	// Call LLM
	resp, err := a.llmClient.Complete(ctx, &llm.CompletionRequest{
		SystemPrompt: systemPrompt,
		UserPrompt:   userPrompt,
		MaxTokens:    1024,
		Temperature:  0.0,
	})
	if err != nil {
		return nil, fmt.Errorf("LLM analysis failed: %w", err)
	}

	// Parse response
	analysis, err := a.parseAnalysis(resp.Content, task)
	if err != nil {
		return nil, fmt.Errorf("failed to parse analysis: %w", err)
	}

	return analysis, nil
}

// buildSystemPrompt builds the system prompt for task analysis.
func (a *Analyzer) buildSystemPrompt() string {
	return `You are a task analysis expert for a graph-based workflow orchestration system.

Your role is to analyze natural language task descriptions and extract:
1. Task complexity (simple, moderate, complex)
2. Whether the task requires external tools
3. Whether the task requires conditional routing/branching
4. Suggested node types (executor nodes, router nodes)
5. Key entities mentioned in the task
6. Overall intent of the task

Respond with a JSON object in this exact format:
{
  "complexity": "simple|moderate|complex",
  "requires_tools": true|false,
  "requires_routing": true|false,
  "suggested_node_types": ["executor", "router"],
  "key_entities": ["entity1", "entity2"],
  "intent": "Brief description of task intent",
  "reasoning": "Explanation of your analysis"
}

Be concise and accurate. Focus on extracting actionable insights for graph planning.`
}

// buildUserPrompt builds the user prompt for task analysis.
func (a *Analyzer) buildUserPrompt(task string) string {
	return fmt.Sprintf("Analyze this task:\n\n%s", task)
}

// parseAnalysis parses the LLM response into a TaskAnalysis.
func (a *Analyzer) parseAnalysis(content, task string) (*models.TaskAnalysis, error) {
	// Extract JSON from response
	jsonStr := extractJSON(content)
	if jsonStr == "" {
		return nil, fmt.Errorf("no JSON found in response")
	}

	// Parse JSON
	var result struct {
		Complexity         string   `json:"complexity"`
		RequiresTools      bool     `json:"requires_tools"`
		RequiresRouting    bool     `json:"requires_routing"`
		SuggestedNodeTypes []string `json:"suggested_node_types"`
		KeyEntities        []string `json:"key_entities"`
		Intent             string   `json:"intent"`
		Reasoning          string   `json:"reasoning"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert to TaskAnalysis
	analysis := &models.TaskAnalysis{
		TaskID:             uuid.New().String(),
		Complexity:         models.ComplexityLevel(result.Complexity),
		RequiresTools:      result.RequiresTools,
		RequiresRouting:    result.RequiresRouting,
		SuggestedNodeTypes: result.SuggestedNodeTypes,
		KeyEntities:        result.KeyEntities,
		Intent:             result.Intent,
		Reasoning:          result.Reasoning,
		AnalyzedAt:         time.Now(),
	}

	return analysis, nil
}
