package api

import (
	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/auth"
	"github.com/victornm/es-backend/user"
	"net/http"
)

/*
 * Header: Authorization: Basic YWRtaW4xQGVzLmNvbTphZG1pbg==
 * Body: None
 */
func (s *Server) createSignInHandler() func(c *gin.Context) {
	authService := auth.NewService(user.NewMemoryDAO(), s.env.secretKey, s.env.jwtExpiredHour)

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
