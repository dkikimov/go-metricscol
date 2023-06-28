package main

import (
	"go-metricscol/internal/agent"
	"log"
	"time"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

func main() {
	metrics := agent.CreateMetrics()

	pollTimer := time.NewTicker(pollInterval)
	reportTimer := time.NewTicker(reportInterval)

	for {
		select {
		case <-pollTimer.C:
			log.Println("Update metrics")
			agent.UpdateMetrics(metrics)
		case <-reportTimer.C:
			log.Println("Send to server")
			metrics.SendToServer("http://127.0.0.1:8080")
		}
	}
}
