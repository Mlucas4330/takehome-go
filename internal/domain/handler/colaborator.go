package handler

import (
	"net/http"
	"takehome-go/internal/domain/dto"
	"takehome-go/internal/domain/model"
	"takehome-go/internal/domain/repository"
	"takehome-go/internal/domain/service"

	httpx "takehome-go/internal/http"

	"github.com/gin-gonic/gin"
)

type ColaboradorHandler struct {
	svc    service.ColaboradorService
	depSvc service.DepartamentoService
}

func NewColaboradorHandler(svc service.ColaboradorService, depSvc service.DepartamentoService) *ColaboradorHandler {
	return &ColaboradorHandler{svc: svc, depSvc: depSvc}
}

// Create
// @Summary Cria colaborador
// @Tags Colaboradores
// @Accept json
// @Produce json
// @Param payload body CreateColaboradorRequest true "Create Colaborator"
// @Success 201 {object} model.Colaborator
// @Failure 422 {object} http.AppError
// @Failure 409 {object} http.AppError
// @Router /api/v1/colaboradores [post]
func (h *ColaboradorHandler) Create(c *gin.Context) {
	var req dto.CreateColaboradorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, err)
		return
	}
	col := &model.Colaborator{
		Nome:           req.Nome,
		CPF:            req.CPF,
		RG:             req.RG,
		DepartamentoID: req.DepartamentoID,
	}
	out, err := h.svc.Create(c.Request.Context(), col)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusCreated, out)
}

// Get
// @Summary Retorna colaborador por ID com nome do gerente do departamento
// @Tags Colaboradores
// @Produce json
// @Param id path string true "Colaborator ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} http.AppError
// @Router /api/v1/colaboradores/{id} [get]
func (h *ColaboradorHandler) Get(c *gin.Context) {
	id := c.Param("id")
	col, gerenteNome, err := h.svc.GetByIDWithGerenteNome(c.Request.Context(), id)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	resp := gin.H{
		"id":              col.ID,
		"nome":            col.Nome,
		"cpf":             col.CPF,
		"rg":              col.RG,
		"departamento_id": col.DepartamentoID,
		"gerente_nome":    gerenteNome,
	}
	c.JSON(http.StatusOK, resp)
}

// Update
// @Summary Atualiza colaborador
// @Tags Colaboradores
// @Accept json
// @Produce json
// @Param id path string true "Colaborator ID"
// @Param payload body UpdateColaboradorRequest true "Update Colaborator"
// @Success 200 {object} model.Colaborator
// @Failure 404 {object} http.AppError
// @Failure 422 {object} http.AppError
// @Failure 409 {object} http.AppError
// @Router /api/v1/colaboradores/{id} [put]
func (h *ColaboradorHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateColaboradorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, err)
		return
	}
	col := &model.Colaborator{
		ID:             id,
		Nome:           req.Nome,
		CPF:            req.CPF,
		RG:             req.RG,
		DepartamentoID: req.DepartamentoID,
	}
	out, err := h.svc.Update(c.Request.Context(), col)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}

// Delete
// @Summary Remove colaborador
// @Tags Colaboradores
// @Param id path string true "Colaborator ID"
// @Success 204
// @Failure 409 {object} http.AppError
// @Router /api/v1/colaboradores/{id} [delete]
func (h *ColaboradorHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

// List
// @Summary Lista colaboradores com filtros e paginação
// @Tags Colaboradores
// @Accept json
// @Produce json
// @Param payload body ListColaboradorRequest true "Filtros"
// @Success 200 {object} repository.PageResult[model.Colaborator]
// @Failure 400 {object} http.AppError
// @Router /api/v1/colaboradores/listar [post]
func (h *ColaboradorHandler) List(c *gin.Context) {
	var req dto.ListColaboradorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, err)
		return
	}
	f := repository.ColabFilter{
		Nome:           req.Nome,
		CPF:            req.CPF,
		RG:             req.RG,
		DepartamentoID: req.DepartamentoID,
	}
	p := repository.Page{Page: req.Page, PageSize: req.PageSize}
	out, err := h.svc.List(c.Request.Context(), f, p)
	if err != nil {
		httpx.WriteError(c, err)
		return
	}
	c.JSON(http.StatusOK, out)
}
