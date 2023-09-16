package server

import (
	"time"
)

// Config describes parameters required for Server.
type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	HashKey       string        `env:"KEY"`
	DatabaseDSN   string        `env:"DATABASE_DSN"`
}

// NewConfig returns new instance of Config with given parameters.
func NewConfig(address string, storeInterval time.Duration, storeFile string, restore bool, hashKey string, databaseDSN string) *Config {
	return &Config{
		Address:       address,
		StoreInterval: storeInterval,
		StoreFile:     storeFile,
		Restore:       restore,
		HashKey:       hashKey,
		DatabaseDSN:   databaseDSN,
	}
}
