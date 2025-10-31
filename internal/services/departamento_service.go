package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mlucas4330/takehome-go/internal/domain"
	"github.com/mlucas4330/takehome-go/internal/repositories"
	"gorm.io/gorm"
)

type DepartamentoService struct {
	repo      *repositories.DepartamentoRepository
	colabRepo *repositories.ColaboradorRepository
}

func NewDepartamentoService(repo *repositories.DepartamentoRepository, colabRepo *repositories.ColaboradorRepository) *DepartamentoService {
	return &DepartamentoService{
		repo:      repo,
		colabRepo: colabRepo,
	}
}

func (s *DepartamentoService) Create(departamento *domain.Departamento) error {
	gerente, err := s.colabRepo.FindByID(departamento.GerenteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("gerente não encontrado")
		}
		return err
	}

	if departamento.DepartamentoSuperiorID != nil {
		_, err := s.repo.FindByID(*departamento.DepartamentoSuperiorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("departamento superior não encontrado")
			}
			return err
		}
	}

	err = s.repo.Create(departamento)
	if err != nil {
		return err
	}

	if gerente.DepartamentoID != departamento.ID {
		gerente.DepartamentoID = departamento.ID
		s.colabRepo.Update(gerente)
	}

	return nil
}

func (s *DepartamentoService) GetByID(id uuid.UUID) (*domain.DepartamentoResponse, error) {
	departamento, err := s.repo.FindByIDWithHierarchy(id)
	if err != nil {
		return nil, err
	}

	return s.buildDepartamentoResponse(departamento), nil
}

func (s *DepartamentoService) buildDepartamentoResponse(dept *domain.Departamento) *domain.DepartamentoResponse {
	response := &domain.DepartamentoResponse{
		ID:                     dept.ID,
		Nome:                   dept.Nome,
		GerenteID:              dept.GerenteID,
		DepartamentoSuperiorID: dept.DepartamentoSuperiorID,
		CreatedAt:              dept.CreatedAt,
		UpdatedAt:              dept.UpdatedAt,
	}

	if dept.Gerente != nil {
		response.NomeGerente = dept.Gerente.Nome
	}

	if len(dept.Subdepartamentos) > 0 {
		response.Subdepartamentos = make([]domain.DepartamentoResponse, len(dept.Subdepartamentos))
		for i, sub := range dept.Subdepartamentos {
			response.Subdepartamentos[i] = *s.buildDepartamentoResponse(&sub)
		}
	}

	return response
}

func (s *DepartamentoService) Update(id uuid.UUID, departamento *domain.Departamento) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	gerente, err := s.colabRepo.FindByID(departamento.GerenteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("gerente não encontrado")
		}
		return err
	}

	if gerente.DepartamentoID != id {
		return errors.New("gerente deve pertencer ao departamento")
	}

	if departamento.DepartamentoSuperiorID != nil {
		if *departamento.DepartamentoSuperiorID == id {
			return errors.New("departamento não pode ser superior de si mesmo")
		}

		wouldCycle, err := s.repo.WouldCreateCycle(id, *departamento.DepartamentoSuperiorID)
		if err != nil {
			return err
		}
		if wouldCycle {
			return errors.New("operação criaria ciclo na hierarquia")
		}

		_, err = s.repo.FindByID(*departamento.DepartamentoSuperiorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("departamento superior não encontrado")
			}
			return err
		}
	}

	existing.Nome = departamento.Nome
	existing.GerenteID = departamento.GerenteID
	existing.DepartamentoSuperiorID = departamento.DepartamentoSuperiorID

	return s.repo.Update(existing)
}

func (s *DepartamentoService) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *DepartamentoService) List(req domain.DepartamentoListRequest) ([]domain.DepartamentoResponse, int64, error) {
	departamentos, total, err := s.repo.List(req)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]domain.DepartamentoResponse, len(departamentos))
	for i, d := range departamentos {
		responses[i] = domain.DepartamentoResponse{
			ID:                     d.ID,
			Nome:                   d.Nome,
			GerenteID:              d.GerenteID,
			DepartamentoSuperiorID: d.DepartamentoSuperiorID,
			CreatedAt:              d.CreatedAt,
			UpdatedAt:              d.UpdatedAt,
		}
		if d.Gerente != nil {
			responses[i].NomeGerente = d.Gerente.Nome
		}
	}

	return responses, total, nil
}

func (s *DepartamentoService) GetGerenteColaboradores(gerenteID uuid.UUID) ([]domain.Colaborador, error) {
	_, err := s.colabRepo.FindByID(gerenteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("gerente não encontrado")
		}
		return nil, err
	}

	return s.repo.GetSubordinateColaboradores(gerenteID)
}
