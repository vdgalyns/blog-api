package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/vdgalyns/blog-api/internal/domain"
	usecase "github.com/vdgalyns/blog-api/internal/usecase"
)

// MockPostRepository для тестирования
type MockPostRepository struct {
	createFunc  func(ctx context.Context, post *domain.Post) error
	getByIDFunc func(ctx context.Context, id int) (*domain.Post, error)
	getAllFunc  func(ctx context.Context) ([]*domain.Post, error)
	updateFunc  func(ctx context.Context, post *domain.Post) error
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *MockPostRepository) Create(ctx context.Context, post *domain.Post) error {
	return m.createFunc(ctx, post)
}

func (m *MockPostRepository) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *MockPostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	return m.getAllFunc(ctx)
}

func (m *MockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	return m.updateFunc(ctx, post)
}

func (m *MockPostRepository) Delete(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}

func TestPostUseCase_Create_Success(t *testing.T) {
	mockRepo := &MockPostRepository{
		createFunc: func(ctx context.Context, post *domain.Post) error {
			post.ID = 1
			return nil
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	post, err := useCase.Create(context.Background(), "Test Title", "Test Content")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if post == nil {
		t.Error("Expected post, got nil")
	}

	if post.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got %s", post.Title)
	}

	if post.Content != "Test Content" {
		t.Errorf("Expected content 'Test Content', got %s", post.Content)
	}
}

func TestPostUseCase_Create_TitleTooShort(t *testing.T) {
	mockRepo := &MockPostRepository{
		createFunc: func(ctx context.Context, post *domain.Post) error {
			return nil
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	post, err := useCase.Create(context.Background(), "AB", "Valid content")

	if err != domain.ErrTitleTooShort {
		t.Errorf("Expected ErrTitleTooShort, got %v", err)
	}

	if post != nil {
		t.Error("Expected post to be nil on validation error")
	}
}

func TestPostUseCase_Create_TitleTooLong(t *testing.T) {
	mockRepo := &MockPostRepository{
		createFunc: func(ctx context.Context, post *domain.Post) error {
			return nil
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	longTitle := string(make([]byte, 256))
	post, err := useCase.Create(context.Background(), longTitle, "Valid content")

	if err != domain.ErrTitleTooLong {
		t.Errorf("Expected ErrTitleTooLong, got %v", err)
	}

	if post != nil {
		t.Error("Expected post to be nil on validation error")
	}
}

func TestPostUseCase_Create_EmptyContent(t *testing.T) {
	mockRepo := &MockPostRepository{
		createFunc: func(ctx context.Context, post *domain.Post) error {
			return nil
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	post, err := useCase.Create(context.Background(), "Valid Title", "")

	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}

	if post != nil {
		t.Error("Expected post to be nil on validation error")
	}
}

func TestPostUseCase_GetByID_Success(t *testing.T) {
	expectedPost := &domain.Post{
		ID:      1,
		Title:   "Test Title",
		Content: "Test Content",
	}

	mockRepo := &MockPostRepository{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Post, error) {
			if id == 1 {
				return expectedPost, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	post, err := useCase.GetByID(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if post == nil {
		t.Error("Expected post, got nil")
	}

	if post.ID != 1 {
		t.Errorf("Expected ID 1, got %d", post.ID)
	}
}

func TestPostUseCase_GetByID_NotFound(t *testing.T) {
	mockRepo := &MockPostRepository{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Post, error) {
			return nil, domain.ErrNotFound
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	post, err := useCase.GetByID(context.Background(), 999)

	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	if post != nil {
		t.Error("Expected post to be nil when not found")
	}
}

func TestPostUseCase_Update_Success(t *testing.T) {
	mockRepo := &MockPostRepository{
		updateFunc: func(ctx context.Context, post *domain.Post) error {
			return nil
		},
		getByIDFunc: func(ctx context.Context, id int) (*domain.Post, error) {
			return &domain.Post{
				ID:      id,
				Title:   "Updated Title",
				Content: "Updated Content",
			}, nil
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	post, err := useCase.Update(context.Background(), 1, "Updated Title", "Updated Content")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if post == nil {
		t.Error("Expected post, got nil")
	}

	if post.Title != "Updated Title" {
		t.Errorf("Expected title 'Updated Title', got %s", post.Title)
	}
}

func TestPostUseCase_Update_NotFound(t *testing.T) {
	mockRepo := &MockPostRepository{
		updateFunc: func(ctx context.Context, post *domain.Post) error {
			return domain.ErrNotFound
		},
		getByIDFunc: func(ctx context.Context, id int) (*domain.Post, error) {
			return nil, domain.ErrNotFound
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	post, err := useCase.Update(context.Background(), 999, "Title", "Content")

	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	if post != nil {
		t.Error("Expected post to be nil when not found")
	}
}

func TestPostUseCase_Delete_Success(t *testing.T) {
	mockRepo := &MockPostRepository{
		deleteFunc: func(ctx context.Context, id int) error {
			if id == 1 {
				return nil
			}
			return domain.ErrNotFound
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	err := useCase.Delete(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestPostUseCase_Delete_NotFound(t *testing.T) {
	mockRepo := &MockPostRepository{
		deleteFunc: func(ctx context.Context, id int) error {
			return domain.ErrNotFound
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	err := useCase.Delete(context.Background(), 999)

	if err != domain.ErrNotFound {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestPostUseCase_GetAll_Success(t *testing.T) {
	expectedPosts := []*domain.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1"},
		{ID: 2, Title: "Post 2", Content: "Content 2"},
	}

	mockRepo := &MockPostRepository{
		getAllFunc: func(ctx context.Context) ([]*domain.Post, error) {
			return expectedPosts, nil
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	posts, err := useCase.GetAll(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}
}

func TestPostUseCase_GetAll_Error(t *testing.T) {
	mockRepo := &MockPostRepository{
		getAllFunc: func(ctx context.Context) ([]*domain.Post, error) {
			return nil, errors.New("database error")
		},
	}

	useCase := usecase.NewPostUseCase(mockRepo)
	posts, err := useCase.GetAll(context.Background())

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if posts != nil {
		t.Error("Expected posts to be nil on error")
	}
}
