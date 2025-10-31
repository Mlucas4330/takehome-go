package dtos

type CreateCollaboratorRequest struct {
	Nome           string  `json:"nome" example:"João Silva"`
	CPF            string  `json:"cpf" example:"123.456.789-09"`
	RG             *string `json:"rg,omitempty" example:"12.345.678-9"`
	DepartamentoID string  `json:"departamento_id" example:"b4a3f6a4-6a37-4e1e-9c2f-4a7c6f9e0a1b"`
}

type CreateCollaboratorResponse struct {
	Nome           string  `json:"nome" example:"João Silva"`
	CPF            string  `json:"cpf" example:"123.456.789-09"`
	RG             *string `json:"rg,omitempty" example:"12.345.678-9"`
	DepartamentoID string  `json:"departamento_id" example:"b4a3f6a4-6a37-4e1e-9c2f-4a7c6f9e0a1b"`
}

type UpdateCollaboratorRequest struct {
	Nome           string  `json:"nome"`
	CPF            string  `json:"cpf"`
	RG             *string `json:"rg"`
	DepartamentoID string  `json:"departamento_id"`
}

type ListCollaboratorRequest struct {
	Nome           *string `json:"nome"`
	CPF            *string `json:"cpf"`
	RG             *string `json:"rg"`
	DepartamentoID *string `json:"departamento_id"`
	Page           int     `json:"page"`
	PageSize       int     `json:"page_size"`
}
