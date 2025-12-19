package main

import (
	"log"

	"github.com/radio-pool/backend/config"
	"github.com/radio-pool/backend/internal/app"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app.Run(cfg)
}
