package dto

type CreateColaboradorRequest struct {
	Nome           string  `json:"nome" binding:"required"`
	CPF            string  `json:"cpf" binding:"required"`
	RG             *string `json:"rg"`
	DepartamentoID string  `json:"departamento_id" binding:"required"`
}

type UpdateColaboradorRequest struct {
	Nome           string  `json:"nome"`
	CPF            string  `json:"cpf"`
	RG             *string `json:"rg"`
	DepartamentoID string  `json:"departamento_id"`
}

type ListColaboradorRequest struct {
	Nome           *string `json:"nome"`
	CPF            *string `json:"cpf"`
	RG             *string `json:"rg"`
	DepartamentoID *string `json:"departamento_id"`
	Page           int     `json:"page"`
	PageSize       int     `json:"page_size"`
}