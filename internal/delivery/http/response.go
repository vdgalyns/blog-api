package http

import (
	"encoding/json"
	"net/http"

	"github.com/radio-pool/backend/internal/domain"
)

type Response struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset,omitempty"`
}

func respondJSON(w http.ResponseWriter, status int, data interface{}, pagination *Pagination) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success:    status >= 200 && status <= 299,
		Data:       data,
		Pagination: pagination,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Success: false,
		Error:   message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func isValidationError(err error) bool {
	if err == nil {
		return false
	}

	return err == domain.ErrInvalidInput ||
		err == domain.ErrTitleTooShort ||
		err == domain.ErrTitleTooLong ||
		err == domain.ErrContentTooShort ||
		err == domain.ErrCommentTooShort ||
		err == domain.ErrCommentTooLong
}
