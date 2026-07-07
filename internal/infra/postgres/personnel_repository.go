package postgres

import (
	"context"
	"errors"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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
	err := r.queries.CreatePersonnel(ctx, db.CreatePersonnelParams{
		ID:               string(personnel.ID()),
		FullName:         personnel.FullName(),
		Alias:            personnel.Alias(),
		Rank:             string(personnel.Rank()),
		RegistrationID:   personnel.RegistrationID().String(),
		Section:          string(personnel.Section()),
		OrganizationUnit: string(personnel.OrganizationUnit()),
	})
	if err != nil {
		if isUniqueViolation(err, "personnel_registration_id_unique") {
			return domain.ErrDuplicateRegistrationID
		}

		return err
	}

	return nil
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

	registrationID, err := domain.NewRegistrationID(row.RegistrationID)
	if err != nil {
		return domain.Personnel{}, err
	}

	return domain.ReconstitutePersonnel(
		domain.PersonnelID(row.ID),
		row.FullName,
		row.Alias,
		domain.PersonnelRank(row.Rank),
		registrationID,
		domain.PersonnelSection(row.Section),
		domain.OrganizationUnit(row.OrganizationUnit),
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
		registrationID, err := domain.NewRegistrationID(row.RegistrationID)
		if err != nil {
			return nil, err
		}

		item, err := domain.ReconstitutePersonnel(
			domain.PersonnelID(row.ID),
			row.FullName,
			row.Alias,
			domain.PersonnelRank(row.Rank),
			registrationID,
			domain.PersonnelSection(row.Section),
			domain.OrganizationUnit(row.OrganizationUnit),
			row.Active,
		)
		if err != nil {
			return nil, err
		}

		personnel = append(personnel, item)
	}

	return personnel, nil
}

func isUniqueViolation(err error, constraintName string) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return false
	}

	return pgErr.Code == "23505" && pgErr.ConstraintName == constraintName
}
