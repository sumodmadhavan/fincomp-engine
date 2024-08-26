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

	if err := params.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the existing Calculate function to preserve the response format
	result, err := goalseek.Calculate(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) RunoutHandler(c *gin.Context) {
	var params runout.RunoutParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := params.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use the existing Calculate function to preserve the response format
	result, err := runout.Calculate(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
