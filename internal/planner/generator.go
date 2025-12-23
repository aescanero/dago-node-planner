package planner

import (
	"context"
	"fmt"

	"github.com/aescanero/dago-libs/pkg/schema"
	"github.com/aescanero/dago-node-planner/internal/llm"
	"github.com/aescanero/dago-node-planner/pkg/models"
	"go.uber.org/zap"
)

// GenerateRequest represents a request to generate a graph.
type GenerateRequest struct {
	Task        string
	Context     map[string]any
	Analysis    *models.TaskAnalysis
	Constraints *models.Constraints
}

// GenerateResponse represents the result of graph generation.
type GenerateResponse struct {
	Graph          any      // The generated graph (parsed JSON)
	GraphJSON      string   // Raw JSON string
	Reasoning      string   // LLM's reasoning
	Iterations     int      // Number of iterations performed
	ValidationLogs []string // Validation logs from each iteration
}

// Generator orchestrates graph generation with iterative refinement.
type Generator struct {
	llmClient       *llm.Client
	prompter        *Prompter
	extractor       *Extractor
	schemaValidator *schema.Validator
	iterator        *Iterator
	logger          *zap.Logger
}

// NewGenerator creates a new graph generator.
func NewGenerator(
	llmClient *llm.Client,
	prompter *Prompter,
	extractor *Extractor,
	schemaValidator *schema.Validator,
	iterator *Iterator,
	logger *zap.Logger,
) *Generator {
	return &Generator{
		llmClient:       llmClient,
		prompter:        prompter,
		extractor:       extractor,
		schemaValidator: schemaValidator,
		iterator:        iterator,
		logger:          logger,
	}
}

// Generate generates a graph from a task with iterative refinement.
func (g *Generator) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	g.logger.Debug("starting graph generation")

	// Get schemas for prompt
	schemas, err := g.getSchemas()
	if err != nil {
		return nil, fmt.Errorf("failed to get schemas: %w", err)
	}

	// Initial generation
	prompt, err := g.prompter.BuildPlanningPrompt(req.Task, req.Analysis, schemas, req.Constraints)
	if err != nil {
		return nil, fmt.Errorf("failed to build planning prompt: %w", err)
	}

	var graphJSON string
	var reasoning string
	var validationLogs []string
	iteration := 0

	// Iterative refinement loop
	err = g.iterator.Iterate(ctx, func(ctx context.Context, attempt int) error {
		iteration = attempt

		var llmResp *llm.CompletionResponse
		var llmErr error

		if attempt == 1 {
			// First attempt: use planning prompt
			llmResp, llmErr = g.llmClient.Complete(ctx, &llm.CompletionRequest{
				SystemPrompt: g.prompter.GetSystemPrompt(),
				UserPrompt:   prompt,
				MaxTokens:    4096,
				Temperature:  0.0,
			})
		} else {
			// Subsequent attempts: use error-fixing prompt
			validationErrs := validationLogs[len(validationLogs)-1]
			fixPrompt, err := g.prompter.BuildErrorFixingPrompt(req.Task, graphJSON, []string{validationErrs}, attempt)
			if err != nil {
				return fmt.Errorf("failed to build error-fixing prompt: %w", err)
			}

			llmResp, llmErr = g.llmClient.Complete(ctx, &llm.CompletionRequest{
				SystemPrompt: g.prompter.GetSystemPrompt(),
				UserPrompt:   fixPrompt,
				MaxTokens:    4096,
				Temperature:  0.0,
			})
		}

		if llmErr != nil {
			return fmt.Errorf("LLM request failed: %w", llmErr)
		}

		// Extract graph JSON
		extractedJSON, extractedReasoning, err := g.extractor.Extract(llmResp.Content)
		if err != nil {
			validationLogs = append(validationLogs, fmt.Sprintf("Extraction error: %s", err))
			return err
		}

		graphJSON = extractedJSON
		reasoning = extractedReasoning

		// Validate graph
		if err := g.schemaValidator.ValidateGraph([]byte(graphJSON)); err != nil {
			errMsg := fmt.Sprintf("Validation failed: %s", err.Error())
			validationLogs = append(validationLogs, errMsg)
			return err
		}

		// Success!
		validationLogs = append(validationLogs, "Validation successful")
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("graph generation failed after %d iterations: %w", iteration, err)
	}

	// Parse final graph
	graph, err := g.extractor.ParseGraph(graphJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to parse final graph: %w", err)
	}

	resp := &GenerateResponse{
		Graph:          graph,
		GraphJSON:      graphJSON,
		Reasoning:      reasoning,
		Iterations:     iteration,
		ValidationLogs: validationLogs,
	}

	g.logger.Debug("graph generation completed",
		zap.Int("iterations", iteration),
	)

	return resp, nil
}

// getSchemas returns schema information for prompts.
// Note: Full schemas are embedded in dago-libs and used for validation.
// For prompts, we provide a simplified description.
func (g *Generator) getSchemas() (map[string]string, error) {
	schemas := map[string]string{
		"graph.schema.json": "Graph schema with nodes, edges, and entry_point",
		"executor-node.schema.json": "Executor node schema with modes: agent, llm, tool",
		"router-node.schema.json": "Router node schema with modes: deterministic, llm, hybrid",
	}

	return schemas, nil
}
