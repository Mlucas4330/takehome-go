package domain

type ErrorResponse struct {
	Error string `json:"error" example:"Mensagem de erro"`
}

type ErrIDInvalido struct {
	Error string `json:"error" example:"ID inválido"`
}

type ErrColaboradorNaoEncontrado struct {
	Error string `json:"error" example:"Colaborador não encontrado"`
}

type ErrDepartamentoNaoEncontrado struct {
	Error string `json:"error" example:"Departamento não encontrado"`
}

type ErrCPFInvalido struct {
	Error string `json:"error" example:"CPF inválido"`
}

type ErrCPFJaCadastrado struct {
	Error string `json:"error" example:"CPF já cadastrado"`
}

type ErrRGJaCadastrado struct {
	Error string `json:"error" example:"RG já cadastrado"`
}

type ErrDepartamentoSuperiorInvalido struct {
	Error string `json:"error" example:"departamento superior inválido"`
}

type ErrCicloHierarquia struct {
	Error string `json:"error" example:"mudança/criação causaria ciclo na hierarquia"`
}

type ErrPayloadInvalido struct {
	Error string `json:"error" example:"json: cannot unmarshal"`
}
