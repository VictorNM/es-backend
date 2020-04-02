package api

import (
	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/auth"
	"github.com/victornm/es-backend/store/memory"
	"net/http"
	"strings"
)

/*
 * Header: Authorization: Basic YWRtaW4xQGVzLmNvbTphZG1pbg==
 * Body: None
 */
func (s *Server) createSignInHandler() func(c *gin.Context) {
	authService := s.createAuthBasicSignIner()

	return func(c *gin.Context) {
		email, password, ok := c.Request.BasicAuth()
		if !ok {
			reject(c, http.StatusUnauthorized, auth.ErrNotAuthenticated)
			return
		}

		tokenString, err := authService.BasicSignIn(email, password)
		if err != nil {
			reject(c, http.StatusUnauthorized, err)
			return
		}

		response(c, http.StatusOK, map[string]string{"token": tokenString})
	}
}

func (s *Server) createAuthMiddleware() gin.HandlerFunc {
	tokenParser := s.createAuthTokenParser()

	return func(c *gin.Context) {
		// Look for an Authorization header
		if authHeader := c.GetHeader("Authorization"); authHeader != "" {
			// Should be a bearer token
			if len(authHeader) > 6 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
				userAuth, err := tokenParser.ParseToken(authHeader[7:])
				if err != nil {
					reject(c, http.StatusUnauthorized, err)
					return
				}

				c.Set("user", userAuth)
			}
		}
	}
}

func (s *Server) createAuthTokenParser() auth.TokenParser {
	return auth.NewService(memory.NewUserStore(), s.config.JWTSecret, s.config.JWTExpiredHours)
}

func (s *Server) createAuthBasicSignIner() auth.BasicSignIner {
	return auth.NewService(memory.NewUserStore(), s.config.JWTSecret, s.config.JWTExpiredHours)
}
