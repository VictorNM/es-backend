package api

import (
	"github.com/stretchr/testify/suite"
	"github.com/victornm/es-backend/api/internal/jsontest"
	"github.com/victornm/es-backend/pkg/store"
	"github.com/victornm/es-backend/pkg/store/memory"
	"github.com/victornm/es-backend/pkg/user"
	"net/http"
	"net/http/httptest"
	"testing"
)

type userTester struct {
	*baseTester

	gateway *mockUserGateway
}

func (tester *userTester) Do(req *http.Request) *httptest.ResponseRecorder {
	origin := createUserFinder
	createUserFinder = func(srv *realServer) user.Finder {
		return tester.gateway
	}
	defer func() {
		createUserFinder = origin
	}()

	return tester.baseTester.Do(req)
}

func (tester *userTester) TestGetProfile() {
	tester.gateway.Seed([]*store.UserRow{
		{
			Email:          "admin@es.com",
			Username:       "admin",
			HashedPassword: mustHashPassword("admin"),
			IsSuperAdmin:   true,
			IsActive:       true,
		},
	})

	tests := map[string]struct {
		isUser   bool
		email    string
		password string

		wantedStatus int
		shouldError  bool
	}{
		"admin should get profile succeed": {
			isUser:   true,
			email:    "admin@es.com",
			password: "admin",

			wantedStatus: http.StatusOK,
			shouldError:  false,
		},

		"not a user should get profile failed": {
			isUser: false,

			wantedStatus: http.StatusUnauthorized,
			shouldError:  true,
		},
	}

	for name, test := range tests {
		tester.T().Run(name, func(t *testing.T) {
			req := jsontest.GET("/api/users/profile")

			if test.isUser {
				token := signIn(tester.srv, test.email, test.password)
				jsontest.SetBearerAuth(req, token)
			}

			w := tester.Do(req)

			tester.Assert().Equal(test.wantedStatus, w.Code)
			tester.Assert().Equal(test.shouldError, getErrors(w) != nil)
		})
	}
}

func TestUser(t *testing.T) {
	suite.Run(t, &userTester{
		baseTester: newBaseTester(initUnittestServer()),
		gateway:    &mockUserGateway{UserGateway: memory.NewUserGateway()},
	})
}

type mockUserGateway struct {
	*memory.UserGateway
}

func (gw *mockUserGateway) Seed(users []*store.UserRow) {
	gw.UserGateway.Seed(users)
}
