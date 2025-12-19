package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vdgalyns/blog-api/internal/domain"
)

func NewRouter(
	postUseCase domain.PostUseCase,
	commentUseCase domain.CommentUseCase,
) chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second))

	postHandler := NewPostHandler(postUseCase)
	commentHandler := NewCommentHandler(commentUseCase)

	router.Route("/api", func(router chi.Router) {
		router.Route("/posts", func(router chi.Router) {
			router.Post("/", postHandler.Create)
			router.Get("/", postHandler.GetAll)
			router.Get("/{id}", postHandler.GetByID)
			router.Put("/{id}", postHandler.Update)
			router.Delete("/{id}", postHandler.Delete)

			router.Route("/{postID}/comments", func(router chi.Router) {
				router.Post("/", commentHandler.Create)
				router.Get("/", commentHandler.GetByPostID)
			})
		})
	})

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return router
}
