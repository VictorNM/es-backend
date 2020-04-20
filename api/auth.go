package api

import (
	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/pkg/auth"
	authMemory "github.com/victornm/es-backend/pkg/auth/repository/memory"
	"github.com/victornm/es-backend/pkg/store/memory"
	"net/http"
	"strings"
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

func (s *realServer) createOauth2RegisterHandler() gin.HandlerFunc {
	service := s.createOAuth2Service()

	return func(c *gin.Context) {
		authURL, err := service.OAuth2Register(c.Param("provider"))
		if err != nil {
			reject(c, http.StatusBadRequest, err)
			return
		}

		c.Redirect(302, authURL)
	}
}

func (s *realServer) createOauth2RegisterCallbackHandler() gin.HandlerFunc {
	service := s.createOAuth2Service()

	return func(c *gin.Context) {
		err := service.OAuth2RegisterCallback(c.Query("state"), c.Query("code"))
		if err != nil {
			reject(c, http.StatusUnauthorized, err)
			return
		}

		response(c, 200, c.Query("code"))
	}
}

func (s *realServer) createAuthMiddleware() gin.HandlerFunc {
	tokenParser := s.createAuthTokenParser()

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if len(authHeader) < 7 || strings.ToUpper(authHeader[0:7]) != "BEARER " {
			abort(c, http.StatusUnauthorized, auth.ErrNotAuthenticated)
			return
		}

		userAuth, err := tokenParser.ParseToken(authHeader[7:])
		if err != nil {
			abort(c, http.StatusUnauthorized, err)
			return
		}

		c.Set("user", userAuth)
	}
}

func (s *realServer) createAuthTokenParser() auth.JWTParserService {
	return auth.NewJWTParserService(s.config.JWTSecret)
}

func (s *realServer) createBasicSignInService() auth.BasicSignInService {
	return auth.NewBasicSignInService(createAuthUserRepository(s), s.config.JWTSecret, s.config.JWTExpiredHours)
}

func (s *realServer) createRegisterService() auth.RegisterService {
	repository := createAuthUserRepository(s)
	return auth.NewRegisterService(repository, auth.NewConsoleSender(repository, s.config.FrontendBaseURL))
}

func (s *realServer) createOAuth2Service() auth.OAuth2RegisterService {
	return auth.NewOAuth2RegisterService(
		createOAuth2StateRepository(s),
		createAuthUserRepository(s),
		s.createOAuth2ClientFactory(),
	)
}

func (s *realServer) createOAuth2ClientFactory() auth.OAuth2ClientFactory {
	return auth.NewOAuth2ClientFactory(auth.WithGoogle(
		s.config.OAuth2GoogleClientID,
		s.config.OAuth2GoogleClientSecret,
		s.config.APIBaseURL+"/api/oauth2/callback"),
	)
}

// TODO: Change to real repository
var createAuthUserRepository = func(srv *realServer) auth.UserRepository {
	return authMemory.NewRepository(memory.GlobalUserStore)
}

// TODO: Change to real repository
var createOAuth2StateRepository = func(srv *realServer) auth.OAuth2StateRepository {
	return authMemory.NewOauth2StateRepository()
}
