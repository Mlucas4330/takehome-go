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
	"takehome-go/internal/validator"
)

type ColaboradorService interface {
	Create(ctx context.Context, req *dto.CreateColaboradorRequest) (*model.Colaborador, error)
	GetByID(ctx context.Context, id uuid.UUID) (*dto.ColaboradorResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *dto.UpdateColaboradorRequest) (*model.Colaborador, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) (*dto.ListColaboradoresResponse, error)
}

type colaboradorService struct {
	repo     repository.ColaboradorRepository
	deptRepo repository.DepartamentoRepository
	cache    database.Cache
	logger   *zap.Logger
}

func NewColaboradorService(
	repo repository.ColaboradorRepository,
	deptRepo repository.DepartamentoRepository,
	cache database.Cache,
	logger *zap.Logger,
) ColaboradorService {
	return &colaboradorService{
		repo:     repo,
		deptRepo: deptRepo,
		cache:    cache,
		logger:   logger,
	}
}

func (s *colaboradorService) Create(ctx context.Context, req *dto.CreateColaboradorRequest) (*model.Colaborador, error) {
	s.logger.Info("Creating colaborador", zap.String("nome", req.Nome), zap.String("cpf", req.CPF))

	if !validator.ValidateCPF(req.CPF) {
		s.logger.Warn("Invalid CPF provided", zap.String("cpf", req.CPF))
		return nil, errors.New("CPF inválido")
	}

	exists, err := s.repo.ExistsByCPF(ctx, req.CPF, nil)
	if err != nil {
		s.logger.Error("Failed to check CPF existence", zap.Error(err))
		return nil, errors.New("Erro ao verificar CPF")
	}
	if exists {
		s.logger.Warn("CPF already exists", zap.String("cpf", req.CPF))
		return nil, errors.New("CPF já cadastrado")
	}

	if req.RG != nil && *req.RG != "" {
		if !validator.ValidateRG(*req.RG) {
			s.logger.Warn("Invalid RG provided", zap.String("rg", *req.RG))
			return nil, errors.New("RG inválido")
		}

		exists, err := s.repo.ExistsByRG(ctx, *req.RG, nil)
		if err != nil {
			s.logger.Error("Failed to check RG existence", zap.Error(err))
			return nil, errors.New("Erro ao verificar RG")
		}
		if exists {
			s.logger.Warn("RG already exists", zap.String("rg", *req.RG))
			return nil, errors.New("RG já cadastrado")
		}
	}

	_, err = s.deptRepo.GetByID(ctx, req.DepartamentoID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Department not found", zap.String("departamento_id", req.DepartamentoID.String()))
			return nil, errors.New("Departamento não encontrado")
		}
		s.logger.Error("Failed to get department", zap.Error(err))
		return nil, errors.New("Erro ao buscar departamento")
	}

	colaborador := &model.Colaborador{
		Nome:           req.Nome,
		CPF:            req.CPF,
		RG:             req.RG,
		DepartamentoID: req.DepartamentoID,
	}

	if err := s.repo.Create(ctx, colaborador); err != nil {
		s.logger.Error("Failed to create colaborador", zap.Error(err))
		return nil, errors.New("Erro ao criar colaborador")
	}

	s.logger.Info("Colaborador created successfully", zap.String("id", colaborador.ID.String()))
	return colaborador, nil
}

func (s *colaboradorService) GetByID(ctx context.Context, id uuid.UUID) (*dto.ColaboradorResponse, error) {
	s.logger.Info("Getting colaborador by ID", zap.String("id", id.String()))

	cacheKey := fmt.Sprintf("colaborador:%s", id.String())
	var cached dto.ColaboradorResponse
	if err := s.cache.Get(ctx, cacheKey, &cached); err == nil {
		s.logger.Debug("Colaborador found in cache", zap.String("id", id.String()))
		return &cached, nil
	}

	colaborador, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Colaborador not found", zap.String("id", id.String()))
			return nil, errors.New("Colaborador não encontrado")
		}
		s.logger.Error("Failed to get colaborador", zap.Error(err))
		return nil, errors.New("Erro ao buscar colaborador")
	}

	response := &dto.ColaboradorResponse{
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

	s.cache.Set(ctx, cacheKey, response, 5*time.Minute)
	s.logger.Info("Colaborador retrieved successfully", zap.String("id", id.String()))

	return response, nil
}

func (s *colaboradorService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateColaboradorRequest) (*model.Colaborador, error) {
	s.logger.Info("Updating colaborador", zap.String("id", id.String()))

	colaborador, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Colaborador not found", zap.String("id", id.String()))
			return nil, errors.New("Colaborador não encontrado")
		}
		s.logger.Error("Failed to get colaborador", zap.Error(err))
		return nil, errors.New("Erro ao buscar colaborador")
	}

	if req.Nome != "" {
		colaborador.Nome = req.Nome
	}

	if req.CPF != "" {
		if !validator.ValidateCPF(req.CPF) {
			s.logger.Warn("Invalid CPF provided", zap.String("cpf", req.CPF))
			return nil, errors.New("CPF inválido")
		}
		exists, err := s.repo.ExistsByCPF(ctx, req.CPF, &id)
		if err != nil {
			s.logger.Error("Failed to check CPF existence", zap.Error(err))
			return nil, errors.New("Erro ao verificar CPF")
		}
		if exists {
			s.logger.Warn("CPF already exists", zap.String("cpf", req.CPF))
			return nil, errors.New("CPF já cadastrado")
		}
		colaborador.CPF = req.CPF
	}

	if req.RG != nil && *req.RG != "" {
		if !validator.ValidateRG(*req.RG) {
			s.logger.Warn("Invalid RG provided", zap.String("rg", *req.RG))
			return nil, errors.New("RG inválido")
		}

		exists, err := s.repo.ExistsByRG(ctx, *req.RG, &id)
		if err != nil {
			s.logger.Error("Failed to check RG existence", zap.Error(err))
			return nil, errors.New("Erro ao verificar RG")
		}
		if exists {
			s.logger.Warn("RG already exists", zap.String("rg", *req.RG))
			return nil, errors.New("RG já cadastrado")
		}
		colaborador.RG = req.RG
	}

	if req.DepartamentoID != nil {
		_, err := s.deptRepo.GetByID(ctx, *req.DepartamentoID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.logger.Warn("Department not found", zap.String("departamento_id", req.DepartamentoID.String()))
				return nil, errors.New("Departamento não encontrado")
			}
			s.logger.Error("Failed to get department", zap.Error(err))
			return nil, errors.New("Erro ao buscar departamento")
		}
		colaborador.DepartamentoID = *req.DepartamentoID
	}

	if err := s.repo.Update(ctx, colaborador); err != nil {
		s.logger.Error("Failed to update colaborador", zap.Error(err))
		return nil, errors.New("Erro ao atualizar colaborador")
	}

	cacheKey := fmt.Sprintf("colaborador:%s", id.String())
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info("Colaborador updated successfully", zap.String("id", id.String()))
	return colaborador, nil
}

func (s *colaboradorService) Delete(ctx context.Context, id uuid.UUID) error {
	s.logger.Info("Deleting colaborador", zap.String("id", id.String()))

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Colaborador not found", zap.String("id", id.String()))
			return errors.New("Colaborador não encontrado")
		}
		s.logger.Error("Failed to get colaborador", zap.Error(err))
		return errors.New("Erro ao buscar colaborador")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete colaborador", zap.Error(err))
		return errors.New("Erro ao deletar colaborador")
	}

	cacheKey := fmt.Sprintf("colaborador:%s", id.String())
	s.cache.Delete(ctx, cacheKey)

	s.logger.Info("Colaborador deleted successfully", zap.String("id", id.String()))
	return nil
}

func (s *colaboradorService) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) (*dto.ListColaboradoresResponse, error) {
	s.logger.Info("Listing colaboradores", zap.Int("page", page), zap.Int("page_size", pageSize))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	colaboradores, total, err := s.repo.List(ctx, filters, page, pageSize)
	if err != nil {
		s.logger.Error("Failed to list colaboradores", zap.Error(err))
		return nil, errors.New("Erro ao listar colaboradores")
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	s.logger.Info("Colaboradores listed successfully", zap.Int64("total", total))

	return &dto.ListColaboradoresResponse{
		Data:       colaboradores,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
