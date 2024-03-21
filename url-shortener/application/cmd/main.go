package main

import (
	"fmt"
	"log"
	"net/http"
	"url-shortener/cmd/handlers"

	"github.com/kelseyhightower/envconfig"
)

func main() {

	var s handlers.Config
	err := envconfig.Process("app", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("main : app started")

	api := http.Server{
		Addr:    fmt.Sprintf("localhost:%s", s.Port),
		Handler: handlers.NewServer(s).Router,
	}

	if err := api.ListenAndServe(); err != nil {
		log.Fatal("error: %s", err)
	}
}
