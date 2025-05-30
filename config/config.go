package config

type Config struct {
	TreeOrder   int
	StoragePath string
}

func LoadConfig() *Config {
	return &Config{
		TreeOrder:   4,
		StoragePath: "data/tree.db",
	}
}
