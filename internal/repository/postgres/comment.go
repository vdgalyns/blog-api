package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radio-pool/backend/internal/domain"
)

type CommentRepository struct {
	database *pgxpool.Pool
}

func NewCommentRepository(database *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		database: database,
	}
}

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := "INSERT INTO comments (post_id, content, created_at) VALUES ($1, $2, $3) RETURNING id"
	err := r.database.QueryRow(ctx, query, comment.PostID, comment.Content, comment.CreatedAt).Scan(&comment.ID)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

func (r *CommentRepository) GetByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	query := "SELECT id, post_id, content, created_at FROM comments WHERE post_id = $1 ORDER BY created_at ASC"
	rows, err := r.database.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}
	defer rows.Close()

	comments := make([]*domain.Comment, 0)
	for rows.Next() {
		c := &domain.Comment{}
		if err := rows.Scan(&c.ID, &c.PostID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, nil
}
