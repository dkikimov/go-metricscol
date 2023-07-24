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
	databaseDSN   string
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("couldn't parse config with error: %s", err)
	}

	log.Printf("Starting server on %s", cfg.Address)

	s, err := server.NewServer(cfg, memory.NewMemStorage(cfg.HashKey))
	if err != nil {
		log.Fatalf("couldn't create server with error: %s", err)
	}

	log.Fatal(s.ListenAndServe())
}

func init() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&storeInterval, "i", 300*time.Second, "Interval to store metrics")
	flag.StringVar(&storeFile, "f", "/tmp/devops-metrics-db.json", "File to store metrics")
	flag.BoolVar(&restore, "r", true, "Restore metrics from file")
	flag.StringVar(&hashKey, "k", "", "Key to encrypt metrics")
	flag.StringVar(&databaseDSN, "d", "", "Database DSN")
}

func parseConfig() (*server.Config, error) {
	flag.Parse()
	config := server.NewConfig(address, storeInterval, storeFile, restore, hashKey, databaseDSN)

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
