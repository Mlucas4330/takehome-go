package models

import (
	"net/http"
	"time"
	"unicode"

	"github.com/mlucas4330/takehome-go/internal/application"
)

type Collaborator struct {
	ID             string    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Nome           string    `gorm:"not null" json:"nome"`
	CPF            string    `gorm:"type:varchar(11);uniqueIndex;not null" json:"cpf"`
	RG             *string   `gorm:"uniqueIndex" json:"rg,omitempty"`
	DepartamentoID string    `gorm:"type:uuid;not null" json:"departamento_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func NewCollaborator(nome, cpf string, rg *string, depID string) (*Collaborator, error) {
	c := &Collaborator{Nome: nome, CPF: cpf, RG: rg, DepartamentoID: depID}
	if err := c.ValidateCPF(); err != nil {
		return nil, err
	}
	if err := c.ValidateRG(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Collaborator) ValidateCPF() error {
	d := make([]int, 0, 11)
	for _, r := range c.CPF {
		if unicode.IsDigit(r) {
			d = append(d, int(r-'0'))
		}
	}

	if len(d) != 11 {
		return application.NewError(
			http.StatusUnprocessableEntity,
			"invalid_cpf",
			"CPF deve conter 11 dígitos",
			map[string]any{"field": "cpf"},
		)
	}

	allEq := true
	for i := 1; i < 11; i++ {
		if d[i] != d[0] {
			allEq = false
			break
		}
	}
	if allEq {
		return application.NewError(
			http.StatusUnprocessableEntity,
			"invalid_cpf",
			"CPF inválido (todos os dígitos são iguais)",
			map[string]any{"field": "cpf"},
		)
	}

	sum := 0
	for i := 0; i < 9; i++ {
		sum += d[i] * (10 - i)
	}
	dv1 := (sum * 10) % 11
	if dv1 == 10 {
		dv1 = 0
	}
	if dv1 != d[9] {
		return application.NewError(
			http.StatusUnprocessableEntity,
			"invalid_cpf",
			"CPF inválido (primeiro dígito verificador)",
			map[string]any{"field": "cpf"},
		)
	}

	sum = 0
	for i := 0; i < 10; i++ {
		sum += d[i] * (11 - i)
	}
	dv2 := (sum * 10) % 11
	if dv2 == 10 {
		dv2 = 0
	}
	if dv2 != d[10] {
		return application.NewError(
			http.StatusUnprocessableEntity,
			"invalid_cpf",
			"CPF inválido (segundo dígito verificador)",
			map[string]any{"field": "cpf"},
		)
	}

	return nil
}

func (c *Collaborator) ValidateRG() error {
	if c.RG == nil {
		return nil
	}

	d := make([]int, 0, 12)
	for _, r := range *c.RG {
		if unicode.IsDigit(r) {
			d = append(d, int(r-'0'))
		}
	}

	if len(d) < 7 || len(d) > 9 {
		return application.NewError(
			http.StatusUnprocessableEntity,
			"invalid_rg",
			"RG deve conter entre 7 e 9 dígitos",
			map[string]any{"field": "rg"},
		)
	}

	allEq := true
	for i := 1; i < len(d); i++ {
		if d[i] != d[0] {
			allEq = false
			break
		}
	}
	if allEq {
		return application.NewError(
			http.StatusUnprocessableEntity,
			"invalid_rg",
			"RG inválido (todos os dígitos são iguais)",
			map[string]any{"field": "rg"},
		)
	}

	return nil
}
