package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/radio-pool/backend/config"
	"github.com/radio-pool/backend/internal/delivery/http"
	postgresRepository "github.com/radio-pool/backend/internal/repository/postgres"
	"github.com/radio-pool/backend/internal/usecase"
	"github.com/radio-pool/backend/pkg/httpserver"
	"github.com/radio-pool/backend/pkg/postgres"
)

func Run(cfg *config.Config) {
	// Postgres
	postgresDB, err := postgres.New(cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer postgresDB.Close()

	log.Println("database connected")

	// Repository
	postRepository := postgresRepository.NewPostRepository(postgresDB.Pool)
	commentRepository := postgresRepository.NewCommentRepository(postgresDB.Pool)

	// UseCase
	postUseCase := usecase.NewPostUseCase(postRepository)
	commentUseCase := usecase.NewCommentUseCase(commentRepository)

	// HTTP Server
	httpRouter := http.NewRouter(postUseCase, commentUseCase)
	httpServer := httpserver.New(httpRouter, cfg.HTTP.Port)
	httpServer.Start()

	log.Printf("server listening on :%s", cfg.HTTP.Port)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("shutting down gracefully...")
	case err = <-httpServer.Notify():
		log.Printf("server error: %v", err)
	}

	if err = httpServer.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
