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
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v9"

	"go-metricscol/internal/models"
	"go-metricscol/internal/server"
)

// go run -ldflags "-X main.buildVersion=v1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d')' -X 'main.buildCommit=$(git rev-parse --short HEAD)'" main.go
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildProperties()

	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("couldn't parse config with error: %s", err)
	}

	log.Printf("Starting server on %s", cfg.Address)

	s, err := server.NewServer(cfg)
	serverContext, serverContextCancel := context.WithCancel(context.Background())
	if err != nil {
		log.Fatalf("couldn't create server with error: %s", err)
	}

	idleConnsClosed := make(chan struct{})
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		<-sigChan

		serverContextCancel()
		close(idleConnsClosed)
	}()

	group := sync.WaitGroup{}
	group.Add(1)
	go func() {
		if err := s.ListenAndServe(serverContext); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
		group.Done()
	}()

	<-idleConnsClosed
	group.Wait()
	log.Println("Server Shutdown gracefully")
}

var jsonParsedArguments commandLineArguments
var arguments commandLineArguments

// Declare variables in which the values of the flags will be written.
func init() {
	flag.StringVar(&arguments.Address, "a", "127.0.0.1:8080", "Address to listen")
	flag.Var(&arguments.StoreInterval, "i", "Interval to store metrics")
	flag.StringVar(&arguments.StoreFile, "f", "/tmp/devops-metrics-db.json", "File to store metrics")
	flag.BoolVar(&arguments.Restore, "r", true, "Restore metrics from file")
	flag.StringVar(&arguments.HashKey, "k", "", "Key to encrypt metrics")
	flag.StringVar(&arguments.DatabaseDSN, "d", "", "Database DSN")
	flag.StringVar(&arguments.CryptoKeyFilePath, "crypto-key", "", "Private crypto key for asymmetric encryption")
	flag.StringVar(&arguments.JSONConfigPath, "c", "", "Path to json config")

	arguments.StoreInterval = models.Duration{Duration: 300 * time.Second}
}

// Parses server.Config from environment variables or flags.
func parseConfig() (*server.Config, error) {
	flag.Parse()

	// Parse from JSON configuration file.
	if len(jsonParsedArguments.JSONConfigPath) != 0 {
		jsonConfig, err := os.ReadFile(jsonParsedArguments.JSONConfigPath)
		if err != nil {
			return nil, fmt.Errorf("couldn't read config file")
		}

		if err := json.Unmarshal(jsonConfig, &jsonParsedArguments); err != nil {
			return nil, fmt.Errorf("couldn't unmarshal json config")
		}
	}

	// Parse from flags
	arguments.Merge(jsonParsedArguments)

	// Parse from environment variables.
	opts := env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(arguments.StoreInterval): models.ParseDurationFromEnv,
		},
	}
	if err := env.ParseWithOptions(&arguments, opts); err != nil {
		return nil, fmt.Errorf("couldn't parse config from env: %s", err)
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
