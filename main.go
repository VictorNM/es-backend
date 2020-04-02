package main

import (
	"fmt"
	"github.com/victornm/es-backend/api"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	config := &api.ServerConfig{
		JWTSecret:       os.Getenv("SECRET"),
		JWTExpiredHours: envAsInt("TOKEN_EXPIRED_HOURS"),
	}

	s := api.NewServer(config)
	s.Init()

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", envAsInt("HTTP_PORT")), s))
}

func envAsInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		log.Fatal(err)
	}

	return value
}
