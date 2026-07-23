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
	if err != nil {
		if isUniqueViolation(err, "assets_name_unique_idx") {
			return domain.ErrDuplicateAssetName
		}

		return err
	}

	return nil
}

// Update updates editable asset fields.
func (r *AssetRepository) Update(ctx context.Context, asset domain.Asset) error {
	err := r.queries.UpdateAsset(ctx, db.UpdateAssetParams{
		ID:   string(asset.ID()),
		Name: asset.Name(),
	})
	if err != nil {
		if isUniqueViolation(err, "assets_name_unique_idx") {
			return domain.ErrDuplicateAssetName
		}

		return err
	}

	return nil
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

	return assetFromRow(row)
}

// FindByName finds an asset by name.
func (r *AssetRepository) FindByName(
	ctx context.Context,
	name string,
) (domain.Asset, bool, error) {
	row, err := r.queries.FindAssetByName(ctx, domain.NormalizeAssetName(name))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Asset{}, false, nil
		}

		return domain.Asset{}, false, err
	}

	asset, err := assetFromRow(row)
	if err != nil {
		return domain.Asset{}, false, err
	}

	return asset, true, nil
}

// FindByNameExcludingID finds an asset by name excluding one asset ID.
func (r *AssetRepository) FindByNameExcludingID(
	ctx context.Context,
	name string,
	excludedID domain.AssetID,
) (domain.Asset, bool, error) {
	row, err := r.queries.FindAssetByNameExcludingID(ctx, db.FindAssetByNameExcludingIDParams{
		Name:       domain.NormalizeAssetName(name),
		ExcludedID: string(excludedID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Asset{}, false, nil
		}

		return domain.Asset{}, false, err
	}

	asset, err := assetFromRow(row)
	if err != nil {
		return domain.Asset{}, false, err
	}

	return asset, true, nil
}

func assetFromRow(row db.Asset) (domain.Asset, error) {
	return domain.ReconstituteAsset(
		domain.AssetID(row.ID),
		row.Name,
		row.Active,
	)
}

// Search retrieves asset records matching a search query.
func (r *AssetRepository) Search(
	ctx context.Context,
	query string,
	limit int,
	statusFilter ports.RecordStatusFilter,
) ([]domain.Asset, error) {
	rows, err := r.queries.SearchAssets(ctx, db.SearchAssetsParams{
		SearchPatterns: buildTextSearchPatterns(query),
		StatusFilter:   string(statusFilter),
		LimitCount:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	assets := make([]domain.Asset, 0, len(rows))

	for _, row := range rows {
		item, err := assetFromRow(row)
		if err != nil {
			return nil, err
		}

		assets = append(assets, item)
	}

	return assets, nil
}

var _ ports.AssetRepository = (*AssetRepository)(nil)

// List retrieves recent asset records.
func (r *AssetRepository) List(
	ctx context.Context,
	limit int,
	statusFilter ports.RecordStatusFilter,
) ([]domain.Asset, error) {
	rows, err := r.queries.ListAssets(ctx, db.ListAssetsParams{
		StatusFilter: string(statusFilter),
		LimitCount:   int32(limit),
	})
	if err != nil {
		return nil, err
	}

	assets := make([]domain.Asset, 0, len(rows))

	for _, row := range rows {
		item, err := assetFromRow(row)
		if err != nil {
			return nil, err
		}

		assets = append(assets, item)
	}

	return assets, nil
}

// Deactivate marks an asset record as inactive.
func (r *AssetRepository) Deactivate(ctx context.Context, id domain.AssetID) (bool, error) {
	updatedCount, err := r.queries.DeactivateAsset(ctx, string(id))
	if err != nil {
		return false, err
	}

	return updatedCount > 0, nil
}

// Reactivate marks an asset record as active.
func (r *AssetRepository) Reactivate(ctx context.Context, id domain.AssetID) (bool, error) {
	updatedCount, err := r.queries.ReactivateAsset(ctx, string(id))
	if err != nil {
		return false, err
	}

	return updatedCount > 0, nil
}
