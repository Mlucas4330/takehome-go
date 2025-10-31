package repositories

import "errors"

var (
	ErrNotFound        = errors.New("record not found")
	ErrUniqueViolation = errors.New("unique constraint violation")
	ErrForeignKey      = errors.New("foreign key constraint violation")
	ErrDeadlock        = errors.New("database deadlock")
	ErrConnection      = errors.New("database connection error")
)
