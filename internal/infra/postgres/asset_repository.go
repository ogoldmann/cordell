package postgres

import (
	"context"
	"errors"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5"
)

// AssetRepository persists asset records in PostgreSQL.
type AssetRepository struct {
	queries *db.Queries
}

// NewAssetRepository creates a PostgreSQL asset repository.
func NewAssetRepository(queries *db.Queries) *AssetRepository {
	return &AssetRepository{
		queries: queries,
	}
}

// Save persists an asset record.
func (r *AssetRepository) Save(ctx context.Context, asset domain.Asset) error {
	_, err := r.queries.CreateAsset(ctx, db.CreateAssetParams{
		ID:     string(asset.ID()),
		Name:   asset.Name(),
		Active: asset.Active(),
	})
	return err
}

// FindByID retrieves an asset record by identifier.
func (r *AssetRepository) FindByID(ctx context.Context, id domain.AssetID) (domain.Asset, error) {
	row, err := r.queries.GetAsset(ctx, string(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Asset{}, ports.ErrNotFound
		}

		return domain.Asset{}, err
	}

	return domain.ReconstituteAsset(
		domain.AssetID(row.ID),
		row.Name,
		row.Active,
	)
}

var _ ports.AssetRepository = (*AssetRepository)(nil)
