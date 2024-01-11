package main

import (
	"go-metricscol/internal/models"
)

type commandLineArguments struct {
	Address           string          `json:"address,omitempty" env:"ADDRESS"`
	StoreInterval     models.Duration `json:"store_interval,omitempty" env:"STORE_INTERVAL"`
	StoreFile         string          `json:"store_file,omitempty" env:"STORE_FILE"`
	Restore           bool            `json:"restore,omitempty" env:"RESTORE"`
	HashKey           string          `json:"hash_key,omitempty" env:"KEY"`
	DatabaseDSN       string          `json:"database_dsn,omitempty" env:"DATABASE_DSN"`
	CryptoKeyFilePath string          `json:"crypto_key_file_path,omitempty" env:"CRYPTO_KEY"`
	TrustedSubnet     string          `json:"trusted_subnet,omitempty" env:"TRUSTED_SUBNET"`
	JSONConfigPath    string          `env:"CONFIG"`
}

// Merge writes values of parameter to default same-named values.
func (c *commandLineArguments) Merge(other commandLineArguments) {
	if len(c.Address) == 0 {
		c.Address = other.Address
	}

	if c.StoreInterval.Duration == 0 {
		c.StoreInterval = other.StoreInterval
	}

	if len(c.StoreFile) == 0 {
		c.StoreFile = other.StoreFile
	}

	if !c.Restore {
		c.Restore = other.Restore
	}

	if len(c.HashKey) == 0 {
		c.HashKey = other.HashKey
	}

	if len(c.DatabaseDSN) == 0 {
		c.DatabaseDSN = other.DatabaseDSN
	}

	if len(c.CryptoKeyFilePath) == 0 {
		c.CryptoKeyFilePath = other.CryptoKeyFilePath
	}

	if len(c.JSONConfigPath) == 0 {
		c.JSONConfigPath = other.JSONConfigPath
	}

	if len(c.TrustedSubnet) == 0 {
		c.TrustedSubnet = other.TrustedSubnet
	}
}
