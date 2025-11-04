package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"takehome-go/internal/model"
)

type ColaboradorRepository interface {
	Create(ctx context.Context, colaborador *model.Colaborador) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Colaborador, error)
	Update(ctx context.Context, colaborador *model.Colaborador) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]any, page, pageSize int) ([]model.Colaborador, int64, error)
	ExistsByCPF(ctx context.Context, cpf string, excludeID *uuid.UUID) (bool, error)
	ExistsByRG(ctx context.Context, rg string, excludeID *uuid.UUID) (bool, error)
	GetByDepartamentoIDs(ctx context.Context, ids []uuid.UUID) ([]model.Colaborador, error)
}

type colaboradorRepository struct {
	db *gorm.DB
}

func NewColaboradorRepository(db *gorm.DB) ColaboradorRepository {
	return &colaboradorRepository{db: db}
}

func (r *colaboradorRepository) Create(ctx context.Context, colaborador *model.Colaborador) error {
	return r.db.WithContext(ctx).Create(colaborador).Error
}

func (r *colaboradorRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Colaborador, error) {
	var colaborador model.Colaborador
	err := r.db.WithContext(ctx).
		Preload("Departamento.Gerente").
		First(&colaborador, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &colaborador, nil
}

func (r *colaboradorRepository) Update(ctx context.Context, colaborador *model.Colaborador) error {
	return r.db.WithContext(ctx).Save(colaborador).Error
}

func (r *colaboradorRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Colaborador{}, "id = ?", id).Error
}

func (r *colaboradorRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]model.Colaborador, int64, error) {
	var colaboradores []model.Colaborador
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Colaborador{})

	if nome, ok := filters["nome"].(string); ok && nome != "" {
		query = query.Where("nome ILIKE ?", "%"+nome+"%")
	}
	if cpf, ok := filters["cpf"].(string); ok && cpf != "" {
		query = query.Where("cpf = ?", cpf)
	}
	if rg, ok := filters["rg"].(string); ok && rg != "" {
		query = query.Where("rg = ?", rg)
	}
	if deptID, ok := filters["departamento_id"].(string); ok && deptID != "" {
		query = query.Where("departamento_id = ?", deptID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Preload("Departamento").
		Find(&colaboradores).Error

	return colaboradores, total, err
}

func (r *colaboradorRepository) ExistsByCPF(ctx context.Context, cpf string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Colaborador{}).Where("cpf = ?", cpf)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *colaboradorRepository) ExistsByRG(ctx context.Context, rg string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&model.Colaborador{}).Where("rg = ?", rg)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *colaboradorRepository) GetByDepartamentoIDs(ctx context.Context, ids []uuid.UUID) ([]model.Colaborador, error) {
	var colaboradores []model.Colaborador
	err := r.db.WithContext(ctx).
		Where("departamento_id IN ?", ids).
		Preload("Departamento").
		Find(&colaboradores).Error
	return colaboradores, err
}
