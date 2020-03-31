package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Server struct {
	router *gin.Engine

	env env
}

func (s *Server) Init() {
	s.router = gin.Default()
	rootAPI := s.router.Group("/api")

	// auth
	authGroup := rootAPI.Group("/auth")
	authGroup.POST("/sign-in", s.createSignInHandler())

	jwtExpiredHour, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRED_HOURS"))
	if err != nil {
		log.Fatal(errors.New(`env invalid: "TOKEN_EXPIRED_HOURS"`))
	}

	s.env = env{
		secretKey:      os.Getenv("SECRET"),
		jwtExpiredHour: jwtExpiredHour,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

type env struct {
	secretKey      string
	jwtExpiredHour int
}
