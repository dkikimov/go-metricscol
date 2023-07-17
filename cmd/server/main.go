package main

import (
	"github.com/caarlos0/env/v9"
	"go-metricscol/internal/server"
	"log"
)

func main() {
	cfg := server.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Couldn't parse config with error: %s", err)
	}

	log.Printf("Starting server on %s", cfg.Address)
	srv := server.New(cfg.Address)
	log.Fatal(srv.ListenAndServe())
}
