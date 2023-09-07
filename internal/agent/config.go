package agent

import "time"

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	HashKey        string        `env:"KEY"`
	RateLimit      int           `env:"RATE_LIMIT"`
}

func NewConfig(address string, reportInterval time.Duration, pollInterval time.Duration, hashKey string, rateLimit int) *Config {
	return &Config{Address: address, ReportInterval: reportInterval, PollInterval: pollInterval, HashKey: hashKey, RateLimit: rateLimit}
}
