package ports

import (
	"context"

	"cordell/internal/domain"
)

// IDGenerator creates unique identifiers for new domain objects.
type IDGenerator interface {
	NewID() (string, error)
}

// PersonnelRepository persists and retrieves personnel records.
type PersonnelRepository interface {
	Save(ctx context.Context, personnel domain.Personnel) error
	FindByID(ctx context.Context, id domain.PersonnelID) (domain.Personnel, error)
	List(ctx context.Context, limit int) ([]domain.Personnel, error)
}

// AssetRepository persists and retrieves asset records.
type AssetRepository interface {
	Save(ctx context.Context, asset domain.Asset) error
	FindByID(ctx context.Context, id domain.AssetID) (domain.Asset, error)
	List(ctx context.Context, limit int) ([]domain.Asset, error)
}

// CustodyRepository persists custody transactions and reads current custody state.
type CustodyRepository interface {
	SaveTransaction(ctx context.Context, transaction domain.CustodyTransaction) error
	CurrentQuantity(ctx context.Context, personnelID domain.PersonnelID, assetID domain.AssetID) (int, error)
}
