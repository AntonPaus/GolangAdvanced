package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	DatabaseDsn     string `env:"DATABASE_DSN" envDefault:""`
	LogLevel        string `env:"LOG_LEVEL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
	StoreInterval   uint   `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	address := flag.String("a", "localhost:8080", "server endpoint")
	databaseDsn := flag.String("d", "", "database DSN")
	logLevel := flag.String("l", "info", "log level")
	storeInterval := flag.Uint("i", 300, "store interval in seconds")
	fileStoragePath := flag.String("f", "", "file storage path")
	restore := flag.Bool("r", true, "restore config")
	flag.Parse()
	if cfg.Address == "" {
		cfg.Address = *address
	}
	if cfg.DatabaseDsn == "" {
		cfg.DatabaseDsn = *databaseDsn
	}
	if cfg.FileStoragePath == "" {
		cfg.FileStoragePath = *fileStoragePath
	}
	if cfg.LogLevel == "" {
		cfg.LogLevel = *logLevel
	}
	if cfg.StoreInterval == 0 {
		cfg.StoreInterval = *storeInterval
	}
	if !cfg.Restore {
		cfg.Restore = *restore
	}
	return cfg, nil
}
