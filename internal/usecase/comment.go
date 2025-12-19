package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/radio-pool/backend/internal/domain"
)

type CommentUseCase struct {
	repository domain.CommentRepository
}

func NewCommentUseCase(repository domain.CommentRepository) *CommentUseCase {
	return &CommentUseCase{
		repository: repository,
	}
}

func (u *CommentUseCase) Create(ctx context.Context, postID int, content string) (*domain.Comment, error) {
	content = strings.TrimSpace(content)

	if err := validateComment(postID, content); err != nil {
		return nil, err
	}

	comment := &domain.Comment{
		PostID:    postID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := u.repository.Create(ctx, comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (u *CommentUseCase) GetByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return u.repository.GetByPostID(ctx, postID)
}

func validateComment(postID int, content string) error {
	if postID <= 0 {
		return domain.ErrInvalidInput
	}
	if content == "" {
		return domain.ErrInvalidInput
	}
	if len(content) < 1 {
		return domain.ErrCommentTooShort
	}
	if len(content) > 1000 {
		return domain.ErrCommentTooLong
	}
	return nil
}
