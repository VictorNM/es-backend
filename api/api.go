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
	router *gin.Engine

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
	s.router = gin.Default()

	s.initRouter()
}

func (s *Server) initRouter() {
	s.router.Use(cors.Default()) // TODO: change this setting later

	rootAPI := s.router.Group("/api")

	// testing purpose: ping => pong
	rootAPI.GET("/ping", createPingHandler())

	// ===== user =====
	userGroup := rootAPI.Group("/users")
	{
		// not auth handlers
		userGroup.POST("/sign-in", s.createSignInHandler())
		userGroup.POST("/register", s.createRegisterHandler())
	}
	{
		// auth handler
		userGroup.GET("/profile", s.createAuthMiddleware(), s.createGetProfileHandler())
	}

	// swagger API documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ServerConfig struct {
	FrontendBaseURL string // the domain where the frontend live

	JWTSecret       string
	JWTExpiredHours int
}

// @Summary PING PONG
// @Description For testing
// @Tags ping
// @Produce json
// @Success 200 {object} api.BaseResponse "PING PONG"
// @Router /api/ping [get]
func createPingHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		response(c, 200, "PONG")
	}
}