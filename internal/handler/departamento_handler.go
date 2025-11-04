package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"takehome-go/internal/dto"
	"takehome-go/internal/service"
)

type DepartamentoHandler struct {
	service service.DepartamentoService
	logger  *zap.Logger
}

func NewDepartamentoHandler(service service.DepartamentoService, logger *zap.Logger) *DepartamentoHandler {
	return &DepartamentoHandler{
		service: service,
		logger:  logger,
	}
}

// Create godoc
// @Summary Criar departamento
// @Description Cria um novo departamento
// @Tags departamentos
// @Accept json
// @Produce json
// @Param departamento body dto.CreateDepartamentoRequest true "Dados do departamento"
// @Success 201 {object} model.Departamento
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /departamentos [post]
func (h *DepartamentoHandler) Create(c *gin.Context) {
	var req dto.CreateDepartamentoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		HandleError(c, http.StatusBadRequest, "Dados inválidos")
		return
	}

	departamento, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "Gerente não encontrado", "Departamento superior não encontrado":
			HandleError(c, http.StatusNotFound, err.Error())
		case "Gerente deve pertencer ao mesmo departamento":
			HandleError(c, http.StatusUnprocessableEntity, err.Error())
		default:
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, departamento)
}

// GetByID godoc
// @Summary Buscar departamento por ID
// @Description Retorna um departamento com sua árvore hierárquica completa
// @Tags departamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do departamento"
// @Success 200 {object} dto.DepartamentoResponse
// @Failure 404 {object} ErrorResponse
// @Router /departamentos/{id} [get]
func (h *DepartamentoHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID", zap.String("id", idStr))
		HandleError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	departamento, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "Departamento não encontrado" {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
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
// @Param departamento body dto.UpdateDepartamentoRequest true "Dados do departamento"
// @Success 200 {object} model.Departamento
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /departamentos/{id} [put]
func (h *DepartamentoHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID", zap.String("id", idStr))
		HandleError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var req dto.UpdateDepartamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		HandleError(c, http.StatusBadRequest, "Dados inválidos")
		return
	}

	departamento, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err.Error() {
		case "Departamento não encontrado", "Gerente não encontrado", "Departamento superior não encontrado":
			HandleError(c, http.StatusNotFound, err.Error())
		case "Gerente deve pertencer ao mesmo departamento", "Operação criaria um ciclo na hierarquia de departamentos":
			HandleError(c, http.StatusUnprocessableEntity, err.Error())
		default:
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, departamento)
}

// Delete godoc
// @Summary Deletar departamento
// @Description Remove um departamento
// @Tags departamentos
// @Accept json
// @Produce json
// @Param id path string true "ID do departamento"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /departamentos/{id} [delete]
func (h *DepartamentoHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID", zap.String("id", idStr))
		HandleError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err.Error() == "Departamento não encontrado" {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
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
// @Param filters body map[string]interface{} false "Filtros (nome, gerente_nome, departamento_superior_id)"
// @Param page query int false "Página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {object} dto.ListDepartamentosResponse
// @Failure 400 {object} ErrorResponse
// @Router /departamentos/listar [post]
func (h *DepartamentoHandler) List(c *gin.Context) {
	var filters map[string]interface{}
	if err := c.ShouldBindJSON(&filters); err != nil {
		filters = make(map[string]interface{})
	}

	page := 1
	pageSize := 10

	if p, ok := filters["page"].(float64); ok {
		page = int(p)
		delete(filters, "page")
	}
	if ps, ok := filters["page_size"].(float64); ok {
		pageSize = int(ps)
		delete(filters, "page_size")
	}

	response, err := h.service.List(c.Request.Context(), filters, page, pageSize)
	if err != nil {
		HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetColaboradoresByGerente godoc
// @Summary Buscar colaboradores por gerente
// @Description Retorna todos os colaboradores dos departamentos subordinados ao gerente
// @Tags gerentes
// @Accept json
// @Produce json
// @Param id path string true "ID do gerente"
// @Success 200 {array} model.Colaborador
// @Failure 404 {object} ErrorResponse
// @Router /gerentes/{id}/colaboradores [get]
func (h *DepartamentoHandler) GetColaboradoresByGerente(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID", zap.String("id", idStr))
		HandleError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	colaboradores, err := h.service.GetColaboradoresByGerente(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "Gerente não encontrado" {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, colaboradores)
}