package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"go-metricscol/internal/models"
)

// ServerConfig describes parameters required for Server.
type ServerConfig struct {
	Address       string
	StoreInterval time.Duration
	StoreFile     string
	Restore       bool
	HashKey       string
	DatabaseDSN   string
	CryptoKey     *rsa.PrivateKey
	TrustedSubnet string
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

// NewServerConfig returns new instance of ServerConfig with given parameters.
func NewServerConfig(
	address string,
	storeInterval models.Duration,
	storeFile string,
	restore bool,
	hashKey string,
	databaseDSN string,
	cryptoKeyFilePath string,
	trustedSubnet string,
) (*ServerConfig, error) {
	cryptoKey, err := rsaPrivateKeyParser(cryptoKeyFilePath)
	if err != nil {
		return nil, fmt.Errorf("couldn't create config: %s", err)
	}

	return &ServerConfig{
		Address:       address,
		StoreInterval: storeInterval.Duration,
		StoreFile:     storeFile,
		Restore:       restore,
		HashKey:       hashKey,
		DatabaseDSN:   databaseDSN,
		CryptoKey:     cryptoKey,
		TrustedSubnet: trustedSubnet,
	}, nil
}
