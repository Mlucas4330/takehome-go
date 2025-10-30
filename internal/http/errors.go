package http

import (
	"errors"
	"net/http"
	"takehome-go/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": AppError{Code: "not_found", Message: err.Error()}})
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": AppError{Code: "unique_conflict", Message: err.Error()}})
	case errors.Is(err, service.ErrValidation):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": AppError{Code: "validation_error", Message: err.Error()}})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": AppError{Code: "bad_request", Message: err.Error()}})
	}
}
