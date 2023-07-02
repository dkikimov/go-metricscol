package main

import (
	"go-metricscol/internal/agent"
	"go-metricscol/internal/models"
	"log"
	"time"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

func main() {
	metrics := models.Metrics{}

	pollTimer := time.NewTicker(pollInterval)
	reportTimer := time.NewTicker(reportInterval)

	for {
		select {
		case <-pollTimer.C:
			log.Println("Update metrics")
			agent.UpdateMetrics(metrics)
		case <-reportTimer.C:
			log.Println("Send to server")
			if err := agent.SendMetricsToServer("http://127.0.0.1:8080", metrics); err != nil {
				log.Fatalf(err.Error())
			}
		}
	}
}
