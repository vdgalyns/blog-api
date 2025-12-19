package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vdgalyns/blog-api/internal/domain"
)

type PostHandler struct {
	useCase domain.PostUseCase
}

func NewPostHandler(useCase domain.PostUseCase) *PostHandler {
	return &PostHandler{
		useCase: useCase,
	}
}

type createPostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.useCase.Create(r.Context(), req.Title, req.Content)
	if err != nil {
		if isValidationError(err) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create post")
		return
	}

	respondJSON(w, http.StatusCreated, post, nil)
}

func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.useCase.GetByID(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "post not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to get post")
		return
	}

	respondJSON(w, http.StatusOK, post, nil)
}

func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.useCase.GetAll(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get posts")
		return
	}

	respondJSON(w, http.StatusOK, posts, nil)
}

type updatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	var req updatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	post, err := h.useCase.Update(r.Context(), id, req.Title, req.Content)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "post not found")
			return
		}
		if isValidationError(err) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to update post")
		return
	}

	respondJSON(w, http.StatusOK, post, nil)
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	err = h.useCase.Delete(r.Context(), id)
	if err != nil {
		if err == domain.ErrNotFound {
			respondError(w, http.StatusNotFound, "post not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to delete post")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "post deleted"}, nil)
}
