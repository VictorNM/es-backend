package api

import (
	"github.com/stretchr/testify/suite"
	"github.com/victornm/es-backend/api/internal/jsontest"
	"github.com/victornm/es-backend/pkg/auth"
	"github.com/victornm/es-backend/pkg/auth/mock"
	gatewayMemory "github.com/victornm/es-backend/pkg/store/memory"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authTester struct {
	*baseTester

	repository TestAuthUserRepository
	provider   TestOAuth2Provider
}

func mustHashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(hashed)
}

func (tester *authTester) Do(req *http.Request) *httptest.ResponseRecorder {
	originRepository := createAuthUserRepository
	createAuthUserRepository = func(srv *realServer) auth.UserRepository {
		return tester.repository
	}

	originFactory := createOAuth2ClientFactory
	createOAuth2ClientFactory = func(s *realServer) auth.OAuth2ProviderFactory {
		return auth.NewOAuth2ClientFactory(tester.provider)
	}

	defer func() {
		createAuthUserRepository = originRepository
		createOAuth2ClientFactory = originFactory
	}()

	return tester.baseTester.Do(req)
}

func (tester *authTester) TestSignIn() {
	tester.repository.Seed([]*auth.User{
		{
			Email:          "admin@es.com",
			HashedPassword: mustHashPassword("admin"),
			IsActive:       true,
		},
	})

	tests := map[string]struct {
		email    string
		password string

		wantedStatus int
		shouldError  bool
	}{
		"admin sign-in should succeed": {
			email:    "admin@es.com",
			password: "admin",

			wantedStatus: http.StatusOK,
			shouldError:  false,
		},

		"email not exist should failed": {
			email:    "victornm@es.com",
			password: "victornm",

			wantedStatus: http.StatusUnauthorized,
			shouldError:  true,
		},
	}

	for name, test := range tests {
		tester.T().Run(name, func(t *testing.T) {
			req := jsontest.WrapPOST("/api/users/sign-in", nil).
				SetBasicAuth(test.email, test.password).
				Unwrap()

			w := tester.Do(req)

			tester.Assert().Equal(test.wantedStatus, w.Code)
			if test.shouldError {
				tester.Assert().NotNil(getErrors(w))
			} else {
				tester.Assert().Nil(getErrors(w))
				tester.Assert().NotEmpty(getDataAsMap(w)["token"])
			}
		})
	}

	tester.repository.Clear()
}

func (tester *authTester) TestRegister() {
	tests := map[string]struct {
		body map[string]string

		wantedStatus int
		shouldError  bool
	}{
		"register with valid information should succeed": {
			body: map[string]string{
				"email":                 "victornm@es.com",
				"username":              "victornm",
				"password":              "admin1234",
				"password_confirmation": "admin1234",
				"full_name":             "Nguyen Mau Vinh",
			},

			wantedStatus: http.StatusCreated,
			shouldError:  false,
		},

		"register with empty body should be a bad request": {
			body: map[string]string{},

			wantedStatus: http.StatusBadRequest,
			shouldError:  true,
		},
	}

	for name, test := range tests {
		tester.T().Run(name, func(t *testing.T) {
			req := jsontest.POST("/api/users/register", test.body)

			w := tester.Do(req)

			tester.Assert().Equal(test.wantedStatus, w.Code)
			tester.Assert().Equal(test.shouldError, getErrors(w) != nil)
		})
	}

	tester.repository.Clear()
}

func (tester *authTester) TestOAuth2Register() {
	providerName := mock.ProviderName
	tester.repository.Seed([]*auth.User{
		{
			Email:    "foo@bar.com",
			Provider: providerName,
			IsActive: true,
		},
	})
	tester.provider.Seed(map[string]*auth.User{
		"code_1": {Email: "victornm@es.com", Provider: providerName},
	})

	tests := map[string]struct {
		body map[string]string

		wantedStatus int
		shouldError  bool
	}{
		"": {
			body: map[string]string{
				"provider": mock.ProviderName,
				"code":     "code_1",
			},

			wantedStatus: http.StatusCreated,
			shouldError:  false,
		},
	}

	for name, test := range tests {
		tester.T().Run(name, func(t *testing.T) {
			req := jsontest.POST("/api/oauth2/register", test.body)

			w := tester.Do(req)

			tester.Assert().Equal(test.wantedStatus, w.Code)
			tester.Assert().Equal(test.shouldError, getErrors(w) != nil)
		})
	}

	tester.repository.Clear()
}

func TestAuth(t *testing.T) {
	tester := &authTester{
		baseTester: newBaseTester(newUnittestServer()),
		repository: mock.NewRepository(gatewayMemory.NewUserGateway()),
		provider:   mock.NewOAuth2Provider(),
	}
	suite.Run(t, tester)
}

type TestAuthUserRepository interface {
	auth.UserRepository

	// Seed prepare database for testing
	Seed(users []*auth.User)

	// Clear reset to an empty database
	Clear()
}

type TestOAuth2Provider interface {
	auth.OAuth2Provider

	Seed(users map[string]*auth.User)

	Clear()
}
