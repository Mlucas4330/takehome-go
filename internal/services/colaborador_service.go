package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/mlucas4330/takehome-go/internal/domain"
	"github.com/mlucas4330/takehome-go/internal/repositories"
	"github.com/mlucas4330/takehome-go/internal/validators"
	"gorm.io/gorm"
)

type ColaboradorService struct {
	repo     *repositories.ColaboradorRepository
	deptRepo *repositories.DepartamentoRepository
}

func NewColaboradorService(repo *repositories.ColaboradorRepository, deptRepo *repositories.DepartamentoRepository) *ColaboradorService {
	return &ColaboradorService{
		repo:     repo,
		deptRepo: deptRepo,
	}
}

func (s *ColaboradorService) Create(colaborador *domain.Colaborador) error {
	if !validators.ValidateCPF(colaborador.CPF) {
		return errors.New("CPF inválido")
	}

	exists, err := s.repo.ExistsByCPF(colaborador.CPF, nil)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("CPF já cadastrado")
	}

	if colaborador.RG != nil && *colaborador.RG != "" {
		exists, err := s.repo.ExistsByRG(*colaborador.RG, nil)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("RG já cadastrado")
		}
	}

	_, err = s.deptRepo.FindByID(colaborador.DepartamentoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("departamento não encontrado")
		}
		return err
	}

	return s.repo.Create(colaborador)
}

func (s *ColaboradorService) GetByID(id uuid.UUID) (*domain.ColaboradorResponse, error) {
	colaborador, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := &domain.ColaboradorResponse{
		ID:             colaborador.ID,
		Nome:           colaborador.Nome,
		CPF:            colaborador.CPF,
		RG:             colaborador.RG,
		DepartamentoID: colaborador.DepartamentoID,
		CreatedAt:      colaborador.CreatedAt,
		UpdatedAt:      colaborador.UpdatedAt,
	}

	if colaborador.Departamento != nil && colaborador.Departamento.Gerente != nil {
		response.NomeGerente = colaborador.Departamento.Gerente.Nome
	}

	return response, nil
}

func (s *ColaboradorService) Update(id uuid.UUID, colaborador *domain.Colaborador) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if !validators.ValidateCPF(colaborador.CPF) {
		return errors.New("CPF inválido")
	}

	exists, err := s.repo.ExistsByCPF(colaborador.CPF, &id)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("CPF já cadastrado")
	}

	if colaborador.RG != nil && *colaborador.RG != "" {
		exists, err := s.repo.ExistsByRG(*colaborador.RG, &id)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("RG já cadastrado")
		}
	}

	_, err = s.deptRepo.FindByID(colaborador.DepartamentoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("departamento não encontrado")
		}
		return err
	}

	existing.Nome = colaborador.Nome
	existing.CPF = colaborador.CPF
	existing.RG = colaborador.RG
	existing.DepartamentoID = colaborador.DepartamentoID

	return s.repo.Update(existing)
}

func (s *ColaboradorService) Delete(id uuid.UUID) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *ColaboradorService) List(req domain.ColaboradorListRequest) ([]domain.ColaboradorResponse, int64, error) {
	colaboradores, total, err := s.repo.List(req)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]domain.ColaboradorResponse, len(colaboradores))
	for i, c := range colaboradores {
		responses[i] = domain.ColaboradorResponse{
			ID:             c.ID,
			Nome:           c.Nome,
			CPF:            c.CPF,
			RG:             c.RG,
			DepartamentoID: c.DepartamentoID,
			CreatedAt:      c.CreatedAt,
			UpdatedAt:      c.UpdatedAt,
		}
		if c.Departamento != nil && c.Departamento.Gerente != nil {
			responses[i].NomeGerente = c.Departamento.Gerente.Nome
		}
	}

	return responses, total, nil
}
