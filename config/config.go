package config

import (
	"errors"
	"os"
)

type Config struct {
	HTTP     HTTP
	Postgres Postgres
}

type HTTP struct {
	Port string
}

type Postgres struct {
	DSN string
}

func Load() (*Config, error) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		return nil, errors.New("environment variable POSTGRES_DSN is not set")
	}

	cfg := &Config{
		HTTP: HTTP{
			Port: port,
		},
		Postgres: Postgres{
			DSN: dsn,
		},
	}

	return cfg, nil
}
