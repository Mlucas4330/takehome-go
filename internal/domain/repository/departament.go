package repository

import (
	"context"
	"takehome-go/internal/domain/model"

	"gorm.io/gorm"
)

type DepFilter struct {
	Nome                   *string `json:"nome"`
	GerenteNome            *string `json:"gerente_nome"`
	DepartamentoSuperiorID *string `json:"departamento_superior_id"`
}

type DepartamentoRepository interface {
	Create(ctx context.Context, d *model.Departament) error
	GetByID(ctx context.Context, id string) (*model.Departament, error)
	Update(ctx context.Context, d *model.Departament) error
	Delete(ctx context.Context, id string) error

	List(ctx context.Context, f DepFilter, p Page) (PageResult[model.Departament], error)

	GetSubtreeIDs(ctx context.Context, rootID string) ([]string, error)
	WouldCreateCycle(ctx context.Context, deptID, newParentID string) (bool, error)
	FindByGerenteID(ctx context.Context, gerenteID string) (*model.Departament, error)
}

type departamentoRepository struct {
	db *gorm.DB
}

func NewDepartamentoRepository(db *gorm.DB) DepartamentoRepository {
	return &departamentoRepository{db: db}
}

func (r *departamentoRepository) Create(ctx context.Context, d *model.Departament) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *departamentoRepository) GetByID(ctx context.Context, id string) (*model.Departament, error) {
	var out model.Departament
	if err := r.db.WithContext(ctx).First(&out, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *departamentoRepository) Update(ctx context.Context, d *model.Departament) error {
	return r.db.WithContext(ctx).Save(d).Error
}

func (r *departamentoRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&model.Departament{}, "id = ?", id).Error
}

func (r *departamentoRepository) List(ctx context.Context, f DepFilter, p Page) (PageResult[model.Departament], error) {
	q := r.db.WithContext(ctx).Model(&model.Departament{})

	if f.Nome != nil && *f.Nome != "" {
		q = q.Where("nome ILIKE ?", "%"+*f.Nome+"%")
	}

	// gerente_nome exige join com colaboradores
	if f.GerenteNome != nil && *f.GerenteNome != "" {
		q = q.Joins("JOIN colaboradores ON colaboradores.id = departamentos.gerente_id").
			Where("colaboradores.nome ILIKE ?", "%"+*f.GerenteNome+"%")
	}

	if f.DepartamentoSuperiorID != nil && *f.DepartamentoSuperiorID != "" {
		q = q.Where("departamento_superior_id = ?", *f.DepartamentoSuperiorID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return PageResult[model.Departament]{}, err
	}

	var items []model.Departament
	if err := paginate(q.Order("nome ASC"), p).Find(&items).Error; err != nil {
		return PageResult[model.Departament]{}, err
	}

	return PageResult[model.Departament]{
		Items:    items,
		Total:    total,
		Page:     p.Page,
		PageSize: p.PageSize,
	}, nil
}

func (r *departamentoRepository) GetSubtreeIDs(ctx context.Context, rootID string) ([]string, error) {
	type row struct{ ID string }
	var rows []row
	sql := `
WITH RECURSIVE subdeps AS (
  SELECT id, departamento_superior_id
  FROM departamentos
  WHERE id = ?
  UNION ALL
  SELECT d.id, d.departamento_superior_id
  FROM departamentos d
  JOIN subdeps s ON d.departamento_superior_id = s.id
)
SELECT id FROM subdeps;
`
	if err := r.db.WithContext(ctx).Raw(sql, rootID).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, v := range rows {
		out = append(out, v.ID)
	}
	return out, nil
}

func (r *departamentoRepository) WouldCreateCycle(ctx context.Context, deptID, newParentID string) (bool, error) {
	// Se o novo pai tem como ancestral o próprio deptID, criaria ciclo.
	type row struct{ ID string }
	var rows []row
	sql := `
WITH RECURSIVE ancestors AS (
  SELECT id, departamento_superior_id
  FROM departamentos
  WHERE id = ?
  UNION ALL
  SELECT d.id, d.departamento_superior_id
  FROM departamentos d
  JOIN ancestors a ON d.id = a.departamento_superior_id
)
SELECT id FROM ancestors;
`
	if err := r.db.WithContext(ctx).Raw(sql, newParentID).Scan(&rows).Error; err != nil {
		return false, err
	}
	for _, v := range rows {
		if v.ID == deptID {
			return true, nil
		}
	}
	return false, nil
}

func (r *departamentoRepository) FindByGerenteID(ctx context.Context, gerenteID string) (*model.Departament, error) {
	var d model.Departament
	if err := r.db.WithContext(ctx).First(&d, "gerente_id = ?", gerenteID).Error; err != nil {
		return nil, err
	}
	return &d, nil
}
