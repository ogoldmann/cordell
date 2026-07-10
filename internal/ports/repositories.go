package ports

import (
	"context"
	"time"

	"cordell/internal/domain"
)

// CurrentCustodyItem represents a current custody balance joined with asset display data.
type CurrentCustodyItem struct {
	PersonnelID domain.PersonnelID
	AssetID     domain.AssetID
	AssetName   string
	Quantity    int
}

// CurrentAssetHolder represents a current custody balance joined with personnel display data.
type CurrentAssetHolder struct {
	AssetID           domain.AssetID
	PersonnelID       domain.PersonnelID
	PersonnelFullName string
	Quantity          int
}

// CustodyHistoryLine represents one asset line inside a custody history entry.
type CustodyHistoryLine struct {
	AssetID   domain.AssetID
	AssetName string
	Quantity  int
}

// CustodyHistoryEntry represents a custody transaction with its lines.
type CustodyHistoryEntry struct {
	ID          domain.CustodyTransactionID
	Type        domain.CustodyTransactionType
	PersonnelID domain.PersonnelID
	Notes       string
	CreatedAt   time.Time
	Lines       []CustodyHistoryLine
}

// IDGenerator creates unique identifiers for new domain objects.
type IDGenerator interface {
	NewID() (string, error)
}

// PersonnelRepository persists and retrieves personnel records.
type PersonnelRepository interface {
	Save(ctx context.Context, personnel domain.Personnel) error
	FindByID(ctx context.Context, id domain.PersonnelID) (domain.Personnel, error)
	List(ctx context.Context, limit int) ([]domain.Personnel, error)
	Search(ctx context.Context, query string, limit int) ([]domain.Personnel, error)
}

// AssetRepository persists and retrieves asset records.
type AssetRepository interface {
	Save(ctx context.Context, asset domain.Asset) error
	FindByID(ctx context.Context, id domain.AssetID) (domain.Asset, error)
	List(ctx context.Context, limit int) ([]domain.Asset, error)
	Search(ctx context.Context, query string, limit int) ([]domain.Asset, error)
}

// CustodyRepository persists custody transactions and reads current custody state.
type CustodyRepository interface {
	SaveTransaction(ctx context.Context, transaction domain.CustodyTransaction) error
	CurrentQuantity(ctx context.Context, personnelID domain.PersonnelID, assetID domain.AssetID) (int, error)
	ListCurrentByPersonnel(ctx context.Context, personnelID domain.PersonnelID) ([]CurrentCustodyItem, error)
	ListCurrentByAsset(ctx context.Context, assetID domain.AssetID) ([]CurrentAssetHolder, error)
	ListHistoryByPersonnel(ctx context.Context, personnelID domain.PersonnelID, limit int) ([]CustodyHistoryEntry, error)
}

// OperatorRepository persists operator records.
type OperatorRepository interface {
	Save(ctx context.Context, operator domain.Operator) error
	FindByID(ctx context.Context, id domain.OperatorID) (domain.Operator, error)
	FindByUsername(ctx context.Context, username string) (domain.Operator, error)
}
