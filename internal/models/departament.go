package models

import "time"

type Departament struct {
	ID                     string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Nome                   string    `gorm:"not null" json:"nome"`
	GerenteID              string    `gorm:"type:uuid;not null" json:"gerente_id"`
	DepartamentoSuperiorID *string   `gorm:"type:uuid" json:"departamento_superior_id,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
