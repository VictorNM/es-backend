package api

import (
	"github.com/victornm/es-backend/pkg/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/pkg/store/memory"
)

// @Summary Basic sign in using email, password
// @Description Sign in using email and password
// @Tags auth
// @Produce json
// @Success 200 {object} api.BaseResponse{data=authToken} "Sign in successfully"
// @Failure 401 {object} api.BaseResponse{errors=[]api.Error} "Not authenticated"
// @Router /users/sign-in [post]
func (s *realServer) createSignInHandler() gin.HandlerFunc {
	authService := s.createBasicSignInService()

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

		response(c, http.StatusOK, authToken{Token: tokenString})
	}
}

type authToken struct {
	Token string `json:"token"`
}

// @Summary Register using email and password
// @Description Register using email and password
// @Tags auth
// @Produce json
// @Param user body auth.RegisterInput true "Register new user"
// @Success 201 {object} api.BaseResponse "Register successfully"
// @Failure 400 {object} api.BaseResponse{errors=[]api.Error} "Bad request"
// @Router /users/register [post]
func (s *realServer) createRegisterHandler() gin.HandlerFunc {
	service := s.createRegisterService()

	return func(c *gin.Context) {
		var input *auth.RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			reject(c, http.StatusBadRequest, auth.ErrInvalidInput)
			return
		}

		if err := service.Register(input); err != nil {
			reject(c, http.StatusBadRequest, err)
			return
		}

		response(c, http.StatusCreated, nil)
	}
}

func (s *realServer) createAuthMiddleware() gin.HandlerFunc {
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

func (s *realServer) createAuthTokenParser() auth.JWTParserService {
	return auth.NewJWTParserService(s.config.JWTSecret)
}

func (s *realServer) createBasicSignInService() auth.BasicSignInService {
	return auth.NewBasicSignInService(s.createAuthUserRepository(), s.config.JWTSecret, s.config.JWTExpiredHours)
}

func (s *realServer) createRegisterService() auth.RegisterService {
	return auth.NewRegisterService(s.createAuthUserRepository(), auth.NewConsoleSender(s.createAuthUserRepository(), s.config.FrontendBaseURL))
}

func (s *realServer) createAuthUserRepository() auth.UserRepository {
	return auth.NewRepository(memory.UserStore)
}