package model

import "time"

type Colaborator struct {
	ID             string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Nome           string    `gorm:"not null" json:"nome"`
	CPF            string    `gorm:"type:varchar(11);uniqueIndex;not null" json:"cpf"`
	RG             *string   `gorm:"uniqueIndex" json:"rg,omitempty"`
	DepartamentoID string    `gorm:"type:uuid;not null" json:"departamento_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
