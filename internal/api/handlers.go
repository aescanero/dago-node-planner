package api

import (
	"net/http"

	"github.com/aescanero/dago-node-planner/pkg/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// planHandler handles POST /api/v1/plan requests.
func (s *Server) planHandler(c *gin.Context) {
	var req models.PlanRequest

	// Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("invalid plan request",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
			"details": err.Error(),
		})
		return
	}

	s.logger.Info("received plan request",
		zap.String("task", req.Task),
	)

	// Execute planning
	resp, err := s.planner.Plan(c.Request.Context(), &req)
	if err != nil {
		s.logger.Error("planning failed",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "planning failed",
			"details": err.Error(),
		})
		return
	}

	s.logger.Info("plan request completed",
		zap.String("plan_id", resp.PlanID),
		zap.Int("iterations", resp.Iterations),
	)

	c.JSON(http.StatusOK, resp)
}

// validateHandler handles POST /api/v1/validate requests.
func (s *Server) validateHandler(c *gin.Context) {
	var req struct {
		GraphJSON string `json:"graph_json" binding:"required"`
	}

	// Parse request
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Warn("invalid validate request",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
			"details": err.Error(),
		})
		return
	}

	s.logger.Debug("received validate request")

	// Validate graph
	err := s.planner.ValidateGraph(req.GraphJSON)
	if err != nil {
		s.logger.Debug("validation failed",
			zap.Error(err),
		)
		c.JSON(http.StatusOK, models.ValidationResult{
			Valid:  false,
			Errors: []string{err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, models.ValidationResult{
		Valid: true,
	})
}
