package api

import (
	"github.com/gin-gonic/gin"
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

func newUnittestServer() *unittestServer {
	return &unittestServer{realServer: NewServer(&ServerConfig{
		FrontendBaseURL: "localhost:3000",
		JWTSecret:       "1232asdasdsd",
		JWTExpiredHours: 12,
	})}
}

func initUnittestServer() *unittestServer {
	s := newUnittestServer()

	s.Init()

	return s
}

type baseTester struct {
	suite.Suite

	srv Server
}

func (tester *baseTester) Do(req *http.Request) *httptest.ResponseRecorder {
	tester.srv.Init()

	w := httptest.NewRecorder()
	tester.srv.ServeHTTP(w, req)
	return w
}

func newBaseTester(srv Server) *baseTester {
	return &baseTester{srv: srv}
}

type pingTester struct {
	*baseTester
}

// TestPingHandler test the route /api/ping
// This is an example for create tests for api
func (tester *pingTester) TestPing() {
	req := newRequest(http.MethodGet, "/api/ping", nil)

	w := tester.Do(req)

	res := getResponse(w)
	tester.Assert().Nil(res.Errors)
	tester.Assert().Equal("PONG", res.Data)
}

func TestPing(t *testing.T) {
	tester := &pingTester{
		baseTester: newBaseTester(initUnittestServer()),
	}

	suite.Run(t, tester)
}
