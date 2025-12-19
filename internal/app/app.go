package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vdgalyns/blog-api/config"
	"github.com/vdgalyns/blog-api/internal/delivery/http"
	postgresRepository "github.com/vdgalyns/blog-api/internal/repository/postgres"
	"github.com/vdgalyns/blog-api/internal/usecase"
	"github.com/vdgalyns/blog-api/pkg/httpserver"
	"github.com/vdgalyns/blog-api/pkg/postgres"
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
