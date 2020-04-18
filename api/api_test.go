package api

import (
	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type unittestServer struct {
	// only re-implement method that access DB of external services
	// re-use all other methods of the realServer
	// so we can create a clear boundary for our application, and test everything inside that boundary
	*realServer
}

// Init for unittestServer ignore any external connection
// only use to init router for testing
func (s *unittestServer) Init() {
	s.router = gin.Default()
	s.initRouter()
}

func initUnittestServer() *unittestServer {
	s := &unittestServer{realServer: NewServer(&ServerConfig{
		FrontendBaseURL: "localhost:3000",
		JWTSecret:       "1232asdasdsd",
		JWTExpiredHours: 12,
	})}

	s.Init()

	return s
}

type pingTester struct {
	suite.Suite
	handler http.Handler
}

// TestPingHandler test the route /api/ping
// This is an example for create tests for api
func (tester *pingTester) TestPing() {
	Convey("Given a new server", tester.T(), func() {
		server := httptest.NewServer(tester.handler)
		e := httpexpect.New(tester.T(), server.URL)

		Convey("When I GET /api/ping", func() {
			req := e.GET("/api/ping")

			Convey("Then response should be 200", func() {
				req.Expect().
					Status(http.StatusOK)
			})

			Convey("And response data should equal PONG", func() {
				req.Expect().JSON().Object().
					ContainsKey("data").
					ValueEqual("data", "PONG")
			})
		})

		Reset(func() {
			server.Close()
		})
	})
}

func TestPing(t *testing.T) {
	tester := &pingTester{handler: initUnittestServer()}

	suite.Run(t, tester)
}
