package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v9"

	"go-metricscol/internal/server"
)

// go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')' -X 'main.buildCommit=$(git rev-parse --short HEAD)'" main.go
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	printBuildProperties()

	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("couldn't parse config with error: %s", err)
	}

	log.Printf("Starting server on %s", cfg.Address)

	s, err := server.NewServer(cfg)
	httpServer := s.GetHttpServer()

	idleConnsClosed := make(chan struct{})
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		<-sigint

		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("couldn't shutdown HTTP server: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err != nil {
		log.Fatalf("couldn't create server with error: %s", err)
	}

	if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
	log.Println("Server Shutdown gracefully")
}

type commandLineArguments struct {
	Address           string        `json:"address,omitempty" env:"ADDRESS"`
	StoreInterval     time.Duration `json:"store_interval,omitempty" env:"STORE_INTERVAL"`
	StoreFile         string        `json:"store_file,omitempty" env:"STORE_FILE"`
	Restore           bool          `json:"restore,omitempty" env:"RESTORE"`
	HashKey           string        `json:"hash_key,omitempty" env:"KEY"`
	DatabaseDSN       string        `json:"database_dsn,omitempty" env:"DATABASE_DSN"`
	CryptoKeyFilePath string        `json:"crypto_key_file_path,omitempty" env:"CRYPTO_KEY"`
	JsonConfigPath    string        `env:"CONFIG"`
}

var arguments commandLineArguments

// Declare variables in which the values of the flags will be written.
func init() {
	flag.StringVar(&arguments.Address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&arguments.StoreInterval, "i", 300*time.Second, "Interval to store metrics")
	flag.StringVar(&arguments.StoreFile, "f", "/tmp/devops-metrics-db.json", "File to store metrics")
	flag.BoolVar(&arguments.Restore, "r", true, "Restore metrics from file")
	flag.StringVar(&arguments.HashKey, "k", "", "Key to encrypt metrics")
	flag.StringVar(&arguments.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&arguments.CryptoKeyFilePath, "crypto-key", "", "Private crypto key for asymmetric encryption")
	flag.StringVar(&arguments.JsonConfigPath, "c", "", "Path to json config")
}

// Parses server.Config from environment variables or flags.
func parseConfig() (*server.Config, error) {
	flag.Parse()

	if err := env.Parse(&arguments); err != nil {
		return nil, fmt.Errorf("couldn't parse config from env: %s", err)
	}

	if len(arguments.JsonConfigPath) != 0 {
		jsonConfig, err := os.ReadFile(arguments.JsonConfigPath)
		if err != nil {
			return nil, fmt.Errorf("couldn't read config file")
		}

		if err := json.Unmarshal(jsonConfig, &arguments); err != nil {
			return nil, fmt.Errorf("couldn't unmarshal json config")
		}
	}

	config, err := server.NewConfig(
		arguments.Address,
		arguments.StoreInterval,
		arguments.StoreFile,
		arguments.Restore,
		arguments.HashKey,
		arguments.DatabaseDSN,
		arguments.CryptoKeyFilePath,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't create config: %s", err)
	}

	return config, nil
}

func printBuildProperties() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
}
