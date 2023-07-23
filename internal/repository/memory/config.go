package memory

type Config struct {
	HashKey string
}

func NewConfig(hashKey string) *Config {
	return &Config{HashKey: hashKey}
}
