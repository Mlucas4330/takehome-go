package collaborator

import (
	"net/http"

	"github.com/mlucas4330/takehome-go/internal/apperror"
)

var (
	ErrCollaboratorNotFound = apperror.New(
		http.StatusNotFound,
		"resource_not_found",
		"Colaborador não encontrado",
		map[string]any{"resource": "colaborador"},
	)

	ErrCPFAlreadyExists = apperror.New(
		http.StatusConflict,
		"duplicate_cpf",
		"CPF já cadastrado no sistema",
		map[string]any{"field": "cpf"},
	)

	ErrRGAlreadyExists = apperror.New(
		http.StatusConflict,
		"duplicate_rg",
		"RG já cadastrado no sistema",
		map[string]any{"field": "rg"},
	)

	ErrCPFOrRGAlreadyExists = apperror.New(
		http.StatusConflict,
		"duplicate_identifier",
		"Já existe um colaborador com um identificador único (CPF ou RG) igual",
		map[string]any{"fields": []string{"cpf", "rg"}},
	)

	ErrDepartamentNotFound = apperror.New(
		http.StatusNotFound,
		"resource_not_found",
		"Departamento não encontrado",
		map[string]any{"field": "departamento_id", "resource": "departamento"},
	)

	ErrManagerNotFound = apperror.New(
		http.StatusNotFound,
		"resource_not_found",
		"Gerente não encontrado",
		map[string]any{"field": "gerente_id", "resource": "colaborador"},
	)

	ErrManagerDifferentDepartament = apperror.New(
		http.StatusUnprocessableEntity,
		"invalid_manager",
		"Gerente deve pertencer ao mesmo departamento do colaborador",
		map[string]any{"field": "gerente_id"},
	)
)
