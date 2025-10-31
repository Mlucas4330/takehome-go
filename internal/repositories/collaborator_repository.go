package repositories

import (
	"context"

	"github.com/mlucas4330/takehome-go/internal/models"

	"gorm.io/gorm"
)

type CollaboratorRepository interface {
	Create(context.Context, *models.Collaborator) error
}

type PostgresCollaboratorRepository struct {
	DB *gorm.DB
}

func NewCollaboratorRepository(db *gorm.DB) *PostgresCollaboratorRepository {
	return &PostgresCollaboratorRepository{DB: db}
}

func (r *PostgresCollaboratorRepository) Create(ctx context.Context, col *models.Collaborator) error {
	return nil
}
