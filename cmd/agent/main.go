package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v9"
	"golang.org/x/sync/errgroup"

	"go-metricscol/internal/agent"
	"go-metricscol/internal/repository/memory"
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

	metrics := memory.NewMetrics()

	pollTimer := time.NewTicker(cfg.PollInterval)
	reportTimer := time.NewTicker(cfg.ReportInterval)

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
				if err := agent.SendMetricsToServer(cfg, &metrics); err != nil {
					log.Printf("Error while sending metrics to server: %s", err)
				}
			}()
			reportTimer.Reset(cfg.ReportInterval)
		}
	}

}

type commandLineArguments struct {
	Address           string        `json:"address,omitempty" env:"ADDRESS"`
	ReportInterval    time.Duration `json:"report_interval,omitempty" env:"REPORT_INTERVAL"`
	PollInterval      time.Duration `json:"poll_interval,omitempty" env:"POLL_INTERVAL"`
	HashKey           string        `json:"hash_key,omitempty" env:"KEY"`
	RateLimit         int           `json:"rate_limit,omitempty" env:"RATE_LIMIT"`
	CryptoKeyFilePath string        `json:"crypto_key_file_path,omitempty" env:"CRYPTO_KEY"`
	JsonConfigPath    string        `env:"CONFIG"`
}

var arguments commandLineArguments

// Declare variables in which the values of the flags will be written.
func init() {
	flag.StringVar(&arguments.Address, "a", agent.DefaultAddress, "Address to listen")
	flag.DurationVar(&arguments.ReportInterval, "r", 10*time.Second, "Interval to report metrics")
	flag.DurationVar(&arguments.PollInterval, "p", 2*time.Second, "Interval to poll metrics")
	flag.StringVar(&arguments.HashKey, "k", "", "Key to encrypt metrics")
	flag.IntVar(&arguments.RateLimit, "l", 1, "Limit the number of requests to the server")
	flag.StringVar(&arguments.CryptoKeyFilePath, "crypto-key", "", "Private crypto key for asymmetric encryption")
	flag.StringVar(&arguments.JsonConfigPath, "c", "", "Path to json config")
}

// Parses agent.Config from environment variables or flags.
func parseConfig() (*agent.Config, error) {
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
