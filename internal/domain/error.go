package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrValidation      = errors.New("validation error")
	ErrTitleTooShort   = errors.New("title must be at least 3 characters long")
	ErrTitleTooLong    = errors.New("title must not exceed 255 characters")
	ErrContentTooShort = errors.New("content must be at least 1 character long")
	ErrCommentTooShort = errors.New("comment must be at least 1 character long")
	ErrCommentTooLong  = errors.New("comment must not exceed 1000 characters")
)
