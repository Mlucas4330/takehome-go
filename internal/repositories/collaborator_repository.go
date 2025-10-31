package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mlucas4330/takehome-go/internal/models"
	"gorm.io/gorm"
)

type CollaboratorRepository interface {
	Create(ctx context.Context, col *models.Collaborator) error
}

type PostgresCollaboratorRepository struct {
	DB *gorm.DB
}

func NewCollaboratorRepository(db *gorm.DB) *PostgresCollaboratorRepository {
	return &PostgresCollaboratorRepository{DB: db}
}

func (r *PostgresCollaboratorRepository) Create(ctx context.Context, col *models.Collaborator) error {
	if err := r.DB.WithContext(ctx).Create(col).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "colaboradores_cpf_unique":
				return ErrUniqueViolation
			case "colaboradores_rg_unique":
				return ErrUniqueViolation
			case "departamentos_gerente_fk":
				return ErrForeignKey
			}
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return ErrUniqueViolation
		}
		return err
	}
	return nil
}
