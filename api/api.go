package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	router http.Handler

	config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{config: config}
}

func (s *Server) Init() {
	router := gin.Default()

	router.Use(cors.Default()) // TODO: change this setting later

	rootAPI := router.Group("/api")

	// sign in
	authGroup := rootAPI.Group("/auth")
	authGroup.POST("/sign-in", s.createSignInHandler())

	// user profile
	rootAPI.GET("/profile", s.createAuthMiddleware(), s.createGetProfileHandler())

	s.router = router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ServerConfig struct {
	JWTSecret       string
	JWTExpiredHours int
}
