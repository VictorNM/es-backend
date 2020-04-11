package main

import (
	"flag"
	"fmt"
	"github.com/victornm/es-backend/api"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	var (
		secret          = flag.String("secret", envString("SECRET", "z91NRBxicpx2qjvO"), "secret key for JWT")
		expiredHours    = flag.Int("token-expired-hours", envInt("TOKEN_EXPIRED_HOURS", 24), "expired duration in hour for JWT token")
		frontEndBaseURL = flag.String("frontend-base-url", envString("FRONT_END_BASE_URL", "localhost:3000"), "")
		httpPort        = flag.Int("http-port", envInt("HTTP_PORT", 8080), "port listening")
	)

	flag.Parse()

	config := &api.ServerConfig{
		FrontendBaseURL: *frontEndBaseURL,

		JWTSecret:       *secret,
		JWTExpiredHours: *expiredHours,
	}

	s := api.NewServer(config)
	s.Init()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *httpPort), s))
}

func envString(key string, value string) string {
	env := os.Getenv(key)
	if len(env) == 0 {
		return value
	}

	return env
}

func envInt(key string, value int) int {
	env, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return value
	}

	return env
}