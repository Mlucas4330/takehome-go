package httputil

import (
	"net/http"

	"github.com/mlucas4330/takehome-go/internal/apperror"

	"github.com/gin-gonic/gin"
)

func WriteError(c *gin.Context, err error) {
	if ae, ok := err.(*apperror.Error); ok {
		c.JSON(ae.Status, gin.H{
			"code":    ae.Code,
			"message": ae.Message,
			"details": ae.Details,
		})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"code":    "internal_error",
		"message": "Erro interno do servidor",
	})
}
