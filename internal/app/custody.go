package app

import (
	"context"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

// CustodyLineCommand contains the input data for one custody transaction line.
type CustodyLineCommand struct {
	AssetID  domain.AssetID
	Quantity int
}

// RegisterCheckoutCommand contains the input data required to register an asset checkout.
type RegisterCheckoutCommand struct {
	PersonnelID domain.PersonnelID
	Lines       []CustodyLineCommand
	Notes       string
}

// RegisterCheckoutService handles the asset checkout use case.
type RegisterCheckoutService struct {
	personnelRepository ports.PersonnelRepository
	assetRepository     ports.AssetRepository
	custodyRepository   ports.CustodyRepository
	idGenerator         ports.IDGenerator
}

// NewRegisterCheckoutService creates a RegisterCheckoutService with its dependencies.
func NewRegisterCheckoutService(
	personnelRepository ports.PersonnelRepository,
	assetRepository ports.AssetRepository,
	custodyRepository ports.CustodyRepository,
	idGenerator ports.IDGenerator,
) *RegisterCheckoutService {
	return &RegisterCheckoutService{
		personnelRepository: personnelRepository,
		assetRepository:     assetRepository,
		custodyRepository:   custodyRepository,
		idGenerator:         idGenerator,
	}
}

// Execute registers and persists a checkout custody transaction.
func (s *RegisterCheckoutService) Execute(
	ctx context.Context,
	cmd RegisterCheckoutCommand,
) (domain.CustodyTransaction, error) {
	if cmd.PersonnelID == "" {
		return domain.CustodyTransaction{}, domain.ErrEmptyPersonnelID
	}

	if _, err := s.personnelRepository.FindByID(ctx, cmd.PersonnelID); err != nil {
		return domain.CustodyTransaction{}, err
	}

	lines := make([]domain.CustodyLine, 0, len(cmd.Lines))

	for _, lineCommand := range cmd.Lines {
		quantity, err := domain.NewQuantity(lineCommand.Quantity)
		if err != nil {
			return domain.CustodyTransaction{}, err
		}

		line, err := domain.NewCustodyLine(lineCommand.AssetID, quantity)
		if err != nil {
			return domain.CustodyTransaction{}, err
		}

		if _, err := s.assetRepository.FindByID(ctx, lineCommand.AssetID); err != nil {
			return domain.CustodyTransaction{}, err
		}

		lines = append(lines, line)
	}

	id, err := s.idGenerator.NewID()
	if err != nil {
		return domain.CustodyTransaction{}, err
	}

	transactionID := domain.CustodyTransactionID(id)

	transaction, err := domain.NewCustodyTransaction(
		transactionID,
		domain.CustodyTransactionTypeCheckout,
		cmd.PersonnelID,
		lines,
		cmd.Notes,
	)
	if err != nil {
		return domain.CustodyTransaction{}, err
	}

	if err := s.custodyRepository.SaveTransaction(ctx, transaction); err != nil {
		return domain.CustodyTransaction{}, err
	}

	return transaction, nil
}

// RegisterReturnCommand contains the input data required to register an asset return.
type RegisterReturnCommand struct {
	PersonnelID domain.PersonnelID
	Lines       []CustodyLineCommand
	Notes       string
}

// RegisterReturnService handles the asset return use case.
type RegisterReturnService struct {
	personnelRepository ports.PersonnelRepository
	assetRepository     ports.AssetRepository
	custodyRepository   ports.CustodyRepository
	idGenerator         ports.IDGenerator
}

// NewRegisterReturnService creates a RegisterReturnService with its dependencies.
func NewRegisterReturnService(
	personnelRepository ports.PersonnelRepository,
	assetRepository ports.AssetRepository,
	custodyRepository ports.CustodyRepository,
	idGenerator ports.IDGenerator,
) *RegisterReturnService {
	return &RegisterReturnService{
		personnelRepository: personnelRepository,
		assetRepository:     assetRepository,
		custodyRepository:   custodyRepository,
		idGenerator:         idGenerator,
	}
}

// Execute registers and persists a return custody transaction.
func (s *RegisterReturnService) Execute(
	ctx context.Context,
	cmd RegisterReturnCommand,
) (domain.CustodyTransaction, error) {
	if cmd.PersonnelID == "" {
		return domain.CustodyTransaction{}, domain.ErrEmptyPersonnelID
	}

	if _, err := s.personnelRepository.FindByID(ctx, cmd.PersonnelID); err != nil {
		return domain.CustodyTransaction{}, err
	}

	lines := make([]domain.CustodyLine, 0, len(cmd.Lines))

	for _, lineCommand := range cmd.Lines {
		quantity, err := domain.NewQuantity(lineCommand.Quantity)
		if err != nil {
			return domain.CustodyTransaction{}, err
		}

		line, err := domain.NewCustodyLine(lineCommand.AssetID, quantity)
		if err != nil {
			return domain.CustodyTransaction{}, err
		}

		if _, err := s.assetRepository.FindByID(ctx, lineCommand.AssetID); err != nil {
			return domain.CustodyTransaction{}, err
		}

		currentQuantity, err := s.custodyRepository.CurrentQuantity(ctx, cmd.PersonnelID, lineCommand.AssetID)
		if err != nil {
			return domain.CustodyTransaction{}, err
		}

		if currentQuantity < quantity.Int() {
			return domain.CustodyTransaction{}, domain.ErrInsufficientCustodyBalance
		}

		lines = append(lines, line)
	}

	id, err := s.idGenerator.NewID()
	if err != nil {
		return domain.CustodyTransaction{}, err
	}

	transactionID := domain.CustodyTransactionID(id)

	transaction, err := domain.NewCustodyTransaction(
		transactionID,
		domain.CustodyTransactionTypeReturn,
		cmd.PersonnelID,
		lines,
		cmd.Notes,
	)
	if err != nil {
		return domain.CustodyTransaction{}, err
	}

	if err := s.custodyRepository.SaveTransaction(ctx, transaction); err != nil {
		return domain.CustodyTransaction{}, err
	}

	return transaction, nil
}

// CurrentCustodyItem contains current custody display data for application use cases.
type CurrentCustodyItem struct {
	PersonnelID domain.PersonnelID
	AssetID     domain.AssetID
	AssetName   string
	Quantity    int
}

// ListCurrentCustodyCommand contains the input data required to list current custody.
type ListCurrentCustodyCommand struct {
	PersonnelID domain.PersonnelID
}

// ListCurrentCustodyService handles current custody listing for personnel.
type ListCurrentCustodyService struct {
	personnelRepository ports.PersonnelRepository
	custodyRepository   ports.CustodyRepository
}

// NewListCurrentCustodyService creates a ListCurrentCustodyService with its dependencies.
func NewListCurrentCustodyService(
	personnelRepository ports.PersonnelRepository,
	custodyRepository ports.CustodyRepository,
) *ListCurrentCustodyService {
	return &ListCurrentCustodyService{
		personnelRepository: personnelRepository,
		custodyRepository:   custodyRepository,
	}
}

// Execute retrieves current custody balances for a personnel record.
func (s *ListCurrentCustodyService) Execute(
	ctx context.Context,
	cmd ListCurrentCustodyCommand,
) ([]CurrentCustodyItem, error) {
	if cmd.PersonnelID == "" {
		return nil, domain.ErrEmptyPersonnelID
	}

	if _, err := s.personnelRepository.FindByID(ctx, cmd.PersonnelID); err != nil {
		return nil, err
	}

	items, err := s.custodyRepository.ListCurrentByPersonnel(ctx, cmd.PersonnelID)
	if err != nil {
		return nil, err
	}

	result := make([]CurrentCustodyItem, 0, len(items))

	for _, item := range items {
		result = append(result, CurrentCustodyItem{
			PersonnelID: item.PersonnelID,
			AssetID:     item.AssetID,
			AssetName:   item.AssetName,
			Quantity:    item.Quantity,
		})
	}

	return result, nil
}

const (
	defaultCustodyHistoryLimit = 50
	maxCustodyHistoryLimit     = 100
)

// CustodyHistoryLine contains display-ready asset data inside a custody history entry.
type CustodyHistoryLine struct {
	AssetID   domain.AssetID
	AssetName string
	Quantity  int
}

// CustodyHistoryEntry contains a custody transaction and its lines.
type CustodyHistoryEntry struct {
	ID          domain.CustodyTransactionID
	Type        domain.CustodyTransactionType
	PersonnelID domain.PersonnelID
	Notes       string
	CreatedAt   time.Time
	Lines       []CustodyHistoryLine
}

// ListCustodyHistoryCommand contains the input data required to list custody history.
type ListCustodyHistoryCommand struct {
	PersonnelID domain.PersonnelID
	Limit       int
}

// ListCustodyHistoryService handles custody history listing for personnel.
type ListCustodyHistoryService struct {
	personnelRepository ports.PersonnelRepository
	custodyRepository   ports.CustodyRepository
}

// NewListCustodyHistoryService creates a ListCustodyHistoryService with its dependencies.
func NewListCustodyHistoryService(
	personnelRepository ports.PersonnelRepository,
	custodyRepository ports.CustodyRepository,
) *ListCustodyHistoryService {
	return &ListCustodyHistoryService{
		personnelRepository: personnelRepository,
		custodyRepository:   custodyRepository,
	}
}

// Execute retrieves custody history for a personnel record.
func (s *ListCustodyHistoryService) Execute(
	ctx context.Context,
	cmd ListCustodyHistoryCommand,
) ([]CustodyHistoryEntry, error) {
	if cmd.PersonnelID == "" {
		return nil, domain.ErrEmptyPersonnelID
	}

	if _, err := s.personnelRepository.FindByID(ctx, cmd.PersonnelID); err != nil {
		return nil, err
	}

	limit := cmd.Limit
	if limit <= 0 {
		limit = defaultCustodyHistoryLimit
	}

	if limit > maxCustodyHistoryLimit {
		limit = maxCustodyHistoryLimit
	}

	entries, err := s.custodyRepository.ListHistoryByPersonnel(ctx, cmd.PersonnelID, limit)
	if err != nil {
		return nil, err
	}

	result := make([]CustodyHistoryEntry, 0, len(entries))

	for _, entry := range entries {
		historyEntry := CustodyHistoryEntry{
			ID:          entry.ID,
			Type:        entry.Type,
			PersonnelID: entry.PersonnelID,
			Notes:       entry.Notes,
			CreatedAt:   entry.CreatedAt,
			Lines:       make([]CustodyHistoryLine, 0, len(entry.Lines)),
		}

		for _, line := range entry.Lines {
			historyEntry.Lines = append(historyEntry.Lines, CustodyHistoryLine{
				AssetID:   line.AssetID,
				AssetName: line.AssetName,
				Quantity:  line.Quantity,
			})
		}

		result = append(result, historyEntry)
	}

	return result, nil
}

// CurrentAssetHolder contains current holder display data for an asset.
type CurrentAssetHolder struct {
	AssetID           domain.AssetID
	PersonnelID       domain.PersonnelID
	PersonnelFullName string
	Quantity          int
}

// ListCurrentAssetHoldersCommand contains the input data required to list current asset holders.
type ListCurrentAssetHoldersCommand struct {
	AssetID domain.AssetID
}

// ListCurrentAssetHoldersService handles current holder listing for assets.
type ListCurrentAssetHoldersService struct {
	assetRepository   ports.AssetRepository
	custodyRepository ports.CustodyRepository
}

// NewListCurrentAssetHoldersService creates a ListCurrentAssetHoldersService with its dependencies.
func NewListCurrentAssetHoldersService(
	assetRepository ports.AssetRepository,
	custodyRepository ports.CustodyRepository,
) *ListCurrentAssetHoldersService {
	return &ListCurrentAssetHoldersService{
		assetRepository:   assetRepository,
		custodyRepository: custodyRepository,
	}
}

// Execute retrieves current custody holders for an asset record.
func (s *ListCurrentAssetHoldersService) Execute(
	ctx context.Context,
	cmd ListCurrentAssetHoldersCommand,
) ([]CurrentAssetHolder, error) {
	if cmd.AssetID == "" {
		return nil, domain.ErrEmptyAssetID
	}

	if _, err := s.assetRepository.FindByID(ctx, cmd.AssetID); err != nil {
		return nil, err
	}

	holders, err := s.custodyRepository.ListCurrentByAsset(ctx, cmd.AssetID)
	if err != nil {
		return nil, err
	}

	result := make([]CurrentAssetHolder, 0, len(holders))

	for _, holder := range holders {
		result = append(result, CurrentAssetHolder{
			AssetID:           holder.AssetID,
			PersonnelID:       holder.PersonnelID,
			PersonnelFullName: holder.PersonnelFullName,
			Quantity:          holder.Quantity,
		})
	}

	return result, nil
}
