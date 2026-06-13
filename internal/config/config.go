package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	ServerAddr      string
	BaseURL         string
	FileStoragePath string
}

type envVars struct {
	ServerAddr      string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}

func New() *Config {
	cfg := &Config{}
	flag.StringVar(&cfg.ServerAddr, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&cfg.BaseURL, "b", "http://localhost:8080", "Base URL for short links")
	flag.StringVar(&cfg.FileStoragePath, "f", "storage.json", "File storage path")
	flag.Parse()

	e := &envVars{}
	if err := env.Parse(e); err != nil {
		log.Fatalf("config error: %v", err)
	}
	if e.ServerAddr != "" {
		cfg.ServerAddr = e.ServerAddr
	}
	if e.BaseURL != "" {
		cfg.BaseURL = e.BaseURL
	}
	if e.FileStoragePath != "" {
		cfg.FileStoragePath = e.FileStoragePath
	}

	return cfg
}
