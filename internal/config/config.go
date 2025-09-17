package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTP HTTP
	PG   PG
}

type HTTP struct {
	Port string
}

type PG struct {
	URL string
}

func New() (*Config, error) {
	cfg := &Config{
		HTTP: HTTP{
			Port: os.Getenv("HTTP_PORT"),
		},
		PG: PG{
			URL: os.Getenv("PG_URL"),
		},
	}

	if cfg.HTTP.Port == "" {
		return nil, errors.New("HTTP_PORT is required")
	}

	if cfg.PG.URL == "" {
		return nil, errors.New("PG_URL is required")
	}

	return cfg, nil
}