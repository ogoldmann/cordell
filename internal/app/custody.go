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

// RegisterCustodyTransactionResult is returned after registering a checkout or return.
type RegisterCustodyTransactionResult struct {
	Transaction domain.CustodyTransaction
	Created     bool
}

func normalizeCustodyLineCommands(commands []CustodyLineCommand) []CustodyLineCommand {
	totalsByAssetID := make(map[domain.AssetID]int)
	orderedAssetIDs := make([]domain.AssetID, 0, len(commands))

	for _, command := range commands {
		if _, ok := totalsByAssetID[command.AssetID]; !ok {
			orderedAssetIDs = append(orderedAssetIDs, command.AssetID)
		}

		totalsByAssetID[command.AssetID] += command.Quantity
	}

	normalizedCommands := make([]CustodyLineCommand, 0, len(orderedAssetIDs))

	for _, assetID := range orderedAssetIDs {
		normalizedCommands = append(normalizedCommands, CustodyLineCommand{
			AssetID:  assetID,
			Quantity: totalsByAssetID[assetID],
		})
	}

	return normalizedCommands
}

func buildCustodyLines(
	ctx context.Context,
	assetRepository ports.AssetRepository,
	commands []CustodyLineCommand,
	requireActiveAssets bool,
) ([]domain.CustodyLine, error) {
	normalizedCommands := normalizeCustodyLineCommands(commands)

	lines := make([]domain.CustodyLine, 0, len(normalizedCommands))

	for _, lineCommand := range normalizedCommands {
		quantity, err := domain.NewQuantity(lineCommand.Quantity)
		if err != nil {
			return nil, err
		}

		line, err := domain.NewCustodyLine(lineCommand.AssetID, quantity)
		if err != nil {
			return nil, err
		}

		asset, err := assetRepository.FindByID(ctx, lineCommand.AssetID)
		if err != nil {
			return nil, err
		}

		if requireActiveAssets && !asset.Active() {
			return nil, domain.ErrInactiveAsset
		}

		lines = append(lines, line)
	}

	return lines, nil
}

func custodyTransactionIDFromCommand(
	commandTransactionID domain.CustodyTransactionID,
	idGenerator ports.IDGenerator,
) (domain.CustodyTransactionID, error) {
	if commandTransactionID != "" {
		return commandTransactionID, nil
	}

	id, err := idGenerator.NewID()
	if err != nil {
		return "", err
	}

	return domain.CustodyTransactionID(id), nil
}

// RegisterCheckoutCommand contains the input data required to register an asset checkout.
type RegisterCheckoutCommand struct {
	TransactionID domain.CustodyTransactionID
	PersonnelID   domain.PersonnelID
	OperatorID    domain.OperatorID
	Lines         []CustodyLineCommand
	Notes         string
}

// RegisterCheckoutService handles the asset checkout use case.
type RegisterCheckoutService struct {
	personnelRepository ports.PersonnelRepository
	assetRepository     ports.AssetRepository
	operatorRepository  ports.OperatorRepository
	custodyRepository   ports.CustodyRepository
	idGenerator         ports.IDGenerator
}

// NewRegisterCheckoutService creates a RegisterCheckoutService with its dependencies.
func NewRegisterCheckoutService(
	personnelRepository ports.PersonnelRepository,
	assetRepository ports.AssetRepository,
	operatorRepository ports.OperatorRepository,
	custodyRepository ports.CustodyRepository,
	idGenerator ports.IDGenerator,
) *RegisterCheckoutService {
	return &RegisterCheckoutService{
		personnelRepository: personnelRepository,
		assetRepository:     assetRepository,
		operatorRepository:  operatorRepository,
		custodyRepository:   custodyRepository,
		idGenerator:         idGenerator,
	}
}

// Execute registers and persists a checkout custody transaction.
func (s *RegisterCheckoutService) Execute(
	ctx context.Context,
	cmd RegisterCheckoutCommand,
) (RegisterCustodyTransactionResult, error) {
	if cmd.PersonnelID == "" {
		return RegisterCustodyTransactionResult{}, domain.ErrEmptyPersonnelID
	}

	personnel, err := s.personnelRepository.FindByID(ctx, cmd.PersonnelID)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	if !personnel.Active() {
		return RegisterCustodyTransactionResult{}, domain.ErrInactivePersonnel
	}

	if cmd.OperatorID == "" {
		return RegisterCustodyTransactionResult{}, domain.ErrEmptyOperatorID
	}

	operator, err := s.operatorRepository.FindByID(ctx, cmd.OperatorID)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	if !operator.Active() {
		return RegisterCustodyTransactionResult{}, ports.ErrNotFound
	}

	lines, err := buildCustodyLines(ctx, s.assetRepository, cmd.Lines, true)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	transactionID, err := custodyTransactionIDFromCommand(cmd.TransactionID, s.idGenerator)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	transaction, err := domain.NewCustodyTransaction(
		transactionID,
		domain.CustodyTransactionTypeCheckout,
		cmd.PersonnelID,
		cmd.OperatorID,
		lines,
		cmd.Notes,
	)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	created, err := s.custodyRepository.SaveTransaction(ctx, transaction)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	return RegisterCustodyTransactionResult{
		Transaction: transaction,
		Created:     created,
	}, nil
}

// RegisterReturnCommand contains the input data required to register an asset return.
type RegisterReturnCommand struct {
	TransactionID domain.CustodyTransactionID
	PersonnelID   domain.PersonnelID
	OperatorID    domain.OperatorID
	Lines         []CustodyLineCommand
	Notes         string
}

// RegisterReturnService handles the asset return use case.
type RegisterReturnService struct {
	personnelRepository ports.PersonnelRepository
	assetRepository     ports.AssetRepository
	operatorRepository  ports.OperatorRepository
	custodyRepository   ports.CustodyRepository
	idGenerator         ports.IDGenerator
}

// NewRegisterReturnService creates a RegisterReturnService with its dependencies.
func NewRegisterReturnService(
	personnelRepository ports.PersonnelRepository,
	assetRepository ports.AssetRepository,
	operatorRepository ports.OperatorRepository,
	custodyRepository ports.CustodyRepository,
	idGenerator ports.IDGenerator,
) *RegisterReturnService {
	return &RegisterReturnService{
		personnelRepository: personnelRepository,
		assetRepository:     assetRepository,
		operatorRepository:  operatorRepository,
		custodyRepository:   custodyRepository,
		idGenerator:         idGenerator,
	}
}

// Execute registers and persists a return custody transaction.
func (s *RegisterReturnService) Execute(
	ctx context.Context,
	cmd RegisterReturnCommand,
) (RegisterCustodyTransactionResult, error) {
	if cmd.PersonnelID == "" {
		return RegisterCustodyTransactionResult{}, domain.ErrEmptyPersonnelID
	}

	if _, err := s.personnelRepository.FindByID(ctx, cmd.PersonnelID); err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	if cmd.OperatorID == "" {
		return RegisterCustodyTransactionResult{}, domain.ErrEmptyOperatorID
	}

	operator, err := s.operatorRepository.FindByID(ctx, cmd.OperatorID)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	if !operator.Active() {
		return RegisterCustodyTransactionResult{}, ports.ErrNotFound
	}

	lines, err := buildCustodyLines(ctx, s.assetRepository, cmd.Lines, false)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	if cmd.TransactionID == "" {
		for _, line := range lines {
			currentQuantity, err := s.custodyRepository.CurrentQuantity(ctx, cmd.PersonnelID, line.AssetID())
			if err != nil {
				return RegisterCustodyTransactionResult{}, err
			}

			if currentQuantity < line.Quantity().Int() {
				return RegisterCustodyTransactionResult{}, domain.ErrInsufficientCustodyBalance
			}
		}
	}

	transactionID, err := custodyTransactionIDFromCommand(cmd.TransactionID, s.idGenerator)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	transaction, err := domain.NewCustodyTransaction(
		transactionID,
		domain.CustodyTransactionTypeReturn,
		cmd.PersonnelID,
		cmd.OperatorID,
		lines,
		cmd.Notes,
	)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	created, err := s.custodyRepository.SaveTransaction(ctx, transaction)
	if err != nil {
		return RegisterCustodyTransactionResult{}, err
	}

	return RegisterCustodyTransactionResult{
		Transaction: transaction,
		Created:     created,
	}, nil
}

// CurrentCustodyItem contains current custody display data for application use cases.
type CurrentCustodyItem struct {
	PersonnelID domain.PersonnelID
	AssetID     domain.AssetID
	AssetName   string
	AssetActive bool
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
			AssetActive: item.AssetActive,
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
	ID            domain.CustodyTransactionID
	Type          domain.CustodyTransactionType
	PersonnelID   domain.PersonnelID
	OperatorID    domain.OperatorID
	OperatorAlias string
	OperatorRank  domain.Rank
	Notes         string
	CreatedAt     time.Time
	Lines         []CustodyHistoryLine
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
			ID:            entry.ID,
			Type:          entry.Type,
			PersonnelID:   entry.PersonnelID,
			OperatorID:    entry.OperatorID,
			OperatorAlias: entry.OperatorAlias,
			OperatorRank:  entry.OperatorRank,
			Notes:         entry.Notes,
			CreatedAt:     entry.CreatedAt,
			Lines:         make([]CustodyHistoryLine, 0, len(entry.Lines)),
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

// CustodyReceiptLine contains one line in a custody receipt.
type CustodyReceiptLine struct {
	AssetID     domain.AssetID
	AssetName   string
	AssetActive bool
	Quantity    int
}

// CustodyCorrectionContextLine contains one line in the latest correction context.
type CustodyCorrectionContextLine struct {
	AssetID     domain.AssetID
	AssetName   string
	AssetActive bool
	Quantity    int
}

// CustodyCorrectionContext contains the latest correction attached to a receipt.
type CustodyCorrectionContext struct {
	ID                               domain.CustodyCorrectionID
	CorrectedTransactionID           domain.CustodyTransactionID
	OperatorID                       domain.OperatorID
	OperatorAlias                    string
	OperatorRank                     domain.Rank
	OperatorRole                     domain.OperatorRole
	OperatorActive                   bool
	CorrectedPersonnelID             domain.PersonnelID
	CorrectedPersonnelFullName       string
	CorrectedPersonnelAlias          string
	CorrectedPersonnelRank           domain.Rank
	CorrectedPersonnelRegistrationID domain.RegistrationID
	CorrectedPersonnelActive         bool
	CorrectedNotes                   string
	CreatedAt                        time.Time
	Lines                            []CustodyCorrectionContextLine
}

// CustodyReceipt contains a complete custody transaction receipt.
type CustodyReceipt struct {
	ID                      domain.CustodyTransactionID
	TransactionType         domain.CustodyTransactionType
	PersonnelID             domain.PersonnelID
	PersonnelFullName       string
	PersonnelAlias          string
	PersonnelRank           domain.Rank
	PersonnelRegistrationID domain.RegistrationID
	PersonnelActive         bool
	OperatorID              domain.OperatorID
	OperatorAlias           string
	OperatorRank            domain.Rank
	OperatorRole            domain.OperatorRole
	OperatorActive          bool
	Notes                   string
	CreatedAt               time.Time
	Lines                   []CustodyReceiptLine
	HasCorrection           bool
	Correction              CustodyCorrectionContext
}

// GetCustodyReceiptCommand contains the input data required to retrieve a custody receipt.
type GetCustodyReceiptCommand struct {
	ID domain.CustodyTransactionID
}

// GetCustodyReceiptService handles custody receipt retrieval.
type GetCustodyReceiptService struct {
	custodyRepository ports.CustodyRepository
}

// NewGetCustodyReceiptService creates a GetCustodyReceiptService.
func NewGetCustodyReceiptService(custodyRepository ports.CustodyRepository) *GetCustodyReceiptService {
	return &GetCustodyReceiptService{
		custodyRepository: custodyRepository,
	}
}

// Execute retrieves a complete custody receipt.
func (s *GetCustodyReceiptService) Execute(
	ctx context.Context,
	cmd GetCustodyReceiptCommand,
) (CustodyReceipt, error) {
	if cmd.ID == "" {
		return CustodyReceipt{}, domain.ErrEmptyTransactionID
	}

	receipt, err := s.custodyRepository.FindReceiptByID(ctx, cmd.ID)
	if err != nil {
		return CustodyReceipt{}, err
	}

	result := CustodyReceipt{
		ID:                      receipt.ID,
		TransactionType:         receipt.TransactionType,
		PersonnelID:             receipt.PersonnelID,
		PersonnelFullName:       receipt.PersonnelFullName,
		PersonnelAlias:          receipt.PersonnelAlias,
		PersonnelRank:           receipt.PersonnelRank,
		PersonnelRegistrationID: receipt.PersonnelRegistrationID,
		PersonnelActive:         receipt.PersonnelActive,
		OperatorID:              receipt.OperatorID,
		OperatorAlias:           receipt.OperatorAlias,
		OperatorRank:            receipt.OperatorRank,
		OperatorRole:            receipt.OperatorRole,
		OperatorActive:          receipt.OperatorActive,
		Notes:                   receipt.Notes,
		CreatedAt:               receipt.CreatedAt,
		Lines:                   make([]CustodyReceiptLine, 0, len(receipt.Lines)),
	}

	for _, line := range receipt.Lines {
		result.Lines = append(result.Lines, CustodyReceiptLine{
			AssetID:     line.AssetID,
			AssetName:   line.AssetName,
			AssetActive: line.AssetActive,
			Quantity:    line.Quantity,
		})
	}

	if receipt.HasCorrection {
		result.HasCorrection = true
		result.Correction = CustodyCorrectionContext{
			ID:                               receipt.Correction.ID,
			CorrectedTransactionID:           receipt.Correction.CorrectedTransactionID,
			OperatorID:                       receipt.Correction.OperatorID,
			OperatorAlias:                    receipt.Correction.OperatorAlias,
			OperatorRank:                     receipt.Correction.OperatorRank,
			OperatorRole:                     receipt.Correction.OperatorRole,
			OperatorActive:                   receipt.Correction.OperatorActive,
			CorrectedPersonnelID:             receipt.Correction.CorrectedPersonnelID,
			CorrectedPersonnelFullName:       receipt.Correction.CorrectedPersonnelFullName,
			CorrectedPersonnelAlias:          receipt.Correction.CorrectedPersonnelAlias,
			CorrectedPersonnelRank:           receipt.Correction.CorrectedPersonnelRank,
			CorrectedPersonnelRegistrationID: receipt.Correction.CorrectedPersonnelRegistrationID,
			CorrectedPersonnelActive:         receipt.Correction.CorrectedPersonnelActive,
			CorrectedNotes:                   receipt.Correction.CorrectedNotes,
			CreatedAt:                        receipt.Correction.CreatedAt,
			Lines:                            make([]CustodyCorrectionContextLine, 0, len(receipt.Correction.Lines)),
		}

		for _, line := range receipt.Correction.Lines {
			result.Correction.Lines = append(result.Correction.Lines, CustodyCorrectionContextLine{
				AssetID:     line.AssetID,
				AssetName:   line.AssetName,
				AssetActive: line.AssetActive,
				Quantity:    line.Quantity,
			})
		}
	}

	return result, nil
}

// RegisterCustodyCorrectionCommand contains the input data required to register a custody correction.
type RegisterCustodyCorrectionCommand struct {
	CorrectionID           domain.CustodyCorrectionID
	CorrectedTransactionID domain.CustodyTransactionID
	OperatorID             domain.OperatorID
	CorrectedPersonnelID   domain.PersonnelID
	Lines                  []CustodyLineCommand
	CorrectedNotes         string
}

// RegisterCustodyCorrectionResult is returned after registering a custody correction.
type RegisterCustodyCorrectionResult struct {
	Correction domain.CustodyCorrection
	Created    bool
}

// RegisterCustodyCorrectionService handles the append-only custody correction use case.
type RegisterCustodyCorrectionService struct {
	personnelRepository ports.PersonnelRepository
	assetRepository     ports.AssetRepository
	operatorRepository  ports.OperatorRepository
	custodyRepository   ports.CustodyRepository
}

// NewRegisterCustodyCorrectionService creates a RegisterCustodyCorrectionService with its dependencies.
func NewRegisterCustodyCorrectionService(
	personnelRepository ports.PersonnelRepository,
	assetRepository ports.AssetRepository,
	operatorRepository ports.OperatorRepository,
	custodyRepository ports.CustodyRepository,
) *RegisterCustodyCorrectionService {
	return &RegisterCustodyCorrectionService{
		personnelRepository: personnelRepository,
		assetRepository:     assetRepository,
		operatorRepository:  operatorRepository,
		custodyRepository:   custodyRepository,
	}
}

// Execute registers an append-only correction and applies custody balance deltas.
func (s *RegisterCustodyCorrectionService) Execute(
	ctx context.Context,
	cmd RegisterCustodyCorrectionCommand,
) (RegisterCustodyCorrectionResult, error) {
	if cmd.CorrectionID == "" {
		return RegisterCustodyCorrectionResult{}, domain.ErrEmptyCustodyCorrectionID
	}

	receipt, err := s.GetCorrectionBaseReceipt(ctx, cmd.CorrectedTransactionID)
	if err != nil {
		return RegisterCustodyCorrectionResult{}, err
	}

	if cmd.OperatorID == "" {
		return RegisterCustodyCorrectionResult{}, domain.ErrEmptyOperatorID
	}

	operator, err := s.operatorRepository.FindByID(ctx, cmd.OperatorID)
	if err != nil {
		return RegisterCustodyCorrectionResult{}, err
	}

	if !operator.Active() {
		return RegisterCustodyCorrectionResult{}, domain.ErrInactiveOperator
	}

	if cmd.CorrectedPersonnelID == "" {
		return RegisterCustodyCorrectionResult{}, domain.ErrEmptyPersonnelID
	}

	if _, err := s.personnelRepository.FindByID(ctx, cmd.CorrectedPersonnelID); err != nil {
		return RegisterCustodyCorrectionResult{}, err
	}

	lines, err := buildCustodyLines(ctx, s.assetRepository, cmd.Lines, false)
	if err != nil {
		return RegisterCustodyCorrectionResult{}, err
	}

	correction, err := domain.NewCustodyCorrection(
		cmd.CorrectionID,
		cmd.CorrectedTransactionID,
		cmd.OperatorID,
		cmd.CorrectedPersonnelID,
		lines,
		cmd.CorrectedNotes,
	)
	if err != nil {
		return RegisterCustodyCorrectionResult{}, err
	}

	previousPersonnelID := receipt.PersonnelID
	previousLines := custodyReceiptLinesToDomainLines(receipt.Lines)

	if receipt.HasCorrection {
		previousPersonnelID = receipt.Correction.CorrectedPersonnelID
		previousLines = custodyCorrectionContextLinesToDomainLines(receipt.Correction.Lines)
	}

	created, err := s.custodyRepository.SaveCorrection(
		ctx,
		correction,
		receipt.TransactionType,
		previousPersonnelID,
		previousLines,
	)
	if err != nil {
		return RegisterCustodyCorrectionResult{}, err
	}

	return RegisterCustodyCorrectionResult{
		Correction: correction,
		Created:    created,
	}, nil
}

// GetCorrectionBaseReceipt retrieves the receipt used as correction base.
func (s *RegisterCustodyCorrectionService) GetCorrectionBaseReceipt(
	ctx context.Context,
	transactionID domain.CustodyTransactionID,
) (CustodyReceipt, error) {
	if transactionID == "" {
		return CustodyReceipt{}, domain.ErrEmptyTransactionID
	}

	receipt, err := s.custodyRepository.FindReceiptByID(ctx, transactionID)
	if err != nil {
		return CustodyReceipt{}, err
	}

	return portCustodyReceiptToAppReceipt(receipt), nil
}

func portCustodyReceiptToAppReceipt(receipt ports.CustodyReceipt) CustodyReceipt {
	result := CustodyReceipt{
		ID:                      receipt.ID,
		TransactionType:         receipt.TransactionType,
		PersonnelID:             receipt.PersonnelID,
		PersonnelFullName:       receipt.PersonnelFullName,
		PersonnelAlias:          receipt.PersonnelAlias,
		PersonnelRank:           receipt.PersonnelRank,
		PersonnelRegistrationID: receipt.PersonnelRegistrationID,
		PersonnelActive:         receipt.PersonnelActive,
		OperatorID:              receipt.OperatorID,
		OperatorAlias:           receipt.OperatorAlias,
		OperatorRank:            receipt.OperatorRank,
		OperatorRole:            receipt.OperatorRole,
		OperatorActive:          receipt.OperatorActive,
		Notes:                   receipt.Notes,
		CreatedAt:               receipt.CreatedAt,
		Lines:                   make([]CustodyReceiptLine, 0, len(receipt.Lines)),
	}

	for _, line := range receipt.Lines {
		result.Lines = append(result.Lines, CustodyReceiptLine{
			AssetID:     line.AssetID,
			AssetName:   line.AssetName,
			AssetActive: line.AssetActive,
			Quantity:    line.Quantity,
		})
	}

	if receipt.HasCorrection {
		result.HasCorrection = true
		result.Correction = CustodyCorrectionContext{
			ID:                               receipt.Correction.ID,
			CorrectedTransactionID:           receipt.Correction.CorrectedTransactionID,
			OperatorID:                       receipt.Correction.OperatorID,
			OperatorAlias:                    receipt.Correction.OperatorAlias,
			OperatorRank:                     receipt.Correction.OperatorRank,
			OperatorRole:                     receipt.Correction.OperatorRole,
			OperatorActive:                   receipt.Correction.OperatorActive,
			CorrectedPersonnelID:             receipt.Correction.CorrectedPersonnelID,
			CorrectedPersonnelFullName:       receipt.Correction.CorrectedPersonnelFullName,
			CorrectedPersonnelAlias:          receipt.Correction.CorrectedPersonnelAlias,
			CorrectedPersonnelRank:           receipt.Correction.CorrectedPersonnelRank,
			CorrectedPersonnelRegistrationID: receipt.Correction.CorrectedPersonnelRegistrationID,
			CorrectedPersonnelActive:         receipt.Correction.CorrectedPersonnelActive,
			CorrectedNotes:                   receipt.Correction.CorrectedNotes,
			CreatedAt:                        receipt.Correction.CreatedAt,
			Lines:                            make([]CustodyCorrectionContextLine, 0, len(receipt.Correction.Lines)),
		}

		for _, line := range receipt.Correction.Lines {
			result.Correction.Lines = append(result.Correction.Lines, CustodyCorrectionContextLine{
				AssetID:     line.AssetID,
				AssetName:   line.AssetName,
				AssetActive: line.AssetActive,
				Quantity:    line.Quantity,
			})
		}
	}

	return result
}

func custodyReceiptLinesToDomainLines(lines []CustodyReceiptLine) []domain.CustodyLine {
	result := make([]domain.CustodyLine, 0, len(lines))

	for _, line := range lines {
		quantity, err := domain.NewQuantity(line.Quantity)
		if err != nil {
			continue
		}

		custodyLine, err := domain.NewCustodyLine(line.AssetID, quantity)
		if err != nil {
			continue
		}

		result = append(result, custodyLine)
	}

	return result
}

func custodyCorrectionContextLinesToDomainLines(lines []CustodyCorrectionContextLine) []domain.CustodyLine {
	result := make([]domain.CustodyLine, 0, len(lines))

	for _, line := range lines {
		quantity, err := domain.NewQuantity(line.Quantity)
		if err != nil {
			continue
		}

		custodyLine, err := domain.NewCustodyLine(line.AssetID, quantity)
		if err != nil {
			continue
		}

		result = append(result, custodyLine)
	}

	return result
}

// CurrentAssetHolder contains current holder display data for an asset.
type CurrentAssetHolder struct {
	AssetID           domain.AssetID
	PersonnelID       domain.PersonnelID
	PersonnelFullName string
	PersonnelAlias    string
	PersonnelRank     domain.Rank
	PersonnelActive   bool
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
			PersonnelAlias:    holder.PersonnelAlias,
			PersonnelRank:     holder.PersonnelRank,
			PersonnelActive:   holder.PersonnelActive,
			Quantity:          holder.Quantity,
		})
	}

	return result, nil
}
