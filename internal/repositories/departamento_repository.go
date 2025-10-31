package repositories

import (
	"github.com/google/uuid"
	"github.com/mlucas4330/takehome-go/internal/domain"
	"gorm.io/gorm"
)

type DepartamentoRepository struct {
	db *gorm.DB
}

func NewDepartamentoRepository(db *gorm.DB) *DepartamentoRepository {
	return &DepartamentoRepository{db: db}
}

func (r *DepartamentoRepository) Create(departamento *domain.Departamento) error {
	return r.db.Create(departamento).Error
}

func (r *DepartamentoRepository) FindByID(id uuid.UUID) (*domain.Departamento, error) {
	var departamento domain.Departamento
	err := r.db.Preload("Gerente").First(&departamento, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &departamento, nil
}

func (r *DepartamentoRepository) FindByIDWithHierarchy(id uuid.UUID) (*domain.Departamento, error) {
	query := `
        WITH RECURSIVE dept_tree AS (
            SELECT id, nome, gerente_id, departamento_superior_id, created_at, updated_at
            FROM departamentos
            WHERE id = @root
            
            UNION ALL
            
            SELECT d.id, d.nome, d.gerente_id, d.departamento_superior_id, d.created_at, d.updated_at
            FROM departamentos d
            INNER JOIN dept_tree dt ON d.departamento_superior_id = dt.id
        )
        SELECT * FROM dept_tree
    `
	var allDepts []domain.Departamento
	err := r.db.Raw(query, map[string]interface{}{"root": id}).Scan(&allDepts).Error
	if err != nil {
		return nil, err
	}
	if len(allDepts) == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	deptMap := make(map[uuid.UUID]*domain.Departamento, len(allDepts))
	for i := range allDepts {
		d := allDepts[i]
		deptMap[d.ID] = &allDepts[i]
	}
	for _, dept := range deptMap {
		var gerente domain.Colaborador
		if err := r.db.First(&gerente, "id = ?", dept.GerenteID).Error; err == nil {
			dept.Gerente = &gerente
		}
	}
	for _, dept := range deptMap {
		if dept.DepartamentoSuperiorID != nil {
			if parent, ok := deptMap[*dept.DepartamentoSuperiorID]; ok {
				parent.Subdepartamentos = append(parent.Subdepartamentos, *dept)
			}
		}
	}
	root, ok := deptMap[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return root, nil
}

func (r *DepartamentoRepository) Update(departamento *domain.Departamento) error {
	return r.db.Save(departamento).Error
}

func (r *DepartamentoRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Departamento{}, "id = ?", id).Error
}

func (r *DepartamentoRepository) List(req domain.DepartamentoListRequest) ([]domain.Departamento, int64, error) {
	var departamentos []domain.Departamento
	var total int64
	query := r.db.Model(&domain.Departamento{}).Preload("Gerente")
	if req.Nome != nil && *req.Nome != "" {
		query = query.Where("nome ILIKE ?", "%"+*req.Nome+"%")
	}
	if req.GerenteNome != nil && *req.GerenteNome != "" {
		query = query.Joins("JOIN colaboradores ON colaboradores.id = departamentos.gerente_id").
			Where("colaboradores.nome ILIKE ?", "%"+*req.GerenteNome+"%")
	}
	if req.DepartamentoSuperiorID != nil {
		query = query.Where("departamento_superior_id = ?", *req.DepartamentoSuperiorID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Find(&departamentos).Error
	return departamentos, total, err
}

func (r *DepartamentoRepository) WouldCreateCycle(deptID uuid.UUID, superiorID uuid.UUID) (bool, error) {
	query := `
        WITH RECURSIVE dept_hierarchy AS (
            SELECT id, departamento_superior_id
            FROM departamentos
            WHERE id = @start
            
            UNION ALL
            
            SELECT d.id, d.departamento_superior_id
            FROM departamentos d
            INNER JOIN dept_hierarchy dh ON d.id = dh.departamento_superior_id
        )
        SELECT COUNT(*) FROM dept_hierarchy WHERE id = @target
    `
	var count int64
	err := r.db.Raw(query, map[string]interface{}{"start": superiorID, "target": deptID}).Scan(&count).Error
	return count > 0, err
}

func (r *DepartamentoRepository) GetSubordinateColaboradores(gerenteID uuid.UUID) ([]domain.Colaborador, error) {
	query := `
        WITH RECURSIVE dept_tree AS (
            SELECT id
            FROM departamentos
            WHERE gerente_id = @gerente
            
            UNION ALL
            
            SELECT d.id
            FROM departamentos d
            INNER JOIN dept_tree dt ON d.departamento_superior_id = dt.id
        )
        SELECT c.* FROM colaboradores c
        WHERE c.departamento_id IN (SELECT id FROM dept_tree)
    `
	var colaboradores []domain.Colaborador
	err := r.db.Raw(query, map[string]interface{}{"gerente": gerenteID}).Scan(&colaboradores).Error
	return colaboradores, err
}
