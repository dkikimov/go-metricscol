package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/caarlos0/env/v9"

	"go-metricscol/internal/server"
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

	log.Printf("Starting server on %s", cfg.Address)

	s, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("couldn't create server with error: %s", err)
	}

	log.Fatal(s.ListenAndServe())
}

var (
	address       string
	storeInterval time.Duration
	storeFile     string
	restore       bool
	hashKey       string
	databaseDSN   string
	cryptoKey     *rsa.PrivateKey
)

func rsaPrivateKeyParser(input string) (interface{}, error) {
	var result *rsa.PrivateKey
	if len(input) != 0 {
		cryptoKeyBytes, err := os.ReadFile(input)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(cryptoKeyBytes)
		result, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	return *result, nil
}

// Declare variables in which the values of the flags will be written.
func init() {
	flag.StringVar(&address, "a", "127.0.0.1:8080", "Address to listen")
	flag.DurationVar(&storeInterval, "i", 300*time.Second, "Interval to store metrics")
	flag.StringVar(&storeFile, "f", "/tmp/devops-metrics-db.json", "File to store metrics")
	flag.BoolVar(&restore, "r", true, "Restore metrics from file")
	flag.StringVar(&hashKey, "k", "", "Key to encrypt metrics")
	flag.StringVar(&databaseDSN, "d", "", "Database DSN")
	flag.Func("crypto-key", "Private crypto key for asymmetric encryption", func(input string) error {
		parseResult, err := rsaPrivateKeyParser(input)
		if err != nil {
			return err
		}

		cryptoKey = parseResult.(*rsa.PrivateKey)
		return nil
	})
}

// Parses server.Config from environment variables or flags.
func parseConfig() (*server.Config, error) {
	flag.Parse()
	config := server.NewConfig(address, storeInterval, storeFile, restore, hashKey, databaseDSN, cryptoKey)

	opts := env.Options{
		FuncMap: map[reflect.Type]env.ParserFunc{
			reflect.TypeOf(rsa.PrivateKey{}): rsaPrivateKeyParser,
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
