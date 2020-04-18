package api

import (
	"github.com/gavv/httpexpect/v2"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type userTester struct {
	suite.Suite
	handler http.Handler
}

func (tester *userTester) TestGetProfile() {
	t := tester.T()
	Convey("Given a sign-in user", t, func() {
		server := httptest.NewServer(tester.handler)
		e := httpexpect.New(t, server.URL)

		token := e.POST("/api/users/sign-in", nil).
			WithBasicAuth("admin@es.com", "admin").
			Expect().
			JSON().Object().
			Value("data").Object().
			Value("token").String().Raw()

		Convey("When get profile with the return token", func() {
			res := e.GET("/api/users/profile").
				WithHeader("Authorization", "Bearer " + token).
				Expect()

			Convey("Then should not error", func() {
				res.Status(http.StatusOK)
				res.JSON().Object().
					ContainsKey("errors").
					ValueEqual("errors", nil)
			})

			Convey("Then should return correct user profile", func() {
				res.JSON().Object().
					ContainsKey("data").
					Value("data").Object().
					ContainsKey("email").
					Value("email").
					String().Equal("admin@es.com")
			})

			Reset(func() {
				server.Close()
			})
		})
	})
}

func TestUser(t *testing.T) {
	suite.Run(t, &userTester{handler: initUnittestServer()})
}
