package httputil

type ErrorResponse struct {
	Code    string `json:"code" example:"invalid_json"`
	Message string `json:"message" example:"JSON malformado ou inválido"`
	Details any    `json:"details,omitempty" example:"{\"field\":\"cpf\"}"`
}
