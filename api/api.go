package api

import (
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/victornm/es-backend/docs"
)

type Server struct {
	router *gin.Engine
	db     *sqlx.DB

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
	s.connectDB()

	s.router = gin.Default()
	s.initRouter()
}

func (s *Server) initRouter() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true
	s.router.Use(cors.New(corsConfig)) // TODO: change this setting later

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

func (s *Server) connectDB() {
	u, err := url.Parse(s.config.SqlConnString)
	if err != nil {
		log.Fatalf("invalid SQL URL %v", err)
	}

	log.Println(u.Scheme)

	db, err := sqlx.Open(u.Scheme, s.config.SqlConnString)
	if err != nil {
		log.Fatalf("open database failed: %v", err)
	}

	s.db = db
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type ServerConfig struct {
	FrontendBaseURL string // the domain where the frontend live

	JWTSecret       string
	JWTExpiredHours int

	SqlConnString string
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
