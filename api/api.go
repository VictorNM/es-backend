package api

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "github.com/victornm/es-backend/docs"
)

// Server is an interface for HTTP Server
// Use for create testing implementations
// For now, we create a testServer, using memory repositories for testing
// In the future, we can create more type of server, like docker integration server, which use a database inside docker
// Or a SQLite server, which use SQLite for memory database
// Or a transaction server, which run every test in a transaction
type Server interface {
	Init()
	ServeHTTP(w http.ResponseWriter, r *http.Request)

	// Only http handler method will be extract to the interface
	createSignInHandler() gin.HandlerFunc
	createRegisterHandler() gin.HandlerFunc
	createOauth2RegisterHandler() gin.HandlerFunc
	createOauth2SignInHandler() gin.HandlerFunc
	createAuthMiddleware() gin.HandlerFunc
	createPingHandler() gin.HandlerFunc
	createGetProfileHandler() gin.HandlerFunc
}

// routeMap create single source of truth when testing API
// Both real server and test server will create routes using this method
// Note that all route here will be create under the route /api
func routeMap(s Server) map[string]map[string][]gin.HandlerFunc {
	return map[string]map[string][]gin.HandlerFunc{
		"/ping": {
			http.MethodGet: []gin.HandlerFunc{s.createPingHandler()},
		},

		// user auth handler
		"/users/sign-in": {
			http.MethodPost: []gin.HandlerFunc{s.createSignInHandler()},
		},

		"/users/register": {
			http.MethodPost: []gin.HandlerFunc{s.createRegisterHandler()},
		},

		"/oauth2/sign-in": {
			http.MethodPost: []gin.HandlerFunc{s.createOauth2SignInHandler()},
		},

		"/oauth2/register": {
			http.MethodPost: []gin.HandlerFunc{s.createOauth2RegisterHandler()},
		},

		// user handler
		"/users/profile": {
			http.MethodGet: []gin.HandlerFunc{s.createAuthMiddleware(), s.createGetProfileHandler()},
		},
	}
}

type realServer struct {
	router *gin.Engine
	db     *sqlx.DB

	config *ServerConfig
}

func NewServer(config *ServerConfig) *realServer {
	return &realServer{config: config}
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
func (s *realServer) Init() {
	s.connectDB()

	s.router = gin.Default()
	s.initRouter()
}

func (s *realServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *realServer) initRouter() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true
	s.router.Use(cors.New(corsConfig)) // TODO: change this setting later

	// main API routes
	rootAPI := s.router.Group("/api")
	routeMap := routeMap(s)
	for path, methodHandler := range routeMap {
		for method, handlerFunc := range methodHandler {
			rootAPI.Handle(method, path, handlerFunc...)
		}
	}

	// swagger API documentation
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (s *realServer) connectDB() {
	u, err := url.Parse(s.config.SqlConnString)
	if err != nil {
		log.Fatalf("invalid SQL URL %v", err)
	}

	db, err := sqlx.Open(u.Scheme, s.config.SqlConnString)
	if err != nil {
		log.Fatalf("open database failed: %v", err)
	}

	s.db = db
}

type ServerConfig struct {
	// config share across packages are define at root level
	//
	FrontendBaseURL string // the domain where the frontend live
	APIBaseURL      string
	SqlConnString   string

	*AuthConfig
}

// @Summary PING PONG
// @Description For testing
// @Tags ping
// @Produce json
// @Success 200 {object} api.BaseResponse "PING PONG"
// @Router /api/ping [get]
func (s *realServer) createPingHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response(c, 200, "PONG")
	}
}
