package handlers

type Config struct {
	HashKey string
}

func NewConfig(hashKey string) *Config {
	return &Config{HashKey: hashKey}
}
