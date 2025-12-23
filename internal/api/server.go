// Package api provides HTTP API for the node planner service.
package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aescanero/dago-node-planner/internal/config"
	"github.com/aescanero/dago-node-planner/internal/planner"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server is the HTTP API server.
type Server struct {
	config  *config.ServerConfig
	planner *planner.Service
	router  *gin.Engine
	logger  *zap.Logger
	server  *http.Server
}

// NewServer creates a new API server.
func NewServer(
	cfg *config.ServerConfig,
	plannerService *planner.Service,
	logger *zap.Logger,
) *Server {
	// Set Gin mode based on log level
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	s := &Server{
		config:  cfg,
		planner: plannerService,
		router:  router,
		logger:  logger,
	}

	// Register routes
	s.registerRoutes()

	return s
}

// registerRoutes registers all API routes.
func (s *Server) registerRoutes() {
	// Health check
	s.router.GET("/health", s.healthHandler)
	s.router.GET("/ready", s.readyHandler)

	// API v1
	v1 := s.router.Group("/api/v1")
	{
		// Planning endpoints
		v1.POST("/plan", s.planHandler)
		v1.POST("/validate", s.validateHandler)
	}
}

// Start starts the HTTP server.
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	s.server = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
	}

	s.logger.Info("starting API server",
		zap.String("address", addr),
	)

	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

// Stop gracefully stops the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping API server")

	if s.server != nil {
		return s.server.Shutdown(ctx)
	}

	return nil
}

// healthHandler handles health check requests.
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now().UTC(),
	})
}

// readyHandler handles readiness check requests.
func (s *Server) readyHandler(c *gin.Context) {
	// TODO: Check dependencies (LLM client, schema repo, etc.)
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"time":   time.Now().UTC(),
	})
}
