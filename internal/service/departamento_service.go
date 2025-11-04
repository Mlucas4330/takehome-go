package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"takehome-go/internal/database"
	"takehome-go/internal/dto"
	"takehome-go/internal/model"
	"takehome-go/internal/repository"
)

type DepartamentoService interface {
	Create(ctx context.Context, req *dto.CreateDepartamentoRequest) (*model.Departamento, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.DepartamentoResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDepartamentoRequest) (*model.Departamento, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) (*dto.ListDepartamentosResponse, error)
	GetColaboradoresByGerente(ctx context.Context, gerenteID uuid.UUID) ([]model.Colaborador, error)
}

type departamentoService struct {
	repo      repository.DepartamentoRepository
	colabRepo repository.ColaboradorRepository
	cache     database.Cache
	logger    *zap.Logger
}

func NewDepartamentoService(
	repo repository.DepartamentoRepository,
	colabRepo repository.ColaboradorRepository,
	cache database.Cache,
	logger *zap.Logger,
) DepartamentoService {
	return &departamentoService{
		repo:      repo,
		colabRepo: colabRepo,
		cache:     cache,
		logger:    logger,
	}
}

func (s *departamentoService) Create(ctx context.Context, req *dto.CreateDepartamentoRequest) (*model.Departamento, error) {
	s.logger.Info("Creating departamento", zap.String("nome", req.Nome))

	gerente, err := s.colabRepo.GetByID(ctx, req.GerenteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Gerente not found", zap.String("gerente_id", req.GerenteID.String()))
			return nil, errors.New("Gerente não encontrado")
		}
		s.logger.Error("Failed to get gerente", zap.Error(err))
		return nil, errors.New("Erro ao buscar gerente")
	}

	if req.DepartamentoSuperiorID != nil {
		_, err := s.repo.GetByID(ctx, *req.DepartamentoSuperiorID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Warn("Superior department not found", zap.String("departamento_superior_id", req.DepartamentoSuperiorID.String()))
				return nil, errors.New("Departamento superior não encontrado")
			}
			s.logger.Error("Failed to get superior department", zap.Error(err))
			return nil, errors.New("Erro ao buscar departamento superior")
		}
	}

	departamento := &model.Departamento{
		Nome:                   req.Nome,
		GerenteID:              req.GerenteID,
		DepartamentoSuperiorID: req.DepartamentoSuperiorID,
	}

	if err := s.repo.Create(ctx, departamento); err != nil {
		s.logger.Error("Failed to create departamento", zap.Error(err))
		return nil, errors.New("Erro ao criar departamento")
	}

	if gerente.DepartamentoID != departamento.ID {
		gerente.DepartamentoID = departamento.ID
		if err := s.colabRepo.Update(ctx, gerente); err != nil {
			s.logger.Error("Failed to update gerente department", zap.Error(err))
		}
	}

	s.logger.Info("Departamento created successfully", zap.String("id", departamento.ID.String()))
	return departamento, nil
}

func (s *departamentoService) GetByID(ctx context.Context, id uuid.UUID) (*dto.DepartamentoResponse, error) {
	s.logger.Info("Getting departamento by ID", zap.String("id", id.String()))

	cacheKey := fmt.Sprintf("departamento:%s", id.String())
	var cached dto.DepartamentoResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("Departamento found in cache", zap.String("id", id.String()))
		return &cached, nil
	}

	departamento, err := s.repo.GetByIDWithHierarchy(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Departamento not found", zap.String("id", id.String()))
			return nil, errors.New("Departamento não encontrado")
		}
		s.logger.Error("Failed to get departamento", zap.Error(err))
		return nil, errors.New("Erro ao buscar departamento")
	}

	response := &dto.DepartamentoResponse{
		ID:                     departamento.ID,
		Nome:                   departamento.Nome,
		Gerente:                departamento.Gerente,
		DepartamentoSuperiorID: departamento.DepartamentoSuperiorID,
		Subdepartamentos:       departamento.Subdepartamentos,
		CreatedAt:              departamento.CreatedAt,
		UpdatedAt:              departamento.UpdatedAt,
	}

	s.cache.Set(ctx, cacheKey, response, 5*time.Minute)
	s.logger.Info("Departamento retrieved successfully", zap.String("id", id.String()))

	return response, nil
}

func (s *departamentoService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateDepartamentoRequest) (*model.Departamento, error) {
	s.logger.Info("Updating departamento", zap.String("id", id.String()))

	departamento, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Departamento not found", zap.String("id", id.String()))
			return nil, errors.New("Departamento não encontrado")
		}
		s.logger.Error("Failed to get departamento", zap.Error(err))
		return nil, errors.New("Erro ao buscar departamento")
	}

	if req.Nome != "" {
		departamento.Nome = req.Nome
	}

	if req.GerenteID != nil {
		gerente, err := s.colabRepo.GetByID(ctx, *req.GerenteID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Warn("Gerente not found", zap.String("gerente_id", req.GerenteID.String()))
				return nil, errors.New("Gerente não encontrado")
			}
			s.logger.Error("Failed to get gerente", zap.Error(err))
			return nil, errors.New("Erro ao buscar gerente")
		}

		if gerente.DepartamentoID != id {
			s.logger.Warn("Gerente not in same department", zap.String("gerente_id", req.GerenteID.String()))
			return nil, errors.New("Gerente deve pertencer ao mesmo departamento")
		}

		departamento.GerenteID = *req.GerenteID
	}

	if req.DepartamentoSuperiorID != nil {
		if *req.DepartamentoSuperiorID != uuid.Nil {
			_, err := s.repo.GetByID(ctx, *req.DepartamentoSuperiorID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					s.logger.Warn("Superior department not found", zap.String("departamento_superior_id", req.DepartamentoSuperiorID.String()))
					return nil, errors.New("Departamento superior não encontrado")
				}
				s.logger.Error("Failed to get superior department", zap.Error(err))
				return nil, errors.New("Erro ao buscar departamento superior")
			}

			hasCycle, err := s.repo.HasCycle(ctx, id, *req.DepartamentoSuperiorID)
			if err != nil {
				s.logger.Error("Failed to check cycle", zap.Error(err))
				return nil, errors.New("Erro ao verificar ciclo na hierarquia")
			}
			if hasCycle {
				s.logger.Warn("Cycle detected in hierarchy", zap.String("departamento_superior_id", req.DepartamentoSuperiorID.String()))
				return nil, errors.New("Operação criaria um ciclo na hierarquia de departamentos")
			}
		}
		departamento.DepartamentoSuperiorID = req.DepartamentoSuperiorID
	}

	if err := s.repo.Update(ctx, departamento); err != nil {
		s.logger.Error("Failed to update departamento", zap.Error(err))
		return nil, errors.New("Erro ao atualizar departamento")
	}

	cacheKey := fmt.Sprintf("departamento:%s", id.String())
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info("Departamento updated successfully", zap.String("id", id.String()))
	return departamento, nil
}

func (s *departamentoService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting departamento", zap.String("id", id.String()))

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Departamento not found", zap.String("id", id.String()))
			return errors.New("Departamento não encontrado")
		}
		s.logger.Error("Failed to get departamento", zap.Error(err))
		return errors.New("Erro ao buscar departamento")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete departamento", zap.Error(err))
		return errors.New("Erro ao deletar departamento")
	}

	cacheKey := fmt.Sprintf("departamento:%s", id.String())
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info("Departamento deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *departamentoService) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) (*dto.ListDepartamentosResponse, error) {
	s.logger.Info("Listing departamentos", zap.Int("page", page), zap.Int("page_size", pageSize))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	departamentos, total, err := s.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list departamentos", zap.Error(err))
		return nil, errors.New("Erro ao listar departamentos")
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	s.logger.Info("Departamentos listed successfully", zap.Int64("total", total))

	return &dto.ListDepartamentosResponse{
		Data:       departamentos,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *departamentoService) GetColaboradoresByGerente(ctx context.Context, gerenteID uuid.UUID) ([]model.Colaborador, error) {
	s.logger.Info("Getting colaboradores by gerente", zap.String("gerente_id", gerenteID.String()))

	gerente, err := s.colabRepo.GetByID(ctx, gerenteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Gerente not found", zap.String("gerente_id", gerenteID.String()))
			return nil, errors.New("Gerente não encontrado")
		}
		s.logger.Error("Failed to get gerente", zap.Error(err))
		return nil, errors.New("Erro ao buscar gerente")
	}

	deptIDs, err := s.repo.GetSubdepartamentosRecursive(ctx, gerente.DepartamentoID)
	if err != nil {
		s.logger.Error("Failed to get subdepartamentos", zap.Error(err))
		return nil, errors.New("Erro ao buscar subdepartamentos")
	}

	deptIDs = append(deptIDs, gerente.DepartamentoID)

	colaboradores, err := s.colabRepo.GetByDepartamentoIDs(ctx, deptIDs)
	if err != nil {
		s.logger.Error("Failed to get colaboradores", zap.Error(err))
		return nil, errors.New("Erro ao buscar colaboradores")
	}

	s.logger.Info("Colaboradores retrieved successfully", zap.Int("count", len(colaboradores)))
	return colaboradores, nil
}
