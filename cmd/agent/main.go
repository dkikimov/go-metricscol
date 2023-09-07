package main

import (
	"flag"
	"github.com/caarlos0/env/v9"
	"go-metricscol/internal/agent"
	"go-metricscol/internal/repository/memory"
	"golang.org/x/sync/errgroup"
	"log"
	"time"
)

var (
	address        string
	reportInterval time.Duration
	pollInterval   time.Duration
	hashKey        string
	rateLimit      int
)

func main() {
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

func init() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&reportInterval, "r", 10*time.Second, "Interval to report metrics")
	flag.DurationVar(&pollInterval, "p", 2*time.Second, "Interval to poll metrics")
	flag.StringVar(&hashKey, "k", "", "Key to encrypt metrics")
	flag.IntVar(&rateLimit, "l", 1, "Limit the number of requests to the server")
}

func parseConfig() (*agent.Config, error) {
	flag.Parse()
	config := agent.NewConfig(address, reportInterval, pollInterval, hashKey, rateLimit)

	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
