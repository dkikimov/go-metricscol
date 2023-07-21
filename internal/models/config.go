package models

type Config struct {
	CryptoKey string `env:"KEY"`
}

func NewConfig(cryptoKey string) *Config {
	return &Config{CryptoKey: cryptoKey}
}
