package handlers

import "crypto/rsa"

type Config struct {
	HashKey          string
	PrivateCryptoKey *rsa.PrivateKey
}

func NewConfig(hashKey string, privateCryptoKey *rsa.PrivateKey) *Config {
	return &Config{
		HashKey:          hashKey,
		PrivateCryptoKey: privateCryptoKey,
	}
}
