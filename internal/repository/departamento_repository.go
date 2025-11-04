package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"takehome-go/internal/model"
)

type DepartamentoRepository interface {
	Create(ctx context.Context, departamento *model.Departamento) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Departamento, error)
	GetByIDWithHierarchy(ctx context.Context, id uuid.UUID) (*model.Departamento, error)
	Update(ctx context.Context, departamento *model.Departamento) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]model.Departamento, int64, error)
	HasCycle(ctx context.Context, id, superiorID uuid.UUID) (bool, error)
	GetSubdepartamentosRecursive(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
}

type departamentoRepository struct {
	db *gorm.DB
}

func NewDepartamentoRepository(db *gorm.DB) DepartamentoRepository {
	return &departamentoRepository{db: db}
}

func (r *departamentoRepository) Create(ctx context.Context, departamento *model.Departamento) error {
	return r.db.WithContext(ctx).Create(departamento).Error
}

func (r *departamentoRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Departamento, error) {
	var departamento model.Departamento
	err := r.db.WithContext(ctx).
		Preload("Gerente").
		First(&departamento, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &departamento, nil
}

func (r *departamentoRepository) GetByIDWithHierarchy(ctx context.Context, id uuid.UUID) (*model.Departamento, error) {
	type DeptResult struct {
		ID                     uuid.UUID
		Nome                   string
		GerenteID              uuid.UUID
		DepartamentoSuperiorID *uuid.UUID
	}

	query := `
		WITH RECURSIVE dept_tree AS (
			SELECT id, nome, gerente_id, departamento_superior_id
			FROM departamentos
			WHERE id = $1
			
			UNION ALL
			
			SELECT d.id, d.nome, d.gerente_id, d.departamento_superior_id
			FROM departamentos d
			INNER JOIN dept_tree dt ON d.departamento_superior_id = dt.id
		)
		SELECT * FROM dept_tree
	`

	var results []DeptResult
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&results).Error; err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	deptMap := make(map[uuid.UUID]*model.Departamento)
	var gerenteIDs []uuid.UUID

	for _, res := range results {
		dept := &model.Departamento{
			ID:                     res.ID,
			Nome:                   res.Nome,
			GerenteID:              res.GerenteID,
			DepartamentoSuperiorID: res.DepartamentoSuperiorID,
			Subdepartamentos:       []model.Departamento{},
		}
		deptMap[res.ID] = dept
		gerenteIDs = append(gerenteIDs, res.GerenteID)
	}

	var gerentes []model.Colaborador
	if err := r.db.WithContext(ctx).Where("id IN ?", gerenteIDs).Find(&gerentes).Error; err != nil {
		return nil, err
	}

	gerenteMap := make(map[uuid.UUID]*model.Colaborador)
	for i := range gerentes {
		gerenteMap[gerentes[i].ID] = &gerentes[i]
	}

	for _, dept := range deptMap {
		if gerente, ok := gerenteMap[dept.GerenteID]; ok {
			dept.Gerente = gerente
		}
	}

	for _, dept := range deptMap {
		if dept.DepartamentoSuperiorID != nil {
			if parent, ok := deptMap[*dept.DepartamentoSuperiorID]; ok {
				parent.Subdepartamentos = append(parent.Subdepartamentos, *dept)
			}
		}
	}

	return deptMap[id], nil
}

func (r *departamentoRepository) Update(ctx context.Context, departamento *model.Departamento) error {
	return r.db.WithContext(ctx).Save(departamento).Error
}

func (r *departamentoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Departamento{}, "id = ?", id).Error
}

func (r *departamentoRepository) List(ctx context.Context, filters map[string]interface{}, page, pageSize int) ([]model.Departamento, int64, error) {
	var departamentos []model.Departamento
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Departamento{})

	if nome, ok := filters["nome"].(string); ok && nome != "" {
		query = query.Where("nome ILIKE ?", "%"+nome+"%")
	}
	if gerenteNome, ok := filters["gerente_nome"].(string); ok && gerenteNome != "" {
		query = query.Joins("JOIN colaboradores ON colaboradores.id = departamentos.gerente_id").
			Where("colaboradores.nome ILIKE ?", "%"+gerenteNome+"%")
	}
	if superiorID, ok := filters["departamento_superior_id"].(string); ok && superiorID != "" {
		query = query.Where("departamento_superior_id = ?", superiorID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).
		Preload("Gerente").
		Preload("DepartamentoSuperior").
		Find(&departamentos).Error

	return departamentos, total, err
}

func (r *departamentoRepository) HasCycle(ctx context.Context, id, superiorID uuid.UUID) (bool, error) {
	query := `
		WITH RECURSIVE dept_hierarchy AS (
			SELECT id, departamento_superior_id
			FROM departamentos
			WHERE id = $1
			
			UNION ALL
			
			SELECT d.id, d.departamento_superior_id
			FROM departamentos d
			INNER JOIN dept_hierarchy dh ON d.id = dh.departamento_superior_id
		)
		SELECT EXISTS(SELECT 1 FROM dept_hierarchy WHERE id = $2)
	`

	var hasCycle bool
	err := r.db.WithContext(ctx).Raw(query, superiorID, id).Scan(&hasCycle).Error
	return hasCycle, err
}

func (r *departamentoRepository) GetSubdepartamentosRecursive(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
	query := `
		WITH RECURSIVE subdepts AS (
			SELECT id
			FROM departamentos
			WHERE departamento_superior_id = $1
			
			UNION ALL
			
			SELECT d.id
			FROM departamentos d
			INNER JOIN subdepts s ON d.departamento_superior_id = s.id
		)
		SELECT id FROM subdepts
	`

	var ids []uuid.UUID
	err := r.db.WithContext(ctx).Raw(query, id).Scan(&ids).Error
	return ids, err
}