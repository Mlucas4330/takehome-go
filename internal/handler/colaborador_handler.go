package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"takehome-go/internal/dto"
	"takehome-go/internal/service"
)

type ColaboradorHandler struct {
	service service.ColaboradorService
	logger  *zap.Logger
}

func NewColaboradorHandler(service service.ColaboradorService, logger *zap.Logger) *ColaboradorHandler {
	return &ColaboradorHandler{
		service: service,
		logger:  logger,
	}
}

// Create godoc
// @Summary Criar colaborador
// @Description Cria um novo colaborador
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param colaborador body dto.CreateColaboradorRequest true "Dados do colaborador"
// @Success 201 {object} model.Colaborador
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /colaboradores [post]
func (h *ColaboradorHandler) Create(c *gin.Context) {
	var req dto.CreateColaboradorRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		HandleError(c, http.StatusBadRequest, "Dados inválidos")
		return
	}

	colaborador, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "CPF inválido":
			HandleError(c, http.StatusUnprocessableEntity, err.Error())
		case "CPF já cadastrado", "RG já cadastrado":
			HandleError(c, http.StatusConflict, err.Error())
		case "Departamento não encontrado":
			HandleError(c, http.StatusNotFound, err.Error())
		default:
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusCreated, colaborador)
}

// GetByID godoc
// @Summary Buscar colaborador por ID
// @Description Retorna um colaborador pelo ID com o nome do gerente
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param id path string true "ID do colaborador"
// @Success 200 {object} dto.ColaboradorResponse
// @Failure 404 {object} ErrorResponse
// @Router /colaboradores/{id} [get]
func (h *ColaboradorHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID", zap.String("id", idStr))
		HandleError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	colaborador, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if err.Error() == "Colaborador não encontrado" {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, colaborador)
}

// Update godoc
// @Summary Atualizar colaborador
// @Description Atualiza os dados de um colaborador
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param id path string true "ID do colaborador"
// @Param colaborador body dto.UpdateColaboradorRequest true "Dados do colaborador"
// @Success 200 {object} model.Colaborador
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Router /colaboradores/{id} [put]
func (h *ColaboradorHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID", zap.String("id", idStr))
		HandleError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	var req dto.UpdateColaboradorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		HandleError(c, http.StatusBadRequest, "Dados inválidos")
		return
	}

	colaborador, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err.Error() {
		case "Colaborador não encontrado", "Departamento não encontrado":
			HandleError(c, http.StatusNotFound, err.Error())
		case "CPF inválido":
			HandleError(c, http.StatusUnprocessableEntity, err.Error())
		case "CPF já cadastrado", "RG já cadastrado":
			HandleError(c, http.StatusConflict, err.Error())
		default:
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.JSON(http.StatusOK, colaborador)
}

// Delete godoc
// @Summary Deletar colaborador
// @Description Remove um colaborador
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param id path string true "ID do colaborador"
// @Success 204
// @Failure 404 {object} ErrorResponse
// @Router /colaboradores/{id} [delete]
func (h *ColaboradorHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.logger.Warn("Invalid UUID", zap.String("id", idStr))
		HandleError(c, http.StatusBadRequest, "ID inválido")
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if err.Error() == "Colaborador não encontrado" {
			HandleError(c, http.StatusNotFound, err.Error())
		} else {
			HandleError(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// List godoc
// @Summary Listar colaboradores
// @Description Lista colaboradores com filtros e paginação
// @Tags colaboradores
// @Accept json
// @Produce json
// @Param filters body map[string]interface{} false "Filtros (nome, cpf, rg, departamento_id)"
// @Param page query int false "Página" default(1)
// @Param page_size query int false "Tamanho da página" default(10)
// @Success 200 {object} dto.ListColaboradoresResponse
// @Failure 400 {object} ErrorResponse
// @Router /colaboradores/listar [post]
func (h *ColaboradorHandler) List(c *gin.Context) {
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