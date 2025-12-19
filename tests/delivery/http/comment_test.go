package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	httpdelivery "github.com/vdgalyns/blog-api/internal/delivery/http"
	"github.com/vdgalyns/blog-api/internal/domain"
)

// MockCommentUseCase для тестирования COMMENT обработчика
type MockCommentUseCase struct {
	createFunc      func(ctx context.Context, postID int, content string) (*domain.Comment, error)
	getByPostIDFunc func(ctx context.Context, postID int) ([]*domain.Comment, error)
}

func (m *MockCommentUseCase) Create(ctx context.Context, postID int, content string) (*domain.Comment, error) {
	return m.createFunc(ctx, postID, content)
}

func (m *MockCommentUseCase) GetByPostID(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return m.getByPostIDFunc(ctx, postID)
}

func TestCommentHandler_Create_Success(t *testing.T) {
	mockUseCase := &MockCommentUseCase{
		createFunc: func(ctx context.Context, postID int, content string) (*domain.Comment, error) {
			return &domain.Comment{
				ID:        1,
				PostID:    postID,
				Content:   content,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := httpdelivery.NewCommentHandler(mockUseCase)

	type createCommentRequest struct {
		Content string `json:"content"`
	}

	body := createCommentRequest{
		Content: "Test Comment",
	}
	bodyBytes, _ := json.Marshal(body)

	router := chi.NewRouter()
	router.Post("/{postID}/comments", handler.Create)

	req := httptest.NewRequest("POST", "/1/comments", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	type response struct {
		Success bool
	}
	var resp response
	json.NewDecoder(w.Body).Decode(&resp)

	if !resp.Success {
		t.Error("Expected response success to be true")
	}
}

func TestCommentHandler_Create_InvalidJSON(t *testing.T) {
	mockUseCase := &MockCommentUseCase{
		createFunc: func(ctx context.Context, postID int, content string) (*domain.Comment, error) {
			return nil, nil
		},
	}

	handler := httpdelivery.NewCommentHandler(mockUseCase)

	router := chi.NewRouter()
	router.Post("/{postID}/comments", handler.Create)

	req := httptest.NewRequest("POST", "/1/comments", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCommentHandler_Create_InvalidPostID(t *testing.T) {
	mockUseCase := &MockCommentUseCase{
		createFunc: func(ctx context.Context, postID int, content string) (*domain.Comment, error) {
			return nil, nil
		},
	}

	handler := httpdelivery.NewCommentHandler(mockUseCase)

	type createCommentRequest struct {
		Content string `json:"content"`
	}

	body := createCommentRequest{
		Content: "Test Comment",
	}
	bodyBytes, _ := json.Marshal(body)

	router := chi.NewRouter()
	router.Post("/{postID}/comments", handler.Create)

	req := httptest.NewRequest("POST", "/invalid/comments", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCommentHandler_Create_ValidationError(t *testing.T) {
	mockUseCase := &MockCommentUseCase{
		createFunc: func(ctx context.Context, postID int, content string) (*domain.Comment, error) {
			return nil, domain.ErrCommentTooLong
		},
	}

	handler := httpdelivery.NewCommentHandler(mockUseCase)

	type createCommentRequest struct {
		Content string `json:"content"`
	}

	body := createCommentRequest{
		Content: string(make([]byte, 1001)),
	}
	bodyBytes, _ := json.Marshal(body)

	router := chi.NewRouter()
	router.Post("/{postID}/comments", handler.Create)

	req := httptest.NewRequest("POST", "/1/comments", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	type response struct {
		Success bool
	}
	var resp response
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Success {
		t.Error("Expected response success to be false")
	}
}

func TestCommentHandler_GetByPostID_Success(t *testing.T) {
	expectedComments := []*domain.Comment{
		{ID: 1, PostID: 1, Content: "Comment 1", CreatedAt: time.Now()},
		{ID: 2, PostID: 1, Content: "Comment 2", CreatedAt: time.Now()},
	}

	mockUseCase := &MockCommentUseCase{
		getByPostIDFunc: func(ctx context.Context, postID int) ([]*domain.Comment, error) {
			if postID == 1 {
				return expectedComments, nil
			}
			return []*domain.Comment{}, nil
		},
	}

	handler := httpdelivery.NewCommentHandler(mockUseCase)

	router := chi.NewRouter()
	router.Get("/{postID}/comments", handler.GetByPostID)

	req := httptest.NewRequest("GET", "/1/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	type response struct {
		Success bool
	}
	var resp response
	json.NewDecoder(w.Body).Decode(&resp)

	if !resp.Success {
		t.Error("Expected response success to be true")
	}
}

func TestCommentHandler_GetByPostID_Empty(t *testing.T) {
	mockUseCase := &MockCommentUseCase{
		getByPostIDFunc: func(ctx context.Context, postID int) ([]*domain.Comment, error) {
			return []*domain.Comment{}, nil
		},
	}

	handler := httpdelivery.NewCommentHandler(mockUseCase)

	router := chi.NewRouter()
	router.Get("/{postID}/comments", handler.GetByPostID)

	req := httptest.NewRequest("GET", "/999/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCommentHandler_GetByPostID_InvalidPostID(t *testing.T) {
	mockUseCase := &MockCommentUseCase{
		getByPostIDFunc: func(ctx context.Context, postID int) ([]*domain.Comment, error) {
			return nil, nil
		},
	}

	handler := httpdelivery.NewCommentHandler(mockUseCase)

	router := chi.NewRouter()
	router.Get("/{postID}/comments", handler.GetByPostID)

	req := httptest.NewRequest("GET", "/invalid/comments", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}
