package main

import (
	"go-metricscol/internal/agent"
	"time"
)

const pollInterval time.Duration = 2 * time.Second
const reportInterval time.Duration = 10 * time.Second

func main() {
	metrics := agent.CreateMetrics()

	pollTimer := time.NewTimer(pollInterval)
	reportTimer := time.NewTimer(reportInterval)

	for {
		select {
		case <-pollTimer.C:
			agent.UpdateMetrics(metrics)
		case <-reportTimer.C:
			metrics.SendToServer(":8080")
		}
	}
}
