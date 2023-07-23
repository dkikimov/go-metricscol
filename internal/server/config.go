package server

import (
	"time"
)

type Config struct {
	Address       string        `env:"ADDRESS"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	StoreFile     string        `env:"STORE_FILE"`
	Restore       bool          `env:"RESTORE"`
	HashKey       string        `env:"KEY"`
}

func NewConfig(address string, storeInterval time.Duration, storeFile string, restore bool, hashKey string) *Config {
	return &Config{Address: address, StoreInterval: storeInterval, StoreFile: storeFile, Restore: restore, HashKey: hashKey}
}
