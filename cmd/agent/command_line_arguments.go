package main

import "go-metricscol/internal/models"

type commandLineArguments struct {
	Address           string          `json:"address,omitempty" env:"ADDRESS"`
	ReportInterval    models.Duration `json:"report_interval,omitempty" env:"REPORT_INTERVAL"`
	PollInterval      models.Duration `json:"poll_interval,omitempty" env:"POLL_INTERVAL"`
	HashKey           string          `json:"hash_key,omitempty" env:"KEY"`
	RateLimit         int             `json:"rate_limit,omitempty" env:"RATE_LIMIT"`
	CryptoKeyFilePath string          `json:"crypto_key_file_path,omitempty" env:"CRYPTO_KEY"`
	JSONConfigPath    string          `env:"CONFIG"`
}

// Merge writes values of parameter to default same-named values.
func (c *commandLineArguments) Merge(other commandLineArguments) {
	if len(c.Address) == 0 {
		c.Address = other.Address
	}

	if c.ReportInterval.Duration == 0 {
		c.ReportInterval = other.ReportInterval
	}

	if c.PollInterval.Duration == 0 {
		c.PollInterval = other.PollInterval
	}

	if len(c.HashKey) == 0 {
		c.HashKey = other.HashKey
	}

	if c.RateLimit == 0 {
		c.RateLimit = other.RateLimit
	}

	if len(c.CryptoKeyFilePath) == 0 {
		c.CryptoKeyFilePath = other.CryptoKeyFilePath
	}

	if len(c.JSONConfigPath) == 0 {
		c.JSONConfigPath = other.JSONConfigPath
	}
}
