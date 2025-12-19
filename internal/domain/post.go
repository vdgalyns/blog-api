package domain

import (
	"context"
	"time"
)

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type PostRepository interface {
	Create(ctx context.Context, post *Post) error
	GetByID(ctx context.Context, id int) (*Post, error)
	GetAll(ctx context.Context) ([]*Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id int) error
}

type PostUseCase interface {
	Create(ctx context.Context, title, content string) (*Post, error)
	GetByID(ctx context.Context, id int) (*Post, error)
	GetAll(ctx context.Context) ([]*Post, error)
	Update(ctx context.Context, id int, title, content string) (*Post, error)
	Delete(ctx context.Context, id int) error
}
