// Package planner provides the core graph planning service.
package planner

import (
	"context"
	"fmt"
	"time"

	"github.com/aescanero/dago-libs/pkg/schema"
	"github.com/aescanero/dago-node-planner/internal/config"
	"github.com/aescanero/dago-node-planner/internal/llm"
	"github.com/aescanero/dago-node-planner/pkg/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service is the main planning service that orchestrates graph generation.
type Service struct {
	llmClient       *llm.Client
	schemaValidator *schema.Validator
	analyzer        *Analyzer
	generator       *Generator
	config          *config.PlanningConfig
	logger          *zap.Logger
}

// NewService creates a new planning service.
func NewService(
	llmClient *llm.Client,
	schemaValidator *schema.Validator,
	cfg *config.PlanningConfig,
	logger *zap.Logger,
) *Service {
	analyzer := NewAnalyzer(llmClient, logger)

	prompter := NewPrompter(cfg.PromptPath, logger)
	extractor := NewExtractor(logger)
	iterator := NewIterator(cfg.MaxIterations, logger)

	generator := NewGenerator(
		llmClient,
		prompter,
		extractor,
		schemaValidator,
		iterator,
		logger,
	)

	return &Service{
		llmClient:       llmClient,
		schemaValidator: schemaValidator,
		analyzer:        analyzer,
		generator:       generator,
		config:          cfg,
		logger:          logger,
	}
}

// Plan generates a graph from a natural language task.
func (s *Service) Plan(ctx context.Context, req *models.PlanRequest) (*models.PlanResponse, error) {
	startTime := time.Now()
	planID := uuid.New().String()

	s.logger.Info("starting graph planning",
		zap.String("plan_id", planID),
		zap.String("task", req.Task),
	)

	// Step 1: Analyze task (optional)
	var analysis *models.TaskAnalysis
	if !req.SkipAnalysis && s.config.EnableAnalysis {
		var err error
		analysis, err = s.analyzer.Analyze(ctx, req.Task)
		if err != nil {
			s.logger.Warn("task analysis failed, continuing without it",
				zap.Error(err),
			)
		} else {
			s.logger.Debug("task analysis completed",
				zap.String("complexity", string(analysis.Complexity)),
				zap.Bool("requires_tools", analysis.RequiresTools),
				zap.Bool("requires_routing", analysis.RequiresRouting),
			)
		}
	}

	// Step 2: Generate graph
	genReq := &GenerateRequest{
		Task:        req.Task,
		Context:     req.Context,
		Analysis:    analysis,
		Constraints: req.Constraints,
	}

	genResp, err := s.generator.Generate(ctx, genReq)
	if err != nil {
		return nil, fmt.Errorf("graph generation failed: %w", err)
	}

	// Step 3: Build response
	duration := time.Since(startTime)
	stats := s.llmClient.GetStats()

	resp := &models.PlanResponse{
		PlanID:         planID,
		Graph:          genResp.Graph,
		GraphJSON:      genResp.GraphJSON,
		Reasoning:      genResp.Reasoning,
		Analysis:       analysis,
		Iterations:     genResp.Iterations,
		ValidationLogs: genResp.ValidationLogs,
		Metadata: &models.PlanMetadata{
			LLMProvider:     "anthropic", // TODO: get from config
			LLMModel:        "unknown",   // TODO: get from llm client
			TokensUsed:      stats.TotalTokens,
			Duration:        duration,
			ConfidenceScore: 0.0, // TODO: implement confidence scoring
			Success:         true,
		},
		CreatedAt: time.Now(),
	}

	s.logger.Info("graph planning completed",
		zap.String("plan_id", planID),
		zap.Int("iterations", genResp.Iterations),
		zap.Int("tokens_used", stats.TotalTokens),
		zap.Duration("duration", duration),
	)

	return resp, nil
}

// ValidateGraph validates a graph JSON string.
func (s *Service) ValidateGraph(graphJSON string) error {
	return s.schemaValidator.ValidateGraph([]byte(graphJSON))
}
