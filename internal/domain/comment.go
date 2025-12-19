package domain

import (
	"context"
	"time"
)

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type CommentRepository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByPostID(ctx context.Context, postID int) ([]*Comment, error)
}

type CommentUseCase interface {
	Create(ctx context.Context, postID int, content string) (*Comment, error)
	GetByPostID(ctx context.Context, postID int) ([]*Comment, error)
}
