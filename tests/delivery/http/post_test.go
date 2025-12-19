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

// MockPostUseCase для тестирования POST обработчика
type MockPostUseCase struct {
	createFunc  func(ctx context.Context, title, content string) (*domain.Post, error)
	getByIDFunc func(ctx context.Context, id int) (*domain.Post, error)
	getAllFunc  func(ctx context.Context) ([]*domain.Post, error)
	updateFunc  func(ctx context.Context, id int, title, content string) (*domain.Post, error)
	deleteFunc  func(ctx context.Context, id int) error
}

func (m *MockPostUseCase) Create(ctx context.Context, title, content string) (*domain.Post, error) {
	return m.createFunc(ctx, title, content)
}

func (m *MockPostUseCase) GetByID(ctx context.Context, id int) (*domain.Post, error) {
	return m.getByIDFunc(ctx, id)
}

func (m *MockPostUseCase) GetAll(ctx context.Context) ([]*domain.Post, error) {
	return m.getAllFunc(ctx)
}

func (m *MockPostUseCase) Update(ctx context.Context, id int, title, content string) (*domain.Post, error) {
	return m.updateFunc(ctx, id, title, content)
}

func (m *MockPostUseCase) Delete(ctx context.Context, id int) error {
	return m.deleteFunc(ctx, id)
}

func TestPostHandler_Create_Success(t *testing.T) {
	mockUseCase := &MockPostUseCase{
		createFunc: func(ctx context.Context, title, content string) (*domain.Post, error) {
			return &domain.Post{
				ID:        1,
				Title:     title,
				Content:   content,
				CreatedAt: time.Now(),
			}, nil
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	type createPostRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	body := createPostRequest{
		Title:   "Test Title",
		Content: "Test Content",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Create(w, req)

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

func TestPostHandler_Create_InvalidJSON(t *testing.T) {
	mockUseCase := &MockPostUseCase{
		createFunc: func(ctx context.Context, title, content string) (*domain.Post, error) {
			return nil, nil
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPostHandler_Create_ValidationError(t *testing.T) {
	mockUseCase := &MockPostUseCase{
		createFunc: func(ctx context.Context, title, content string) (*domain.Post, error) {
			return nil, domain.ErrTitleTooShort
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	type createPostRequest struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	body := createPostRequest{
		Title:   "AB",
		Content: "Valid Content",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	handler.Create(w, req)

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

func TestPostHandler_GetByID_Success(t *testing.T) {
	expectedPost := &domain.Post{
		ID:        1,
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: time.Now(),
	}

	mockUseCase := &MockPostUseCase{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Post, error) {
			if id == 1 {
				return expectedPost, nil
			}
			return nil, domain.ErrNotFound
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	router := chi.NewRouter()
	router.Get("/{id}", handler.GetByID)

	req := httptest.NewRequest("GET", "/1", nil)
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

func TestPostHandler_GetByID_NotFound(t *testing.T) {
	mockUseCase := &MockPostUseCase{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Post, error) {
			return nil, domain.ErrNotFound
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	router := chi.NewRouter()
	router.Get("/{id}", handler.GetByID)

	req := httptest.NewRequest("GET", "/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
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

func TestPostHandler_GetByID_InvalidID(t *testing.T) {
	mockUseCase := &MockPostUseCase{
		getByIDFunc: func(ctx context.Context, id int) (*domain.Post, error) {
			return nil, nil
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	router := chi.NewRouter()
	router.Get("/{id}", handler.GetByID)

	req := httptest.NewRequest("GET", "/invalid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestPostHandler_GetAll_Success(t *testing.T) {
	expectedPosts := []*domain.Post{
		{ID: 1, Title: "Post 1", Content: "Content 1", CreatedAt: time.Now()},
		{ID: 2, Title: "Post 2", Content: "Content 2", CreatedAt: time.Now()},
	}

	mockUseCase := &MockPostUseCase{
		getAllFunc: func(ctx context.Context) ([]*domain.Post, error) {
			return expectedPosts, nil
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

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

func TestPostHandler_Delete_Success(t *testing.T) {
	mockUseCase := &MockPostUseCase{
		deleteFunc: func(ctx context.Context, id int) error {
			if id == 1 {
				return nil
			}
			return domain.ErrNotFound
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	router := chi.NewRouter()
	router.Delete("/{id}", handler.Delete)

	req := httptest.NewRequest("DELETE", "/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestPostHandler_Delete_NotFound(t *testing.T) {
	mockUseCase := &MockPostUseCase{
		deleteFunc: func(ctx context.Context, id int) error {
			return domain.ErrNotFound
		},
	}

	handler := httpdelivery.NewPostHandler(mockUseCase)

	router := chi.NewRouter()
	router.Delete("/{id}", handler.Delete)

	req := httptest.NewRequest("DELETE", "/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}
