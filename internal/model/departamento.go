package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Departamento struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Nome                   string     `gorm:"not null" json:"nome"`
	GerenteID              uuid.UUID  `gorm:"type:uuid;not null" json:"gerente_id"`
	DepartamentoSuperiorID *uuid.UUID `gorm:"type:uuid" json:"departamento_superior_id,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`

	Gerente              *Colaborador   `gorm:"foreignKey:GerenteID" json:"gerente,omitempty"`
	DepartamentoSuperior *Departamento  `gorm:"foreignKey:DepartamentoSuperiorID" json:"departamento_superior,omitempty"`
	Subdepartamentos     []Departamento `gorm:"foreignKey:DepartamentoSuperiorID" json:"subdepartamentos,omitempty"`
}

func (d *Departamento) TableName() string {
	return "departamentos"
}

func (d *Departamento) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}
