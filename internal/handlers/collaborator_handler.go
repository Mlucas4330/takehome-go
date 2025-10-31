package handlers

import (
	"net/http"

	"github.com/mlucas4330/takehome-go/internal/apperror"
	"github.com/mlucas4330/takehome-go/internal/dtos"
	"github.com/mlucas4330/takehome-go/internal/httputil"
	"github.com/mlucas4330/takehome-go/internal/services/collaborator"

	"github.com/gin-gonic/gin"
)

type CollaboratorHandler struct {
	colSvc *collaborator.CollaboratorService
}

func NewCollaboratorHandler(colSvc *collaborator.CollaboratorService) *CollaboratorHandler {
	return &CollaboratorHandler{colSvc: colSvc}
}

// Create godoc
// @Summary      Cria um novo colaborador
// @Description  Cria um colaborador validando CPF/RG e existência de departamento/gerente.
// @Tags         Colaboradores
// @Accept       json
// @Produce      json
// @Param        body  body      dtos.CreateCollaboratorRequest  true  "Dados para criação do colaborador"
// @Success      201   {object}  dtos.CreateCollaboratorResponse
// @Failure      400   {object}  httputil.ErrorResponse  "JSON inválido"
// @Failure      404   {object}  httputil.ErrorResponse  "Departamento ou gerente não encontrado"
// @Failure      409   {object}  httputil.ErrorResponse  "CPF ou RG duplicado"
// @Failure      422   {object}  httputil.ErrorResponse  "Validação de domínio (ex.: CPF/RG inválido, gerente de outro departamento)"
// @Failure      500   {object}  httputil.ErrorResponse  "Erro interno do servidor"
// @Router       /api/v1/colaboradores [post]
func (h *CollaboratorHandler) Create(c *gin.Context) {
	var req dtos.CreateCollaboratorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.WriteError(c, apperror.ErrInvalidJSON())
		return
	}

	out, err := h.colSvc.Create(c.Request.Context(), req)
	if err != nil {
		httputil.WriteError(c, err)
		return
	}

	c.JSON(http.StatusCreated, out)
}
