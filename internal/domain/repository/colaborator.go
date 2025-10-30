package repository

import (
	"context"
	"takehome-go/internal/domain/model"

	"gorm.io/gorm"
)

type ColabFilter struct {
	Nome           *string `json:"nome"`
	CPF            *string `json:"cpf"`
	RG             *string `json:"rg"`
	DepartamentoID *string `json:"departamento_id"`
}

type ColaboradorRepository interface {
	Create(ctx context.Context, c *model.Colaborator) error
	GetByID(ctx context.Context, id string) (*model.Colaborator, error)
	Update(ctx context.Context, c *model.Colaborator) error
	Delete(ctx context.Context, id string) error

	ExistsByCPF(ctx context.Context, cpf string, excludeID *string) (bool, error)
	ExistsByRG(ctx context.Context, rg string, excludeID *string) (bool, error)

	List(ctx context.Context, f ColabFilter, p Page) (PageResult[model.Colaborator], error)
}

type colaboradorRepository struct {
	db *gorm.DB
}

func NewColaboradorRepository(db *gorm.DB) ColaboradorRepository {
	return &colaboradorRepository{db: db}
}

func (r *colaboradorRepository) Create(ctx context.Context, c *model.Colaborator) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *colaboradorRepository) GetByID(ctx context.Context, id string) (*model.Colaborator, error) {
	var out model.Colaborator
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *colaboradorRepository) Update(ctx context.Context, c *model.Colaborator) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *colaboradorRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Colaborator{}, "id = ?", id).Error
}

func (r *colaboradorRepository) ExistsByCPF(ctx context.Context, cpf string, excludeID *string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Colaborator{}).Where("cpf = ?", cpf)
	if excludeID != nil && *excludeID != "" {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *colaboradorRepository) ExistsByRG(ctx context.Context, rg string, excludeID *string) (bool, error) {
	q := r.db.WithContext(ctx).Model(&model.Colaborator{}).Where("rg = ?", rg)
	if excludeID != nil && *excludeID != "" {
		q = q.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func paginate(db *gorm.DB, p Page) *gorm.DB {
	p.Normalize()
	offset := (p.Page - 1) * p.PageSize
	return db.Offset(offset).Limit(p.PageSize)
}

func (r *colaboradorRepository) List(ctx context.Context, f ColabFilter, p Page) (PageResult[model.Colaborator], error) {
	q := r.db.WithContext(ctx).Model(&model.Colaborator{})
	if f.Nome != nil && *f.Nome != "" {
		q = q.Where("nome ILIKE ?", "%"+*f.Nome+"%")
	}
	if f.CPF != nil && *f.CPF != "" {
		q = q.Where("cpf = ?", *f.CPF)
	}
	if f.RG != nil && *f.RG != "" {
		q = q.Where("rg = ?", *f.RG)
	}
	if f.DepartamentoID != nil && *f.DepartamentoID != "" {
		q = q.Where("departamento_id = ?", *f.DepartamentoID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return PageResult[model.Colaborator]{}, err
	}

	var items []model.Colaborator
	if err := paginate(q.Order("nome ASC"), p).Find(&items).Error; err != nil {
		return PageResult[model.Colaborator]{}, err
	}

	return PageResult[model.Colaborator]{
		Items:    items,
		Total:    total,
		Page:     p.Page,
		PageSize: p.PageSize,
	}, nil
}