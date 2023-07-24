package main

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"go-metricscol/internal/agent"
	"go-metricscol/internal/models"
	"log"
	"time"
)

var (
	address        string
	reportInterval time.Duration
	pollInterval   time.Duration
	hashKey        string
)

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("couldn't parse config with error: %s", err)
	}

	metrics := models.NewMetrics()

	pollTimer := time.NewTicker(cfg.PollInterval)
	reportTimer := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-pollTimer.C:
			log.Println("Update metrics")
			agent.UpdateMetrics(&metrics)
		case <-reportTimer.C:
			log.Printf("Send metrics to %s\n", cfg.Address)
			if err := agent.SendMetricsToServer(cfg.Address, &metrics, cfg.HashKey); err != nil {
				log.Printf("Error while sending metrics to server: %s", err)
			}
		}
	}
}

func init() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&reportInterval, "r", 10*time.Second, "Interval to report metrics")
	flag.DurationVar(&pollInterval, "p", 2*time.Second, "Interval to poll metrics")
	flag.StringVar(&hashKey, "k", "", "Key to encrypt metrics")
}

func parseConfig() (*agent.Config, error) {
	flag.Parse()
	config := agent.NewConfig(address, reportInterval, pollInterval, hashKey)

	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
