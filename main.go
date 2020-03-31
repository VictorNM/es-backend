package main

import (
	"github.com/joho/godotenv"
	"github.com/victornm/es-backend/api"
	"log"
	"net/http"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	s := &api.Server{}
	s.Init()

	log.Fatal(http.ListenAndServe(":8080", s))
}
