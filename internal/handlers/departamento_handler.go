package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mlucas4330/takehome-go/internal/domain"
	"github.com/mlucas4330/takehome-go/internal/services"
)

type DepartamentoHandler struct {
	service *services.DepartamentoService
}

func NewDepartamentoHandler(service *services.DepartamentoService) *DepartamentoHandler {
	return &DepartamentoHandler{service: service}
}

// Create godoc
// @Summary Criar departamento
// @Description Cria um novo departamento
// @Tags departamentos
// @Accept json
// @Produce json
// @Param departamento body domain.Departamento true "Dados do departamento"
// @Success 201 {object} domain.Departamento
// @Failure 400 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/v1/departamentos [post]
func (h *DepartamentoHandler) Create(c *gin.Context) {
	var departamento domain.Departamento
	if err := c.ShouldBindJSON(&departamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(&departamento); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, departamento)
}

// GetByID godoc
// @Summary Buscar departamento por ID
// @Description Retorna um departamento com sua árvore hierárquica
// @Tags departamentos
// @Produce json
// @Param id path string true "ID do departamento"
// @Success 200 {object} domain.DepartamentoResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/departamentos/{id} [get]
func (h *DepartamentoHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	departamento, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Departamento não encontrado"})
		return
	}

	c.JSON(http.StatusOK, departamento)
}

// Update godoc
// @Summary Atualizar departamento
// @Description Atualiza os dados de um departamento
// @Tags departamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do departamento"
// @Param departamento body domain.Departamento true "Dados do departamento"
// @Success 200 {object} domain.Departamento
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 422 {object} map[string]string
// @Router /api/v1/departamentos/{id} [put]
func (h *DepartamentoHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var departamento domain.Departamento
	if err := c.ShouldBindJSON(&departamento); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(id, &departamento); err != nil {
		statusCode := http.StatusUnprocessableEntity
		if err.Error() == "Departamento não encontrado" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, departamento)
}

// Delete godoc
// @Summary Deletar departamento
// @Description Remove um departamento
// @Tags departamentos
// @Param id path string true "ID do departamento"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /api/v1/departamentos/{id} [delete]
func (h *DepartamentoHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Departamento não encontrado"})
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary Listar departamentos
// @Description Lista departamentos com filtros e paginação
// @Tags departamentos
// @Accept json
// @Produce json
// @Param request body domain.DepartamentoListRequest true "Filtros e paginação"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /api/v1/departamentos/listar [post]
func (h *DepartamentoHandler) List(c *gin.Context) {
	var req domain.DepartamentoListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	departamentos, total, err := h.service.List(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      departamentos,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// GetGerenteColaboradores godoc
// @Summary Listar colaboradores subordinados ao gerente
// @Description Retorna todos os colaboradores dos departamentos subordinados ao gerente
// @Tags gerentes
// @Produce json
// @Param id path string true "ID do gerente"
// @Success 200 {array} domain.Colaborador
// @Failure 404 {object} map[string]string
// @Router /api/v1/gerentes/{id}/colaboradores [get]
func (h *DepartamentoHandler) GetGerenteColaboradores(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	colaboradores, err := h.service.GetGerenteColaboradores(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, colaboradores)
}
