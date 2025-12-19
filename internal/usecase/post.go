package usecase

import (
	"context"
	"strings"
	"time"

	"github.com/vdgalyns/blog-api/internal/domain"
)

type PostUseCase struct {
	repository domain.PostRepository
}

func NewPostUseCase(repository domain.PostRepository) *PostUseCase {
	return &PostUseCase{
		repository: repository,
	}
}

func (u *PostUseCase) Create(ctx context.Context, title, content string) (*domain.Post, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if err := validatePost(title, content); err != nil {
		return nil, err
	}

	post := &domain.Post{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := u.repository.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (u *PostUseCase) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	return u.repository.GetByID(ctx, id)
}

func (u *PostUseCase) GetAll(ctx context.Context) ([]*domain.Post, error) {
	return u.repository.GetAll(ctx)
}

func (u *PostUseCase) Update(ctx context.Context, id int, title, content string) (*domain.Post, error) {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if err := validatePost(title, content); err != nil {
		return nil, err
	}

	post := &domain.Post{
		ID:      id,
		Title:   title,
		Content: content,
	}

	if err := u.repository.Update(ctx, post); err != nil {
		return nil, err
	}

	return u.repository.GetByID(ctx, id)
}

func (u *PostUseCase) Delete(ctx context.Context, id int) error {
	return u.repository.Delete(ctx, id)
}

func validatePost(title, content string) error {
	if title == "" {
		return domain.ErrInvalidInput
	}
	if len(title) < 3 {
		return domain.ErrTitleTooShort
	}
	if len(title) > 255 {
		return domain.ErrTitleTooLong
	}
	if content == "" {
		return domain.ErrInvalidInput
	}
	if len(content) < 1 {
		return domain.ErrContentTooShort
	}
	return nil
}
