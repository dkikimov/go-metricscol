package agent

import "time"

type Config struct {
	Address        string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

func NewConfig(address string, reportInterval time.Duration, pollInterval time.Duration) *Config {
	return &Config{Address: address, ReportInterval: reportInterval, PollInterval: pollInterval}
}
