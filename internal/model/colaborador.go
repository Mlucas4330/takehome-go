package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Colaborador struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Nome           string    `gorm:"not null" json:"nome"`
	CPF            string    `gorm:"uniqueIndex;not null" json:"cpf"`
	RG             *string   `gorm:"uniqueIndex" json:"rg,omitempty"`
	DepartamentoID uuid.UUID `gorm:"type:uuid;not null" json:"departamento_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	Departamento *Departamento `gorm:"foreignKey:DepartamentoID" json:"departamento,omitempty"`
}

func (c *Colaborador) TableName() string {
	return "colaboradores"
}

func (c *Colaborador) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.Must(uuid.NewV7())
	}
	return nil
}