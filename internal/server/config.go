package server

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"
)

// Config describes parameters required for Server.
type Config struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	HashKey       string
	DatabaseDSN   string
	CryptoKey     *rsa.PrivateKey
}

func rsaPrivateKeyParser(input string) (*rsa.PrivateKey, error) {
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

	return result, nil
}

// NewConfig returns new instance of Config with given parameters.
func NewConfig(
	address string,
	storeInterval time.Duration,
	storeFile string,
	restore bool,
	hashKey string,
	databaseDSN string,
	cryptoKeyFilePath string,
) (*Config, error) {
	cryptoKey, err := rsaPrivateKeyParser(cryptoKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't create config: %s", err)
	}

	return &Config{
		Address:       address,
		StoreInterval: storeInterval,
		StoreFile:     storeFile,
		Restore:       restore,
		HashKey:       hashKey,
		DatabaseDSN:   databaseDSN,
		CryptoKey:     cryptoKey,
	}, nil
}
