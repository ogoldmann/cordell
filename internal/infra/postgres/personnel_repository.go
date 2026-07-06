package postgres

import (
	"context"
	"errors"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5"
)

// PersonnelRepository persists personnel records in PostgreSQL.
type PersonnelRepository struct {
	queries *db.Queries
}

// NewPersonnelRepository creates a PostgreSQL personnel repository.
func NewPersonnelRepository(queries *db.Queries) *PersonnelRepository {
	return &PersonnelRepository{
		queries: queries,
	}
}

// Save persists a personnel record.
func (r *PersonnelRepository) Save(ctx context.Context, personnel domain.Personnel) error {
	_, err := r.queries.CreatePersonnel(ctx, db.CreatePersonnelParams{
		ID:       string(personnel.ID()),
		FullName: personnel.FullName(),
		Active:   personnel.Active(),
	})
	return err
}

// FindByID retrieves a personnel record by identifier.
func (r *PersonnelRepository) FindByID(ctx context.Context, id domain.PersonnelID) (domain.Personnel, error) {
	row, err := r.queries.GetPersonnel(ctx, string(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Personnel{}, ports.ErrNotFound
		}

		return domain.Personnel{}, err
	}

	return domain.ReconstitutePersonnel(
		domain.PersonnelID(row.ID),
		row.FullName,
		row.Active,
	)
}

var _ ports.PersonnelRepository = (*PersonnelRepository)(nil)

// List retrieves recent personnel records.
func (r *PersonnelRepository) List(ctx context.Context, limit int) ([]domain.Personnel, error) {
	rows, err := r.queries.ListPersonnel(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	personnel := make([]domain.Personnel, 0, len(rows))

	for _, row := range rows {
		item, err := domain.ReconstitutePersonnel(
			domain.PersonnelID(row.ID),
			row.FullName,
			row.Active,
		)
		if err != nil {
			return nil, err
		}

		personnel = append(personnel, item)
	}

	return personnel, nil
}
