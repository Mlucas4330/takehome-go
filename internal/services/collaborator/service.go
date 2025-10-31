package collaborator

import (
	"context"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mlucas4330/takehome-go/internal/apperror"
	"github.com/mlucas4330/takehome-go/internal/dtos"
	"github.com/mlucas4330/takehome-go/internal/models"
	"github.com/mlucas4330/takehome-go/internal/repositories"
)

type CollaboratorService struct {
	colRepo  repositories.CollaboratorRepository
	deptRepo repositories.DepartamentRepository
}

func toCreateResponseDTO(col *models.Collaborator) dtos.CreateCollaboratorResponse {
	return dtos.CreateCollaboratorResponse{
		Nome:           col.Nome,
		CPF:            col.CPF,
		RG:             col.RG,
		DepartamentoID: col.DepartamentoID,
	}
}

func NewCollaboratorService(colRepo repositories.CollaboratorRepository, deptRepo repositories.DepartamentRepository) *CollaboratorService {
	return &CollaboratorService{colRepo: colRepo, deptRepo: deptRepo}
}

func (s *CollaboratorService) Create(ctx context.Context, req dtos.CreateCollaboratorRequest) (dtos.CreateCollaboratorResponse, error) {
	col, err := models.NewCollaborator(req.Nome, req.CPF, req.RG, req.DepartamentoID)
	if err != nil {
		return dtos.CreateCollaboratorResponse{}, err
	}

	if err := s.colRepo.Create(ctx, col); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, repositories.ErrUniqueViolation) && errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "colaboradores_cpf_unique":
				return dtos.CreateCollaboratorResponse{}, ErrCPFAlreadyExists
			case "colaboradores_rg_unique":
				return dtos.CreateCollaboratorResponse{}, ErrRGAlreadyExists
			default:
				return dtos.CreateCollaboratorResponse{}, ErrCPFOrRGAlreadyExists
			}
		}

		if errors.Is(err, repositories.ErrForeignKey) {
			return dtos.CreateCollaboratorResponse{}, ErrDepartamentNotFound
		}

		return dtos.CreateCollaboratorResponse{}, apperror.New(
			http.StatusInternalServerError,
			"internal_error",
			"Falha ao criar colaborador",
			nil,
		)
	}

	return toCreateResponseDTO(col), nil
}
