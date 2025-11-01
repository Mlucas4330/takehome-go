package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mlucas4330/takehome-go/internal/domain"
	"github.com/mlucas4330/takehome-go/internal/services"
)

type ColaboradorHandler struct {
	service *services.ColaboradorService
}

func NewColaboradorHandler(service *services.ColaboradorService) *ColaboradorHandler {
	return &ColaboradorHandler{service: service}
}

// Create godoc
// @Summary Criar colaborador
// @Description Cria um novo colaborador. Valida CPF (único e válido), RG (único se informado) e existência do departamento.
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param payload body domain.Colaborador true "Dados do colaborador"
// @Success 201 {object} domain.Colaborador
// @Failure 400 {object} domain.ErrPayloadInvalido
// @Failure 409 {object} domain.ErrCPFJaCadastrado
// @Failure 409 {object} domain.ErrRGJaCadastrado
// @Failure 422 {object} domain.ErrCPFInvalido
// @Failure 422 {object} domain.ErrDepartamentoNaoEncontrado
// @Router /api/v1/colaboradores [post]
func (h *ColaboradorHandler) Create(c *gin.Context) {
	var colaborador domain.Colaborador
	if err := c.ShouldBindJSON(&colaborador); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&colaborador); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "CPF já cadastrado" || err.Error() == "RG já cadastrado" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, colaborador)
}

// GetByID godoc
// @Summary Buscar colaborador por ID
// @Description Retorna um colaborador e o nome do gerente do seu departamento.
// @Tags colaboradores
// @Produce json
// @Param id path string true "ID do colaborador"
// @Success 200 {object} domain.ColaboradorResponse
// @Failure 400 {object} domain.ErrIDInvalido
// @Failure 404 {object} domain.ErrColaboradorNaoEncontrado
// @Router /api/v1/colaboradores/{id} [get]
func (h *ColaboradorHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	colaborador, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Colaborador não encontrado"})
		return
	}

	c.JSON(http.StatusOK, colaborador)
}

// Update godoc
// @Summary Atualizar colaborador
// @Description Atualiza os dados de um colaborador. Aplica as mesmas validações de criação.
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param id path string true "ID do colaborador"
// @Param payload body domain.Colaborador true "Dados do colaborador"
// @Success 200 {object} domain.Colaborador
// @Failure 400 {object} domain.ErrIDInvalido
// @Failure 400 {object} domain.ErrPayloadInvalido
// @Failure 404 {object} domain.ErrColaboradorNaoEncontrado
// @Failure 409 {object} domain.ErrCPFJaCadastrado
// @Failure 409 {object} domain.ErrRGJaCadastrado
// @Failure 422 {object} domain.ErrCPFInvalido
// @Failure 422 {object} domain.ErrDepartamentoNaoEncontrado
// @Router /api/v1/colaboradores/{id} [put]
func (h *ColaboradorHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var colaborador domain.Colaborador
	if err := c.ShouldBindJSON(&colaborador); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(id, &colaborador); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "Colaborador não encontrado" {
			statusCode = http.StatusNotFound
		}
		if err.Error() == "CPF já cadastrado" || err.Error() == "RG já cadastrado" {
			statusCode = http.StatusConflict
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, colaborador)
}

// Delete godoc
// @Summary Deletar colaborador
// @Description Remove um colaborador.
// @Tags colaboradores
// @Param id path string true "ID do colaborador"
// @Success 204 {string} string "Sem conteúdo"
// @Success 204 {string} string "Sem conteúdo"
// @Failure 400 {object} domain.ErrIDInvalido
// @Failure 404 {object} domain.ErrColaboradorNaoEncontrado
// @Router /api/v1/colaboradores/{id} [delete]
func (h *ColaboradorHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Colaborador não encontrado"})
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary Listar colaboradores
// @Description Lista colaboradores com filtros e paginação. Filtros aceitos: nome, cpf, rg, departamento_id.
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param request body domain.ColaboradorListRequest true "Filtros e paginação"
// @Success 200 {object} domain.ColaboradorPage
// @Failure 400 {object} domain.ErrPayloadInvalido
// @Router /api/v1/colaboradores/listar [post]
func (h *ColaboradorHandler) List(c *gin.Context) {
	var req domain.ColaboradorListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	colaboradores, total, err := h.service.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      colaboradores,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}
