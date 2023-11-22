package agent

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"go-metricscol/internal/models"
)

var (
	DefaultAddress = "127.0.0.1:8080"
)

// Config describes parameters required for Agent work.
type Config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	HashKey        string
	RateLimit      int
	CryptoKey      *rsa.PublicKey
}

func rsaPublicKeyParser(input string) (*rsa.PublicKey, error) {
	var result *rsa.PublicKey
	if len(input) != 0 {
		cryptoKeyBytes, err := os.ReadFile(input)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode(cryptoKeyBytes)
		result, err = x509.ParsePKCS1PublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// NewConfig returns new instance of Config with given parameters.
func NewConfig(
	address string,
	reportInterval models.Duration,
	pollInterval models.Duration,
	hashKey string,
	rateLimit int,
	cryptoKeyFilePath string,
) (*Config, error) {
	cryptoKey, err := rsaPublicKeyParser(cryptoKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't create config: %s", err)
	}

	return &Config{
		Address:        address,
		ReportInterval: reportInterval.Duration,
		PollInterval:   pollInterval.Duration,
		HashKey:        hashKey,
		RateLimit:      rateLimit,
		CryptoKey:      cryptoKey,
	}, nil
}
