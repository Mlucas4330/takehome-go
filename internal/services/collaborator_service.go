package services

import (
	"context"
	"log"

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

	dept, err := s.deptRepo.FindByID(col.DepartamentoID)

	log.Fatal(dept)

	if err != nil {
		return dtos.CreateCollaboratorResponse{}, err
	}

	err = s.colRepo.Create(ctx, col)
	if err != nil {
		return dtos.CreateCollaboratorResponse{}, err
	}

	return toCreateResponseDTO(col), nil
}
