package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/vdgalyns/blog-api/internal/domain"
	usecase "github.com/vdgalyns/blog-api/internal/usecase"
)

// MockCommentRepository для тестирования
type MockCommentRepository struct {
	createFunc      func(ctx context.Context, comment *domain.Comment) error
	getByPostIDFunc func(ctx context.Context, postID int) ([]*domain.Comment, error)
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	return m.createFunc(ctx, comment)
}

func (m *MockCommentRepository) GetByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return m.getByPostIDFunc(ctx, postID)
}

func TestCommentUseCase_Create_Success(t *testing.T) {
	mockRepo := &MockCommentRepository{
		createFunc: func(ctx context.Context, comment *domain.Comment) error {
			comment.ID = 1
			return nil
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	comment, err := uc.Create(context.Background(), 1, "Test Comment")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if comment == nil {
		t.Error("Expected comment, got nil")
	}

	if comment.PostID != 1 {
		t.Errorf("Expected post_id 1, got %d", comment.PostID)
	}

	if comment.Content != "Test Comment" {
		t.Errorf("Expected content 'Test Comment', got %s", comment.Content)
	}
}

func TestCommentUseCase_Create_InvalidPostID(t *testing.T) {
	mockRepo := &MockCommentRepository{
		createFunc: func(ctx context.Context, comment *domain.Comment) error {
			return nil
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	comment, err := uc.Create(context.Background(), -1, "Test Comment")

	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}

	if comment != nil {
		t.Error("Expected comment to be nil on validation error")
	}
}

func TestCommentUseCase_Create_EmptyContent(t *testing.T) {
	mockRepo := &MockCommentRepository{
		createFunc: func(ctx context.Context, comment *domain.Comment) error {
			return nil
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	comment, err := uc.Create(context.Background(), 1, "")

	if err != domain.ErrInvalidInput {
		t.Errorf("Expected ErrInvalidInput, got %v", err)
	}

	if comment != nil {
		t.Error("Expected comment to be nil on validation error")
	}
}

func TestCommentUseCase_Create_ContentTooLong(t *testing.T) {
	mockRepo := &MockCommentRepository{
		createFunc: func(ctx context.Context, comment *domain.Comment) error {
			return nil
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	longContent := string(make([]byte, 1001))
	comment, err := uc.Create(context.Background(), 1, longContent)

	if err != domain.ErrCommentTooLong {
		t.Errorf("Expected ErrCommentTooLong, got %v", err)
	}

	if comment != nil {
		t.Error("Expected comment to be nil on validation error")
	}
}

func TestCommentUseCase_GetByPostID_Success(t *testing.T) {
	expectedComments := []*domain.Comment{
		{ID: 1, PostID: 1, Content: "Comment 1"},
		{ID: 2, PostID: 1, Content: "Comment 2"},
	}

	mockRepo := &MockCommentRepository{
		getByPostIDFunc: func(ctx context.Context, postID int) ([]*domain.Comment, error) {
			if postID == 1 {
				return expectedComments, nil
			}
			return []*domain.Comment{}, nil
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	comments, err := uc.GetByPostID(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(comments))
	}
}

func TestCommentUseCase_GetByPostID_Empty(t *testing.T) {
	mockRepo := &MockCommentRepository{
		getByPostIDFunc: func(ctx context.Context, postID int) ([]*domain.Comment, error) {
			return []*domain.Comment{}, nil
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	comments, err := uc.GetByPostID(context.Background(), 999)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(comments) != 0 {
		t.Errorf("Expected 0 comments, got %d", len(comments))
	}
}

func TestCommentUseCase_GetByPostID_Error(t *testing.T) {
	mockRepo := &MockCommentRepository{
		getByPostIDFunc: func(ctx context.Context, postID int) ([]*domain.Comment, error) {
			return nil, errors.New("database error")
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	comments, err := uc.GetByPostID(context.Background(), 1)

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if comments != nil {
		t.Error("Expected comments to be nil on error")
	}
}

func TestCommentUseCase_Create_TrimsWhitespace(t *testing.T) {
	mockRepo := &MockCommentRepository{
		createFunc: func(ctx context.Context, comment *domain.Comment) error {
			comment.ID = 1
			return nil
		},
	}

	uc := usecase.NewCommentUseCase(mockRepo)
	comment, err := uc.Create(context.Background(), 1, "  Test Comment  ")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if comment.Content != "Test Comment" {
		t.Errorf("Expected content 'Test Comment', got '%s'", comment.Content)
	}
}
