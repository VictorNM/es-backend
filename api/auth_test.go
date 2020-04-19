package api

import (
	"github.com/stretchr/testify/suite"
	"github.com/victornm/es-backend/pkg/auth"
	"github.com/victornm/es-backend/pkg/auth/repository/memory"
	"github.com/victornm/es-backend/pkg/store"
	gatewayMemory "github.com/victornm/es-backend/pkg/store/memory"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authTester struct {
	*baseTester

	userGateway *gatewayMemory.UserGateway
}

func genPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	return string(hashed)
}

func (tester *authTester) Do(req *http.Request) *httptest.ResponseRecorder {
	origin := createAuthUserRepository
	createAuthUserRepository = func(srv *realServer) auth.UserRepository {
		return memory.NewRepository(tester.userGateway)
	}

	defer func() {
		createAuthUserRepository = origin
	}()

	return tester.baseTester.Do(req)
}

func (tester *authTester) TestSignIn() {
	method := http.MethodPost
	target := "/api/users/sign-in"

	tester.userGateway.Seed([]*store.UserRow{
		{
			Email:          "admin@es.com",
			HashedPassword: genPassword("admin"),
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
			req := newRequest(method, target, nil)
			req.SetBasicAuth(test.email, test.password)

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
}

func (tester *authTester) TestRegister() {
	method := http.MethodPost
	target := "/api/users/register"

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
			req := newRequest(method, target, test.body)

			w := tester.Do(req)

			tester.Assert().Equal(test.wantedStatus, w.Code)
			tester.Assert().Equal(test.shouldError, getErrors(w) != nil)
		})
	}
}

func TestAuth(t *testing.T) {
	tester := &authTester{
		baseTester:  newBaseTester(newUnittestServer()),
		userGateway: gatewayMemory.NewUserGateway(),
	}
	suite.Run(t, tester)
}