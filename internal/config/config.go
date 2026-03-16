package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	StoreInterval   uint   `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	address := flag.String("a", "localhost:8080", "server endpoint")
	storeInterval := flag.Uint("i", 300, "Store interval")
	fileStoragePath := flag.String("f", "./storage", "f")
	restore := flag.Bool("r", false, "restore config")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}
	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = *storeInterval
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStoragePath
	}
	if !cfg.Restore {
		cfg.Restore = *restore
	}
	return &Config{}, nil
}
