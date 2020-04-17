package api

import (
	"github.com/gavv/httpexpect/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"github.com/victornm/es-backend/auth"
	"github.com/victornm/es-backend/store/memory"
	"net/http"
	"net/http/httptest"
	"testing"
)

type authTester struct {
	suite.Suite
	handler http.Handler
}

// mock database
func (s *unittestServer) createAuthUserRepository() auth.UserRepository {
	return auth.NewRepository(memory.NewUserStore())
}

func (tester *authTester) TestSignIn() {
	t := tester.T()
	Convey("Given a new server", t, func() {
		server := httptest.NewServer(tester.handler)
		e := httpexpect.New(t, server.URL)

		Convey("When user sign-in", func() {
			res := e.POST("/api/users/sign-in", nil).
				WithBasicAuth("admin@es.com", "admin").
				Expect()

			Convey("Then should not error", func() {
				res.JSON().Object().
					ContainsKey("errors").
					ValueEqual("errors", nil)
			})

			Convey("Then should return a token", func() {
				res.Status(http.StatusOK)
				res.JSON().Object().
					ContainsKey("data").
					Value("data").Object().
					ContainsKey("token").
					Value("token").
					String().NotEmpty()
			})
		})
	})
}

func TestAuth(t *testing.T) {
	suite.Run(t, &authTester{handler:initUnittestServer()})
}
