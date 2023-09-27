package main

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v9"

	"go-metricscol/internal/server"
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("couldn't parse config with error: %s", err)
	}

	log.Printf("Starting server on %s", cfg.Address)

	s, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("couldn't create server with error: %s", err)
	}

	log.Fatal(s.ListenAndServe())
}

var (
	address       string
	storeInterval time.Duration
	storeFile     string
	restore       bool
	hashKey       string
	databaseDSN   string
)

// Declare variables in which the values of the flags will be written.
func init() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&storeInterval, "i", 300*time.Second, "Interval to store metrics")
	flag.StringVar(&storeFile, "f", "/tmp/devops-metrics-db.json", "File to store metrics")
	flag.BoolVar(&restore, "r", true, "Restore metrics from file")
	flag.StringVar(&hashKey, "k", "", "Key to encrypt metrics")
	flag.StringVar(&databaseDSN, "d", "", "Database DSN")
}

// Parses server.Config from environment variables or flags.
func parseConfig() (*server.Config, error) {
	flag.Parse()
	config := server.NewConfig(address, storeInterval, storeFile, restore, hashKey, databaseDSN)

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
