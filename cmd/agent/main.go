package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/caarlos0/env/v9"
	"golang.org/x/sync/errgroup"

	"go-metricscol/internal/agent"
	"go-metricscol/internal/models"
	"go-metricscol/internal/repository/memory"
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

	metrics := memory.NewMetrics()
	agentClient, err := agent.NewAgent(cfg, agent.HTTP)
	if err != nil {
		log.Fatalf("couldn't create agent with error: %s", err)
	}

	pollTimer := time.NewTicker(cfg.PollInterval)
	reportTimer := time.NewTicker(cfg.ReportInterval)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	for {
		select {
		case <-pollTimer.C:
			log.Println("Update metrics")

			pollTimer.Stop()
			g := errgroup.Group{}
			g.Go(func() error {
				return agent.UpdateMetrics(&metrics)
			})
			g.Go(func() error {
				return agent.CollectAdditionalMetrics(&metrics)
			})

			if err := g.Wait(); err != nil {
				log.Printf("Couldn't collect metrics: %s", err)
			}
			pollTimer.Reset(cfg.PollInterval)
		case <-reportTimer.C:
			log.Printf("Send metrics to %s\n", cfg.Address)

			reportTimer.Stop()
			go func() {
				if err := agentClient.SendMetricsToServer(&metrics); err != nil {
					log.Printf("Error while sending metrics to server: %s", err)
				}
			}()
			reportTimer.Reset(cfg.ReportInterval)
		case <-sigChan:
			log.Printf("Send metrics to %s before graceful shutdown \n", cfg.Address)

			pollTimer.Stop()
			reportTimer.Stop()

			if err := agentClient.SendMetricsToServer(&metrics); err != nil {
				log.Printf("Error while sending metrics to server: %s", err)
			}

			log.Print("Agent graceful shutdown \n")
			os.Exit(0)
		}
	}

}

var jsonParsedArguments commandLineArguments
var arguments commandLineArguments

// Declare variables in which the values of the flags will be written.
func init() {
	flag.StringVar(&arguments.Address, "a", agent.DefaultAddress, "Address to listen")
	flag.Var(&arguments.ReportInterval, "r", "Interval to report metrics")
	flag.Var(&arguments.PollInterval, "p", "Interval to poll metrics")
	flag.StringVar(&arguments.HashKey, "k", "", "Key to encrypt metrics")
	flag.IntVar(&arguments.RateLimit, "l", 1, "Limit the number of requests to the server")
	flag.StringVar(&arguments.CryptoKeyFilePath, "crypto-key", "", "Private crypto key for asymmetric encryption")
	flag.StringVar(&arguments.JSONConfigPath, "c", "", "Path to json config")

	arguments.ReportInterval = models.Duration{Duration: 10 * time.Second}
	arguments.PollInterval = models.Duration{Duration: 2 * time.Second}
}

// Parses agent.Config from environment variables or flags.
func parseConfig() (*agent.Config, error) {
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
			reflect.TypeOf(arguments.ReportInterval): models.ParseDurationFromEnv,
		},
	}
	if err := env.ParseWithOptions(&arguments, opts); err != nil {
		return nil, fmt.Errorf("couldn't parse config from env: %s", err)
	}

	config, err := agent.NewConfig(
		arguments.Address,
		arguments.ReportInterval,
		arguments.PollInterval,
		arguments.HashKey,
		arguments.RateLimit,
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
