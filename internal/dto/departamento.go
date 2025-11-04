package dto

import (
	"takehome-go/internal/model"
	"time"

	"github.com/google/uuid"
)

type CreateDepartamentoRequest struct {
	Nome                   string     `json:"nome" binding:"required"`
	GerenteID              uuid.UUID  `json:"gerente_id" binding:"required"`
	DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id"`
}

type UpdateDepartamentoRequest struct {
	Nome                   string     `json:"nome"`
	GerenteID              *uuid.UUID `json:"gerente_id"`
	DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id"`
}

type DepartamentoResponse struct {
	ID                     uuid.UUID            `json:"id"`
	Nome                   string               `json:"nome"`
	Gerente                *model.Colaborador   `json:"gerente"`
	DepartamentoSuperiorID *uuid.UUID           `json:"departamento_superior_id,omitempty"`
	Subdepartamentos       []model.Departamento `json:"subdepartamentos"`
	CreatedAt              time.Time            `json:"created_at"`
	UpdatedAt              time.Time            `json:"updated_at"`
}

type ListDepartamentosResponse struct {
	Data       []model.Departamento `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}
