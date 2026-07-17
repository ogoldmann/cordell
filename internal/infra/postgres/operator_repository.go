package postgres

import (
	"context"
	"errors"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5"
)

// OperatorRepository persists operators in PostgreSQL.
type OperatorRepository struct {
	queries *db.Queries
}

// NewOperatorRepository creates an OperatorRepository.
func NewOperatorRepository(queries *db.Queries) *OperatorRepository {
	return &OperatorRepository{
		queries: queries,
	}
}

// Save persists an operator record.
func (r *OperatorRepository) Save(ctx context.Context, operator domain.Operator) error {
	err := r.queries.CreateOperator(ctx, db.CreateOperatorParams{
		ID:           string(operator.ID()),
		Username:     operator.Username(),
		Role:         operator.Role().String(),
		PasswordHash: operator.PasswordHash(),
	})
	if err != nil {
		if isUniqueViolation(err, "operators_username_unique") {
			return domain.ErrDuplicateOperatorUsername
		}

		return err
	}

	return nil
}

// FindByID retrieves an operator by id.
func (r *OperatorRepository) FindByID(ctx context.Context, id domain.OperatorID) (domain.Operator, error) {
	row, err := r.queries.GetOperatorByID(ctx, string(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Operator{}, ports.ErrNotFound
		}

		return domain.Operator{}, err
	}

	return domain.ReconstituteOperator(
		domain.OperatorID(row.ID),
		row.Username,
		domain.OperatorRole(row.Role),
		row.PasswordHash,
		row.Active,
	)
}

// FindByUsername retrieves an operator by username.
func (r *OperatorRepository) FindByUsername(ctx context.Context, username string) (domain.Operator, error) {
	row, err := r.queries.GetOperatorByUsername(ctx, domain.NormalizeOperatorUsername(username))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Operator{}, ports.ErrNotFound
		}

		return domain.Operator{}, err
	}

	return domain.ReconstituteOperator(
		domain.OperatorID(row.ID),
		row.Username,
		domain.OperatorRole(row.Role),
		row.PasswordHash,
		row.Active,
	)
}

var _ ports.OperatorRepository = (*OperatorRepository)(nil)
