package main

import (
	"github.com/caarlos0/env/v9"
	"go-metricscol/internal/agent"
	"go-metricscol/internal/models"
	"log"
	"time"
)

func main() {
	cfg := agent.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Couldn't parse config with error: %s", err)
	}

	metrics := models.MetricsMap{}

	pollTimer := time.NewTicker(cfg.PollInterval)
	reportTimer := time.NewTicker(cfg.ReportInterval)

	for {
		select {
		case <-pollTimer.C:
			log.Println("Update metrics")
			agent.UpdateMetrics(metrics)
		case <-reportTimer.C:
			log.Printf("Send metrics to %s\n", cfg.Address)
			if err := agent.SendMetricsToServer(cfg.Address, metrics); err != nil {
				log.Printf("Error while sending metrics to server: %s", err)
			}
		}
	}
}
