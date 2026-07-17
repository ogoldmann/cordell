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

// List retrieves operator summaries for administration.
func (r *OperatorRepository) List(ctx context.Context, limit int) ([]ports.OperatorSummary, error) {
	rows, err := r.queries.ListOperators(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	operators := make([]ports.OperatorSummary, 0, len(rows))

	for _, row := range rows {
		role, err := domain.NewOperatorRole(row.Role)
		if err != nil {
			return nil, err
		}

		operators = append(operators, ports.OperatorSummary{
			ID:        domain.OperatorID(row.ID),
			Username:  row.Username,
			Role:      role,
			Active:    row.Active,
			CreatedAt: row.CreatedAt.Time,
		})
	}

	return operators, nil
}

// Deactivate marks an operator as inactive.
func (r *OperatorRepository) Deactivate(ctx context.Context, id domain.OperatorID) (bool, error) {
	updatedCount, err := r.queries.DeactivateOperator(ctx, string(id))
	if err != nil {
		return false, err
	}

	return updatedCount > 0, nil
}

// CountActiveAdmins counts active admin operators.
func (r *OperatorRepository) CountActiveAdmins(ctx context.Context) (int, error) {
	count, err := r.queries.CountActiveAdminOperators(ctx)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
