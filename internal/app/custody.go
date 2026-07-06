package app

import (
	"context"

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
