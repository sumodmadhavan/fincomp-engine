package api

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	s := &Server{
		router: gin.Default(),
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.POST("/goalseek", s.GoalSeekHandler)
	s.router.POST("/runout", s.RunoutHandler)
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}
