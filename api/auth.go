package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/pkg/auth"
	"github.com/victornm/es-backend/pkg/auth/mock"
	"github.com/victornm/es-backend/pkg/store/memory"
)

type AuthConfig struct {
	JWTSecret       string
	JWTExpiredHours int

	OAuth2GoogleClientID     string
	OAuth2GoogleClientSecret string
}

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

// @Summary Register using oauth2
// @Description Register using oauth2
// @Tags auth
// @Produce json
// @Param user body auth.OAuth2Input true "Register new user using oauth2"
// @Success 201 {object} api.BaseResponse "Register successfully"
// @Failure 400 {object} api.BaseResponse{errors=[]api.Error} "Bad request"
// @Router /oauth2/register [post]
func (s *realServer) createOauth2RegisterHandler() gin.HandlerFunc {
	service := s.createOAuth2RegisterService()

	return func(c *gin.Context) {
		var input auth.OAuth2Input
		err := c.ShouldBindJSON(&input)
		if err != nil {
			reject(c, http.StatusBadRequest, auth.ErrInvalidInput)
			return
		}

		err = service.OAuth2Register(input)
		if err != nil {
			reject(c, http.StatusUnauthorized, err)
			return
		}

		response(c, http.StatusCreated, nil)
	}
}

// @Summary Sign in using oauth2
// @Description Sign in using oauth2
// @Tags auth
// @Produce json
// @Param user body auth.OAuth2Input true "Sign in using oauth2"
// @Success 200 {object} api.BaseResponse{data=authToken} "Sign in successfully"
// @Failure 401 {object} api.BaseResponse{errors=[]api.Error} "Not authenticated"
// @Router /oauth2/sign-in [post]
func (s *realServer) createOauth2SignInHandler() gin.HandlerFunc {
	service := s.createOAuth2SignInService()

	return func(c *gin.Context) {
		var input auth.OAuth2Input
		err := c.ShouldBindJSON(&input)
		if err != nil {
			reject(c, http.StatusBadRequest, auth.ErrInvalidInput)
			return
		}

		tokenString, err := service.OAuth2SignIn(input)
		if err != nil {
			reject(c, http.StatusUnauthorized, err)
			return
		}

		response(c, http.StatusCreated, authToken{Token: tokenString})
	}
}

func (s *realServer) createAuthMiddleware() gin.HandlerFunc {
	tokenParser := s.createJWTService()

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

func (s *realServer) createJWTService() auth.JWTService {
	return auth.NewJWTService(s.config.JWTSecret, s.config.JWTExpiredHours)
}

func (s *realServer) createBasicSignInService() auth.BasicSignInService {
	return auth.NewBasicSignInService(createAuthUserRepository(s), s.createJWTService())
}

func (s *realServer) createRegisterService() auth.RegisterService {
	repository := createAuthUserRepository(s)
	return auth.NewRegisterService(repository, auth.NewActivationEmailSender(repository, s.config.FrontendBaseURL+"/activate"))
}

func (s *realServer) createOAuth2RegisterService() auth.OAuth2RegisterService {
	return auth.NewOAuth2RegisterService(createAuthUserRepository(s), createOAuth2ClientFactory(s))
}

func (s *realServer) createOAuth2SignInService() auth.OAuth2SignInService {
	return auth.NewOAuth2SignInService(createAuthUserRepository(s), createOAuth2ClientFactory(s), s.createJWTService())
}

// TODO: Change to real repository
var createAuthUserRepository = func(s *realServer) auth.UserRepository {
	return mock.NewRepository(memory.GlobalUserStore)
}

var createOAuth2ClientFactory = func(s *realServer) auth.OAuth2ProviderFactory {
	return auth.NewOAuth2ClientFactory(auth.NewGoogleProvider(
		s.config.OAuth2GoogleClientID,
		s.config.OAuth2GoogleClientSecret,
	))
}
