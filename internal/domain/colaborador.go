package domain

import (
	"time"

	"github.com/google/uuid"
)

type Colaborador struct {
	ID             uuid.UUID     `gorm:"type:uuid;primary_key" json:"id"`
	Nome           string        `gorm:"not null" json:"nome" binding:"required"`
	CPF            string        `gorm:"uniqueIndex;not null" json:"cpf" binding:"required"`
	RG             *string       `gorm:"uniqueIndex" json:"rg,omitempty"`
	DepartamentoID uuid.UUID     `gorm:"type:uuid;not null" json:"departamento_id" binding:"required"`
	Departamento   *Departamento `gorm:"foreignKey:DepartamentoID" json:"departamento,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type ColaboradorResponse struct {
	Nome           string    `json:"nome"`
	CPF            string    `json:"cpf"`
	RG             *string   `json:"rg,omitempty"`
	DepartamentoID uuid.UUID `json:"departamento_id"`
}

type ColaboradorListRequest struct {
	Nome           *string    `json:"nome,omitempty"`
	CPF            *string    `json:"cpf,omitempty"`
	RG             *string    `json:"rg,omitempty"`
	DepartamentoID *uuid.UUID `json:"departamento_id,omitempty"`
	Page           int        `json:"page"`
	PageSize       int        `json:"page_size"`
}

func (Colaborador) TableName() string {
	return "colaboradores"
}

func (c *Colaborador) BeforeCreate() error {
	if c.ID == uuid.Nil {
		c.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}
