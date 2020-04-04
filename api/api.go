package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/victornm/es-backend/docs"
)

type Server struct {
	router http.Handler

	config *ServerConfig
}

func NewServer(config *ServerConfig) *Server {
	return &Server{config: config}
}

// @title ES API
// @version 1.0

// @contact.name VictorNM
// @contact.url https://github.com/VictorNM/

// @host localhost:8080
// @BasePath /api

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization
func (s *Server) Init() {
	router := gin.Default()

	router.Use(cors.Default()) // TODO: change this setting later

	rootAPI := router.Group("/api")

	// sign in
	authGroup := rootAPI.Group("/auth")
	{
		authGroup.POST("/sign-in", s.createSignInHandler())
	}

	// user profile
	rootAPI.GET("/profile", s.createAuthMiddleware(), s.createGetProfileHandler())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	s.router = router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ServerConfig struct {
	JWTSecret       string
	JWTExpiredHours int
}
