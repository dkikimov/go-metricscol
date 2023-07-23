package models

type Config struct {
	CryptoKey string
}

func NewConfig(cryptoKey string) *Config {
	return &Config{CryptoKey: cryptoKey}
}
