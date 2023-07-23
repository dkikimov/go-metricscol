package main

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/server"
	"log"
	"time"
)

var (
	address       string
	storeInterval time.Duration
	storeFile     string
	restore       bool
	hashKey       string
)

func main() {
	cfg := parseConfig()

	log.Printf("Starting server on %s", cfg.Address)

	s := server.NewServer(cfg, memory.NewMemStorage(cfg.HashKey))
	log.Fatal(s.ListenAndServe())
}

func init() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&storeInterval, "i", 300*time.Second, "Interval to store metrics")
	flag.StringVar(&storeFile, "f", "/tmp/devops-metrics-db.json", "File to store metrics")
	flag.BoolVar(&restore, "r", true, "Restore metrics from file")
	flag.StringVar(&hashKey, "k", "", "Key to encrypt metrics")
}

func parseConfig() *server.Config {
	flag.Parse()
	config := server.NewConfig(address, storeInterval, storeFile, restore, hashKey)

	if err := env.Parse(config); err != nil {
		log.Fatalf("Couldn't parse config with error: %s", err)
	}
	return config
}
