package handlers

import (
	"net/http"

	"github.com/mlucas4330/takehome-go/internal/application"
	"github.com/mlucas4330/takehome-go/internal/dtos"
	"github.com/mlucas4330/takehome-go/internal/httputil"
	"github.com/mlucas4330/takehome-go/internal/services"

	"github.com/gin-gonic/gin"
)

type CollaboratorHandler struct {
	colSvc *services.CollaboratorService
}

func NewCollaboratorHandler(colSvc *services.CollaboratorService) *CollaboratorHandler {
	return &CollaboratorHandler{colSvc: colSvc}
}

func (h *CollaboratorHandler) Create(c *gin.Context) {
	var req dtos.CreateCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.WriteError(c, application.NewError(
			http.StatusBadRequest,
			"invalid_json",
			"JSON malformado",
			nil,
		))
		return
	}

	out, err := h.colSvc.Create(c.Request.Context(), req)
	if err != nil {
		httputil.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}
