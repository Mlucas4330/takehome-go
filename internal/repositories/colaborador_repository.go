package repositories

import (
	"github.com/google/uuid"
	"github.com/mlucas4330/takehome-go/internal/domain"
	"gorm.io/gorm"
)

type ColaboradorRepository struct {
	db *gorm.DB
}

func NewColaboradorRepository(db *gorm.DB) *ColaboradorRepository {
	return &ColaboradorRepository{db: db}
}

func (r *ColaboradorRepository) Create(colaborador *domain.Colaborador) error {
	return r.db.Create(colaborador).Error
}

func (r *ColaboradorRepository) FindByID(id uuid.UUID) (*domain.Colaborador, error) {
	var colaborador domain.Colaborador
	err := r.db.Preload("Departamento.Gerente").First(&colaborador, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &colaborador, nil
}

func (r *ColaboradorRepository) Update(colaborador *domain.Colaborador) error {
	return r.db.Save(colaborador).Error
}

func (r *ColaboradorRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Colaborador{}, "id = ?", id).Error
}

func (r *ColaboradorRepository) List(req domain.ColaboradorListRequest) ([]domain.Colaborador, int64, error) {
	var colaboradores []domain.Colaborador
	var total int64

	query := r.db.Model(&domain.Colaborador{}).Preload("Departamento.Gerente")

	if req.Nome != nil {
		query = query.Where("nome ILIKE ?", "%"+*req.Nome+"%")
	}
	if req.CPF != nil {
		query = query.Where("cpf = ?", *req.CPF)
	}
	if req.RG != nil {
		query = query.Where("rg = ?", *req.RG)
	}
	if req.DepartamentoID != nil {
		query = query.Where("departamento_id = ?", *req.DepartamentoID)
	}

	query.Count(&total)

	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	offset := (req.Page - 1) * req.PageSize
	err := query.Offset(offset).Limit(req.PageSize).Find(&colaboradores).Error

	return colaboradores, total, err
}

func (r *ColaboradorRepository) ExistsByCPF(cpf string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Colaborador{}).Where("cpf = ?", cpf)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *ColaboradorRepository) ExistsByRG(rg string, excludeID *uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Model(&domain.Colaborador{}).Where("rg = ?", rg)
	if excludeID != nil {
		query = query.Where("id != ?", *excludeID)
	}
	err := query.Count(&count).Error
	return count > 0, err
}
