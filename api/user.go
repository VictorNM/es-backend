package api

import (
	"github.com/victornm/es-backend/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/store/memory"
)

// @Summary Basic sign in using email, password
// @Description Sign in using email and password
// @Tags user
// @Produce json
// @Success 200 {object} api.BaseResponse{data=authToken} "Sign in successfully"
// @Failure 401 {object} api.BaseResponse{errors=[]api.Error} "Not authenticated"
// @Router /users/sign-in [post]
func (s *Server) createSignInHandler() func(c *gin.Context) {
	authService := s.createBasicSignInService()

	return func(c *gin.Context) {
		email, password, ok := c.Request.BasicAuth()
		if !ok {
			reject(c, http.StatusUnauthorized, user.ErrNotAuthenticated)
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

// @Summary Basic sign in using email, password
// @Description Sign in using email and password
// @Tags user
// @Produce json
// @Param user body user.RegisterMutation true "Register new user"
// @Success 201 {object} api.BaseResponse "Register successfully"
// @Failure 400 {object} api.BaseResponse{errors=[]api.Error} "Bad request"
// @Router /users/register [post]
func (s *Server) createRegisterHandler() func(c *gin.Context) {
	service := s.createRegisterService()

	return func(c *gin.Context) {
		var input *user.RegisterMutation
		if err := c.ShouldBindJSON(&input); err != nil {
			reject(c, http.StatusBadRequest, user.ErrInvalidInput)
			return
		}

		if err := service.Register(input); err != nil {
			reject(c, http.StatusBadRequest, err)
			return
		}

		response(c, http.StatusCreated, nil)
	}
}

// @Summary Get current sign-inned user's profile
// @Description Get profile by user_id in token,
// @Tags user
// @Produce json
// @Success 200 {object} api.BaseResponse{data=user.ProfileDTO} "Get profile successfully"
// @Router /users/profile [get]
func (s *Server) createGetProfileHandler() func(c *gin.Context) {
	userQuery := s.createUserGetProfileQuery()

	return func(c *gin.Context) {
		userAuth := getUser(c)

		u, err := userQuery.GetProfile(userAuth.UserID)
		if err != nil {
			reject(c, http.StatusNotFound, err)
		}

		response(c, http.StatusOK, u)
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

func (s *Server) createAuthTokenParser() user.JWTParserService {
	return user.NewJWTParserService(s.config.JWTSecret)
}

func (s *Server) createBasicSignInService() user.BasicSignInService {
	return user.NewBasicSignInService(memory.UserStore, s.config.JWTSecret, s.config.JWTExpiredHours)
}

func (s *Server) createRegisterService() user.RegisterService {
	return user.NewRegisterService(memory.UserStore, user.NewConsoleSender(memory.UserStore, s.config.FrontendBaseURL))
}

func (s *Server) createUserGetProfileQuery() user.GetProfileQuery {
	return user.NewQueryService(memory.UserStore)
}
