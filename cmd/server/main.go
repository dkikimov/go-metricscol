package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v9"

	"go-metricscol/internal/config"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository"
	"go-metricscol/internal/repository/memory"
	"go-metricscol/internal/repository/postgres"
	"go-metricscol/internal/server"
	"go-metricscol/internal/server/backends"
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

	var repo repository.Repository
	if len(cfg.DatabaseDSN) > 0 {
		repo, err = postgres.New(cfg.DatabaseDSN)
		if err != nil {
			log.Fatalf("couldn't create new postgres db: %s", err)
		}
	} else {
		repo = memory.NewMemStorage()
	}

	createdBackend, err := createBackend(backends.GRPCType, repo, cfg)
	if err != nil {
		log.Fatalf("couldn't create backend with error: %s", err)
	}

	s := server.NewServer(cfg, repo, createdBackend)

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
		if err := s.ListenAndServe(serverContext); err != nil {
			log.Fatalf("Server ListenAndServe: %s", err)
		}
		group.Done()
	}()

	<-idleConnsClosed
	group.Wait()
	log.Println("Server Shutdown gracefully")
}

func createBackend(backendType backends.BackendType, repository repository.Repository, cfg *config.ServerConfig) (backends.Backend, error) {
	switch backendType {
	case backends.GRPCType:
		listen, err := net.Listen("tcp", cfg.Address)
		if err != nil {
			return nil, fmt.Errorf("couldn't listen: %s", err)
		}

		return backends.NewGrpc(repository, cfg, listen)
	case backends.HTTPType:
		return backends.NewHTTP(repository, cfg)
	default:
		return nil, fmt.Errorf("unknown backend type id: %d", backendType)
	}
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
	flag.StringVar(&arguments.TrustedSubnet, "t", "", "Trusted subnet")

	arguments.StoreInterval = models.Duration{Duration: 300 * time.Second}
}

// Parses server.ServerConfig from environment variables or flags.
func parseConfig() (*config.ServerConfig, error) {
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

	cfg, err := config.NewServerConfig(
		arguments.Address,
		arguments.StoreInterval,
		arguments.StoreFile,
		arguments.Restore,
		arguments.HashKey,
		arguments.DatabaseDSN,
		arguments.CryptoKeyFilePath,
		arguments.TrustedSubnet,
	)
	if err != nil {
		return nil, fmt.Errorf("couldn't create config: %s", err)
	}

	return cfg, nil
}

func printBuildProperties() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
}
