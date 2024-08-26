// File: api/handlers.go

package api

import (
	"financialapi/internal/financials"
	"financialapi/internal/goalseek"
	"financialapi/internal/runout"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) GoalSeekHandler(c *gin.Context) {
	var params financials.FinancialParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engine := goalseek.NewGoalSeekCalculator(params)

	if err := engine.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := engine.Compute(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := engine.GetResult()
	c.JSON(http.StatusOK, result)
}

func (s *Server) RunoutHandler(c *gin.Context) {
	var params runout.RunoutParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engine := runout.NewRunoutCalculator(params)

	if err := engine.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := engine.Compute(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := engine.GetResult()
	c.JSON(http.StatusOK, result)
}