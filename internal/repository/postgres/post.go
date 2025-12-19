package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radio-pool/backend/internal/domain"
)

type PostRepository struct {
	database *pgxpool.Pool
}

func NewPostRepository(database *pgxpool.Pool) *PostRepository {
	return &PostRepository{
		database: database,
	}
}

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	query := "INSERT INTO posts (title, content, created_at) VALUES ($1, $2, $3) RETURNING id"
	err := r.database.QueryRow(ctx, query, post.Title, post.Content, post.CreatedAt).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}
	return nil
}

func (r *PostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	query := "SELECT id, title, content, created_at FROM posts WHERE id = $1"
	post := &domain.Post{}
	err := r.database.QueryRow(ctx, query, id).Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	return post, nil
}

func (r *PostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	query := "SELECT id, title, content, created_at FROM posts ORDER BY created_at DESC"
	rows, err := r.database.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	posts := make([]*domain.Post, 0)
	for rows.Next() {
		post := &domain.Post{}
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan post: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := "UPDATE posts SET title = $1, content = $2 WHERE id = $3"
	tag, err := r.database.Exec(ctx, query, post.Title, post.Content, post.ID)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *PostRepository) Delete(ctx context.Context, id int) error {
	query := "DELETE FROM posts WHERE id = $1"
	tag, err := r.database.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
