package repositories

import (
	"errors"

	"github.com/mlucas4330/takehome-go/internal/models"
	"gorm.io/gorm"
)

type DepartamentRepository interface {
	FindByID(string) (*models.Departament, error)
}

type PostgresDepartamentRepository struct {
	db *gorm.DB
}

func NewDepartamentRepository(db *gorm.DB) *PostgresDepartamentRepository {
	return &PostgresDepartamentRepository{db: db}
}

func (r *PostgresDepartamentRepository) FindByID(id string) (*models.Departament, error) {
	var dept models.Departament
	if err := r.db.First(&dept, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &dept, nil
}
