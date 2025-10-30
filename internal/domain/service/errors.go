package service

import "errors"

var (
	ErrValidation = errors.New("validation_error")
	ErrNotFound   = errors.New("not_found")
	ErrConflict   = errors.New("conflict")
)