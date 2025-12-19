package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/vdgalyns/blog-api/internal/domain"
)

type CommentHandler struct {
	useCase domain.CommentUseCase
}

func NewCommentHandler(useCase domain.CommentUseCase) *CommentHandler {
	return &CommentHandler{
		useCase: useCase,
	}
}

type createCommentRequest struct {
	Content string `json:"content"`
}

func (h *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	var req createCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	comment, err := h.useCase.Create(r.Context(), postID, req.Content)
	if err != nil {
		if isValidationError(err) {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondError(w, http.StatusInternalServerError, "failed to create comment")
		return
	}

	respondJSON(w, http.StatusCreated, comment, nil)
}

func (h *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	comments, err := h.useCase.GetByPostID(r.Context(), postID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get comments")
		return
	}

	respondJSON(w, http.StatusOK, comments, nil)
}
