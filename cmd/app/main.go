package main

import (
	"log"

	"github.com/vdgalyns/blog-api/config"
	"github.com/vdgalyns/blog-api/internal/app"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	app.Run(cfg)
}
