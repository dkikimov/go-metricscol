package main

import (
	"crypto/rsa"
	"flag"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/caarlos0/env/v9"
	"golang.org/x/crypto/ssh"
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

var (
	address        string
	reportInterval time.Duration
	pollInterval   time.Duration
	hashKey        string
	rateLimit      int
	cryptoKey      *rsa.PublicKey
)

func rsaPublicKeyParser(input string) (interface{}, error) {
	var result *rsa.PublicKey
	if len(input) != 0 {
		cryptoKeyBytes, err := os.ReadFile(input)
		if err != nil {
			return nil, err
		}

		parsed, _, _, _, err := ssh.ParseAuthorizedKey(cryptoKeyBytes)
		if err != nil {
			return nil, err
		}

		parsedCryptoKey := parsed.(ssh.CryptoPublicKey)
		pubCrypto := parsedCryptoKey.CryptoPublicKey()
		result = pubCrypto.(*rsa.PublicKey)
	}

	return *result, nil
}

// Declare variables in which the values of the flags will be written.
func init() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&reportInterval, "r", 10*time.Second, "Interval to report metrics")
	flag.DurationVar(&pollInterval, "p", 2*time.Second, "Interval to poll metrics")
	flag.StringVar(&hashKey, "k", "", "Key to encrypt metrics")
	flag.IntVar(&rateLimit, "l", 1, "Limit the number of requests to the server")
	flag.Func("crypto-key", "Crypto key for asymmetric encryption", func(input string) error {
		parseResult, err := rsaPublicKeyParser(input)
		if err != nil {
			return err
		}

		cryptoKey = parseResult.(*rsa.PublicKey)
		return nil
	})
}

// Parses agent.Config from environment variables or flags.
func parseConfig() (*agent.Config, error) {
	flag.Parse()

	config := agent.NewConfig(address, reportInterval, pollInterval, hashKey, rateLimit, cryptoKey)

	opts := env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(rsa.PublicKey{}): rsaPublicKeyParser,
		},
	}

	if err := env.ParseWithOptions(config, opts); err != nil {
		return nil, err
	}
	return config, nil
}

func printBuildProperties() {
	log.Printf("Build version: %s", buildVersion)
	log.Printf("Build date: %s", buildDate)
	log.Printf("Build commit: %s", buildCommit)
}
