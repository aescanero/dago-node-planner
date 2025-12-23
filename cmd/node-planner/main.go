// Package main provides the entry point for the dago-node-planner service.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aescanero/dago-libs/pkg/ports"
	"github.com/aescanero/dago-libs/pkg/schema"
	"github.com/aescanero/dago-adapters/pkg/llm/anthropic"
	"github.com/aescanero/dago-node-planner/internal/api"
	"github.com/aescanero/dago-node-planner/internal/config"
	"github.com/aescanero/dago-node-planner/internal/llm"
	"github.com/aescanero/dago-node-planner/internal/planner"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	configPath = flag.String("config", "", "Path to configuration file")
	version    = "dev" // Set by build process
)

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := initLogger(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting dago-node-planner",
		zap.String("version", version),
	)

	// Initialize schema validator (from dago-libs)
	schemaValidator, err := schema.NewValidator()
	if err != nil {
		logger.Fatal("failed to initialize schema validator",
			zap.Error(err),
		)
	}

	// Initialize LLM client
	llmClient, err := initLLMClient(cfg.LLM, logger)
	if err != nil {
		logger.Fatal("failed to initialize LLM client",
			zap.Error(err),
		)
	}

	// Initialize planner service
	plannerService := planner.NewService(
		llmClient,
		schemaValidator,
		&cfg.Planning,
		logger,
	)

	// Initialize API server
	apiServer := api.NewServer(
		&cfg.Server,
		plannerService,
		logger,
	)

	// Start server in background
	go func() {
		if err := apiServer.Start(); err != nil {
			logger.Fatal("server error",
				zap.Error(err),
			)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	logger.Info("received shutdown signal")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := apiServer.Stop(ctx); err != nil {
		logger.Error("error during shutdown",
			zap.Error(err),
		)
		os.Exit(1)
	}

	logger.Info("shutdown complete")
}

// initLogger initializes the logger based on configuration.
func initLogger(cfg config.LoggingConfig) (*zap.Logger, error) {
	var zapCfg zap.Config

	if cfg.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}

	// Set log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	// Set output path
	if cfg.OutputPath == "stdout" || cfg.OutputPath == "" {
		zapCfg.OutputPaths = []string{"stdout"}
	} else {
		zapCfg.OutputPaths = []string{cfg.OutputPath}
	}

	return zapCfg.Build()
}

// initLLMClient initializes the LLM client based on configuration.
func initLLMClient(cfg config.LLMConfig, logger *zap.Logger) (*llm.Client, error) {
	var llmClient ports.LLMClient
	var err error

	switch cfg.Provider {
	case "anthropic":
		llmClient, err = anthropic.NewClient(cfg.APIKey, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create Anthropic client: %w", err)
		}

	case "openai":
		// TODO: Implement OpenAI client when available
		return nil, fmt.Errorf("OpenAI provider not yet implemented")

	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", cfg.Provider)
	}

	return llm.NewClient(llmClient, &cfg, logger), nil
}
