package api

import (
	"github.com/gavv/httpexpect/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"github.com/victornm/es-backend/pkg/auth"
	"github.com/victornm/es-backend/pkg/store/memory"
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

		Convey("When admin sign-in", func() {
			res := e.POST("/api/users/sign-in", nil).
				WithBasicAuth("admin@es.com", "admin").
				Expect()

			Convey("Then should not error", func() {
				res.Status(http.StatusOK)
				res.JSON().Object().
					ContainsKey("errors").
					ValueEqual("errors", nil)
			})

			Convey("Then should return a token", func() {
				res.JSON().Object().
					ContainsKey("data").
					Value("data").Object().
					ContainsKey("token").
					Value("token").
					String().NotEmpty()
			})

			Reset(func() {
				server.Close()
			})
		})
	})
}

func (tester *authTester) TestRegister() {
	t := tester.T()
	Convey("Given a new server", t, func(c C) {
		server := httptest.NewServer(tester.handler)
		e := httpexpect.New(t, server.URL)
		Convey("When register with valid information", func() {
			res := e.POST("/api/users/register").
				WithJSON(map[string]string{
					"email":                 "victornm@es.com",
					"username":              "victornm",
					"password":              "admin1234",
					"password_confirmation": "admin1234",
					"full_name":             "Nguyen Mau Vinh",
				}).
				Expect()

			Convey("Then should not error", func() {
				res.JSON().Object().
					ContainsKey("errors").
					ValueEqual("errors", nil)
				res.Status(http.StatusCreated)
			})
		})

		Convey("When register with empty body", func() {
			res := e.POST("/api/users/register").
				WithJSON(map[string]string{}).
				Expect()

			Convey("Then should be a bad request", func() {
				res.JSON().Object().
					ContainsKey("errors").
					Value("errors").NotNull()
				res.Status(http.StatusBadRequest)
			})
		})

		Reset(func() {
			server.Close()
		})
	})
}

func TestAuth(t *testing.T) {
	suite.Run(t, &authTester{handler: initUnittestServer()})
}
