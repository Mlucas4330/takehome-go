package dto

import (
	"takehome-go/internal/model"
	"time"

	"github.com/google/uuid"
)

type CreateColaboradorRequest struct {
	Nome           string    `json:"nome" binding:"required"`
	CPF            string    `json:"cpf" binding:"required"`
	RG             *string   `json:"rg"`
	DepartamentoID uuid.UUID `json:"departamento_id" binding:"required"`
}

type UpdateColaboradorRequest struct {
	Nome           string     `json:"nome"`
	CPF            string     `json:"cpf"`
	RG             *string    `json:"rg"`
	DepartamentoID *uuid.UUID `json:"departamento_id"`
}

type ColaboradorResponse struct {
	ID             uuid.UUID `json:"id"`
	Nome           string    `json:"nome"`
	CPF            string    `json:"cpf"`
	RG             *string   `json:"rg,omitempty"`
	DepartamentoID uuid.UUID `json:"departamento_id"`
	NomeGerente    string    `json:"nome_gerente"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ListColaboradoresResponse struct {
	Data       []model.Colaborador `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}
