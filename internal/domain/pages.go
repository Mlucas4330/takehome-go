package domain

type ColaboradorPage struct {
	Data     []Colaborador `json:"data"`
	Total    int64         `json:"total" example:"1"`
	Page     int           `json:"page" example:"1"`
	PageSize int           `json:"page_size" example:"10"`
}

type DepartamentoPage struct {
	Data     []Departamento `json:"data"`
	Total    int64          `json:"total" example:"1"`
	Page     int            `json:"page" example:"1"`
	PageSize int            `json:"page_size" example:"10"`
}
