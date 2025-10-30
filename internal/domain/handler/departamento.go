package handler

import (
	"net/http"

	"takehome-go/internal/domain/model"
	"takehome-go/internal/domain/repository"
	"takehome-go/internal/domain/service"
	httpx "takehome-go/internal/http"

	"github.com/gin-gonic/gin"
)

type DepartamentoHandler struct {
	svc service.DepartamentoService
}

func NewDepartamentoHandler(svc service.DepartamentoService) *DepartamentoHandler {
	return &DepartamentoHandler{svc: svc}
}

type CreateDepartamentoRequest struct {
	Nome                   string  `json:"nome" binding:"required"`
	GerenteID              string  `json:"gerente_id" binding:"required"`
	DepartamentoSuperiorID *string `json:"departamento_superior_id"`
}

type UpdateDepartamentoRequest struct {
	Nome                   string  `json:"nome"`
	GerenteID              string  `json:"gerente_id"`
	DepartamentoSuperiorID *string `json:"departamento_superior_id"`
}

type ListDepartamentoRequest struct {
	Nome                   *string `json:"nome"`
	GerenteNome            *string `json:"gerente_nome"`
	DepartamentoSuperiorID *string `json:"departamento_superior_id"`
	Page                   int     `json:"page"`
	PageSize               int     `json:"page_size"`
}

// Create
// @Summary Cria departamento
// @Tags Departamentos
// @Accept json
// @Produce json
// @Param payload body CreateDepartamentoRequest true "Create Departamento"
// @Success 201 {object} model.Departamento
// @Failure 422 {object} http.AppError
// @Router /api/v1/departamentos [post]
func (h *DepartamentoHandler) Create(c *gin.Context) {
	var req CreateDepartamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, err)
		return
	}
	d := &model.Departamento{
		Nome:                   req.Nome,
		GerenteID:              req.GerenteID,
		DepartamentoSuperiorID: req.DepartamentoSuperiorID,
	}
	out, err := h.svc.Create(c.Request.Context(), d)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Get
// @Summary Retorna departamento, gerente e árvore completa de subdepartamentos
// @Tags Departamentos
// @Produce json
// @Param id path string true "Departamento ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} http.AppError
// @Router /api/v1/departamentos/{id} [get]
func (h *DepartamentoHandler) Get(c *gin.Context) {
	id := c.Param("id")
	dep, gerente, tree, err := h.svc.GetWithTree(c.Request.Context(), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	resp := gin.H{
		"id":        dep.ID,
		"nome":      dep.Nome,
		"gerente":   gerente,
		"subarvore": tree,
	}
	c.JSON(http.StatusOK, resp)
}

// Update
// @Summary Atualiza departamento
// @Tags Departamentos
// @Accept json
// @Produce json
// @Param id path string true "Departamento ID"
// @Param payload body UpdateDepartamentoRequest true "Update Departamento"
// @Success 200 {object} model.Departamento
// @Failure 404 {object} http.AppError
// @Failure 422 {object} http.AppError
// @Router /api/v1/departamentos/{id} [put]
func (h *DepartamentoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateDepartamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, err)
		return
	}
	d := &model.Departamento{
		ID:                     id,
		Nome:                   req.Nome,
		GerenteID:              req.GerenteID,
		DepartamentoSuperiorID: req.DepartamentoSuperiorID,
	}
	out, err := h.svc.Update(c.Request.Context(), d)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete
// @Summary Remove departamento
// @Tags Departamentos
// @Param id path string true "Departamento ID"
// @Success 204
// @Router /api/v1/departamentos/{id} [delete]
func (h *DepartamentoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// List
// @Summary Lista departamentos com filtros e paginação
// @Tags Departamentos
// @Accept json
// @Produce json
// @Param payload body ListDepartamentoRequest true "Filtros"
// @Success 200 {object} repository.PageResult[model.Departamento]
// @Router /api/v1/departamentos/listar [post]
func (h *DepartamentoHandler) List(c *gin.Context) {
	var req ListDepartamentoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, err)
		return
	}
	f := repository.DepFilter{
		Nome:                   req.Nome,
		GerenteNome:            req.GerenteNome,
		DepartamentoSuperiorID: req.DepartamentoSuperiorID,
	}
	p := repository.Page{Page: req.Page, PageSize: req.PageSize}
	out, err := h.svc.List(c.Request.Context(), f, p)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
