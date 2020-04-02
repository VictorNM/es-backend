package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	router *gin.Engine

	config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{config: config}
}

func (s *Server) Init() {
	s.router = gin.Default()
	rootAPI := s.router.Group("/api")

	// auth
	authGroup := rootAPI.Group("/auth")
	authGroup.POST("/sign-in", s.createSignInHandler())
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ServerConfig struct {
	JWTSecret       string
	JWTExpiredHours int
}
