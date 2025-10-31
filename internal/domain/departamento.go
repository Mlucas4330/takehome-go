package domain

import (
	"time"

	"github.com/google/uuid"
)

type Departamento struct {
	ID                     uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Nome                   string         `gorm:"not null" json:"nome" binding:"required"`
	GerenteID              uuid.UUID      `gorm:"type:uuid;not null" json:"gerente_id" binding:"required"`
	Gerente                *Colaborador   `gorm:"foreignKey:GerenteID" json:"gerente,omitempty"`
	DepartamentoSuperiorID *uuid.UUID     `gorm:"type:uuid" json:"departamento_superior_id,omitempty"`
	DepartamentoSuperior   *Departamento  `gorm:"foreignKey:DepartamentoSuperiorID" json:"departamento_superior,omitempty"`
	Subdepartamentos       []Departamento `gorm:"foreignKey:DepartamentoSuperiorID" json:"subdepartamentos,omitempty"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
}

type DepartamentoResponse struct {
	ID                     uuid.UUID              `json:"id"`
	Nome                   string                 `json:"nome"`
	GerenteID              uuid.UUID              `json:"gerente_id"`
	NomeGerente            string                 `json:"nome_gerente"`
	DepartamentoSuperiorID *uuid.UUID             `json:"departamento_superior_id,omitempty"`
	Subdepartamentos       []DepartamentoResponse `json:"subdepartamentos,omitempty"`
	CreatedAt              time.Time              `json:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at"`
}

type DepartamentoListRequest struct {
	Nome                   *string    `json:"nome,omitempty"`
	GerenteNome            *string    `json:"gerente_nome,omitempty"`
	DepartamentoSuperiorID *uuid.UUID `json:"departamento_superior_id,omitempty"`
	Page                   int        `json:"page"`
	PageSize               int        `json:"page_size"`
}

func (Departamento) TableName() string {
	return "departamentos"
}

func (d *Departamento) BeforeCreate() error {
	if d.ID == uuid.Nil {
		d.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}
