package app

import (
	"log"

	"github.com/vdgalyns/blog-api/internal/config"
)

func Run() {
	cfg, err := config.New()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Printf("Starting server on port %s", cfg.HTTP.Port)
}