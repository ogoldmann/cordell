package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CustodyRepository persists custody transactions and balances in PostgreSQL.
type CustodyRepository struct {
	pool    *pgxpool.Pool
	queries *db.Queries
}

// NewCustodyRepository creates a PostgreSQL custody repository.
func NewCustodyRepository(pool *pgxpool.Pool, queries *db.Queries) *CustodyRepository {
	return &CustodyRepository{
		pool:    pool,
		queries: queries,
	}
}

// SaveTransaction persists a custody transaction and updates balances atomically.
func (r *CustodyRepository) SaveTransaction(ctx context.Context, transaction domain.CustodyTransaction) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	_, err = qtx.CreateCustodyTransaction(ctx, db.CreateCustodyTransactionParams{
		ID:              string(transaction.ID()),
		TransactionType: string(transaction.Type()),
		PersonnelID:     string(transaction.PersonnelID()),
		OperatorID:      string(transaction.OperatorID()),
		Notes:           transaction.Notes(),
	})
	if err != nil {
		if isUniqueViolation(err, "custody_transactions_pkey") {
			return false, nil
		}

		return false, err
	}

	for _, line := range transaction.Lines() {
		_, err := qtx.CreateCustodyLine(ctx, db.CreateCustodyLineParams{
			CustodyTransactionID: string(transaction.ID()),
			AssetID:              string(line.AssetID()),
			Quantity:             int32(line.Quantity().Int()),
		})
		if err != nil {
			return false, err
		}

		if err := applyCustodyBalanceChange(ctx, qtx, transaction, line); err != nil {
			return false, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func applyCustodyBalanceChange(
	ctx context.Context,
	queries *db.Queries,
	transaction domain.CustodyTransaction,
	line domain.CustodyLine,
) error {
	switch transaction.Type() {
	case domain.CustodyTransactionTypeCheckout:
		return queries.IncreaseCustodyBalance(ctx, db.IncreaseCustodyBalanceParams{
			PersonnelID: string(transaction.PersonnelID()),
			AssetID:     string(line.AssetID()),
			Quantity:    int32(line.Quantity().Int()),
		})

	case domain.CustodyTransactionTypeReturn:
		_, err := queries.DecreaseCustodyBalanceIfAvailable(ctx, db.DecreaseCustodyBalanceIfAvailableParams{
			PersonnelID: string(transaction.PersonnelID()),
			AssetID:     string(line.AssetID()),
			Quantity:    int32(line.Quantity().Int()),
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return domain.ErrInsufficientCustodyBalance
			}

			return err
		}

		return nil

	default:
		return domain.ErrInvalidTransactionType
	}
}

// SaveCorrection persists a custody correction and updates balances atomically.
func (r *CustodyRepository) SaveCorrection(
	ctx context.Context,
	correction domain.CustodyCorrection,
	transactionType domain.CustodyTransactionType,
	previousPersonnelID domain.PersonnelID,
	previousLines []domain.CustodyLine,
) (bool, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`
INSERT INTO custody_corrections (
    id,
    corrected_transaction_id,
    operator_id,
    corrected_personnel_id,
    corrected_notes
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)`,
		string(correction.ID()),
		string(correction.CorrectedTransactionID()),
		string(correction.OperatorID()),
		string(correction.CorrectedPersonnelID()),
		correction.CorrectedNotes(),
	)
	if err != nil {
		if isUniqueViolation(err, "custody_corrections_pkey") {
			return false, nil
		}

		return false, err
	}

	for _, line := range correction.Lines() {
		_, err := tx.Exec(
			ctx,
			`
INSERT INTO custody_correction_lines (
    custody_correction_id,
    asset_id,
    quantity
) VALUES (
    $1,
    $2,
    $3
)`,
			string(correction.ID()),
			string(line.AssetID()),
			line.Quantity().Int(),
		)
		if err != nil {
			return false, err
		}
	}

	qtx := r.queries.WithTx(tx)
	for _, delta := range custodyCorrectionBalanceDeltas(
		transactionType,
		previousPersonnelID,
		previousLines,
		correction.CorrectedPersonnelID(),
		correction.Lines(),
	) {
		if err := applyCustodyBalanceDelta(ctx, qtx, delta); err != nil {
			return false, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

type custodyBalanceDelta struct {
	personnelID domain.PersonnelID
	assetID     domain.AssetID
	quantity    int
}

type custodyBalanceDeltaKey struct {
	personnelID domain.PersonnelID
	assetID     domain.AssetID
}

func custodyCorrectionBalanceDeltas(
	transactionType domain.CustodyTransactionType,
	previousPersonnelID domain.PersonnelID,
	previousLines []domain.CustodyLine,
	correctedPersonnelID domain.PersonnelID,
	correctedLines []domain.CustodyLine,
) []custodyBalanceDelta {
	multiplier := custodyTransactionBalanceMultiplier(transactionType)
	totals := make(map[custodyBalanceDeltaKey]int)
	orderedKeys := make([]custodyBalanceDeltaKey, 0, len(previousLines)+len(correctedLines))

	addDelta := func(personnelID domain.PersonnelID, assetID domain.AssetID, quantity int) {
		key := custodyBalanceDeltaKey{
			personnelID: personnelID,
			assetID:     assetID,
		}
		if _, ok := totals[key]; !ok {
			orderedKeys = append(orderedKeys, key)
		}

		totals[key] += quantity
	}

	for _, line := range previousLines {
		addDelta(previousPersonnelID, line.AssetID(), -multiplier*line.Quantity().Int())
	}

	for _, line := range correctedLines {
		addDelta(correctedPersonnelID, line.AssetID(), multiplier*line.Quantity().Int())
	}

	deltas := make([]custodyBalanceDelta, 0, len(orderedKeys))
	for _, key := range orderedKeys {
		quantity := totals[key]
		if quantity == 0 {
			continue
		}

		deltas = append(deltas, custodyBalanceDelta{
			personnelID: key.personnelID,
			assetID:     key.assetID,
			quantity:    quantity,
		})
	}

	return deltas
}

func custodyTransactionBalanceMultiplier(transactionType domain.CustodyTransactionType) int {
	if transactionType == domain.CustodyTransactionTypeReturn {
		return -1
	}

	return 1
}

func applyCustodyBalanceDelta(ctx context.Context, queries *db.Queries, delta custodyBalanceDelta) error {
	if delta.quantity > 0 {
		return queries.IncreaseCustodyBalance(ctx, db.IncreaseCustodyBalanceParams{
			PersonnelID: string(delta.personnelID),
			AssetID:     string(delta.assetID),
			Quantity:    int32(delta.quantity),
		})
	}

	_, err := queries.DecreaseCustodyBalanceIfAvailable(ctx, db.DecreaseCustodyBalanceIfAvailableParams{
		PersonnelID: string(delta.personnelID),
		AssetID:     string(delta.assetID),
		Quantity:    int32(-delta.quantity),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.ErrInsufficientCustodyBalance
		}

		return err
	}

	return nil
}

// CurrentQuantity returns the current custody quantity for a personnel and asset pair.
func (r *CustodyRepository) CurrentQuantity(
	ctx context.Context,
	personnelID domain.PersonnelID,
	assetID domain.AssetID,
) (int, error) {
	quantity, err := r.queries.GetCustodyBalanceQuantity(ctx, db.GetCustodyBalanceQuantityParams{
		PersonnelID: string(personnelID),
		AssetID:     string(assetID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}

		return 0, err
	}

	return int(quantity), nil
}

var _ ports.CustodyRepository = (*CustodyRepository)(nil)

// ListPersonnelWithCurrentCustody returns personnel that currently hold at least one asset.
func (r *CustodyRepository) ListPersonnelWithCurrentCustody(
	ctx context.Context,
) ([]ports.PersonnelWithCurrentCustody, error) {
	rows, err := r.queries.ListPersonnelWithCurrentCustody(ctx)
	if err != nil {
		return nil, err
	}

	personnel := make([]ports.PersonnelWithCurrentCustody, 0, len(rows))

	for _, row := range rows {
		personnel = append(personnel, ports.PersonnelWithCurrentCustody{
			ID:             domain.PersonnelID(row.ID),
			FullName:       row.FullName,
			Alias:          row.Alias,
			Rank:           domain.Rank(row.Rank),
			RegistrationID: domain.RegistrationID(row.RegistrationID),
			Active:         row.Active,
			TotalQuantity:  int(row.TotalQuantity),
		})
	}

	return personnel, nil
}

// ListCurrentByPersonnel retrieves current custody balances for a personnel record.
func (r *CustodyRepository) ListCurrentByPersonnel(
	ctx context.Context,
	personnelID domain.PersonnelID,
) ([]ports.CurrentCustodyItem, error) {
	rows, err := r.queries.ListCurrentCustodyByPersonnel(ctx, string(personnelID))
	if err != nil {
		return nil, err
	}

	items := make([]ports.CurrentCustodyItem, 0, len(rows))

	for _, row := range rows {
		items = append(items, ports.CurrentCustodyItem{
			PersonnelID: domain.PersonnelID(row.PersonnelID),
			AssetID:     domain.AssetID(row.AssetID),
			AssetName:   row.AssetName,
			AssetActive: row.AssetActive,
			Quantity:    int(row.Quantity),
		})
	}

	return items, nil
}

// ListCurrentByAsset retrieves current custody holders for an asset record.
func (r *CustodyRepository) ListCurrentByAsset(
	ctx context.Context,
	assetID domain.AssetID,
) ([]ports.CurrentAssetHolder, error) {
	rows, err := r.queries.ListCurrentCustodyByAsset(ctx, string(assetID))
	if err != nil {
		return nil, err
	}

	holders := make([]ports.CurrentAssetHolder, 0, len(rows))

	for _, row := range rows {
		holders = append(holders, ports.CurrentAssetHolder{
			AssetID:           domain.AssetID(row.AssetID),
			PersonnelID:       domain.PersonnelID(row.PersonnelID),
			PersonnelFullName: row.PersonnelFullName,
			PersonnelAlias:    row.PersonnelAlias,
			PersonnelRank:     domain.Rank(row.PersonnelRank),
			PersonnelActive:   row.PersonnelActive,
			Quantity:          int(row.Quantity),
		})
	}

	return holders, nil
}

// ListHistoryByPersonnel retrieves custody transaction history for a personnel record.
func (r *CustodyRepository) ListHistoryByPersonnel(
	ctx context.Context,
	personnelID domain.PersonnelID,
	limit int,
) ([]ports.CustodyHistoryEntry, error) {
	rows, err := r.queries.ListCustodyHistoryByPersonnel(ctx, db.ListCustodyHistoryByPersonnelParams{
		PersonnelID: string(personnelID),
		LimitCount:  int32(limit),
	})
	if err != nil {
		return nil, err
	}

	entries := make([]ports.CustodyHistoryEntry, 0)
	entryIndexes := make(map[string]int)

	for _, row := range rows {
		transactionID := row.TransactionID

		entryIndex, ok := entryIndexes[transactionID]
		if !ok {
			createdAt, err := timestamptzToTime(row.TransactionCreatedAt)
			if err != nil {
				return nil, err
			}

			entries = append(entries, ports.CustodyHistoryEntry{
				ID:            domain.CustodyTransactionID(row.TransactionID),
				Type:          domain.CustodyTransactionType(row.TransactionType),
				PersonnelID:   domain.PersonnelID(row.PersonnelID),
				OperatorID:    domain.OperatorID(row.OperatorID),
				OperatorAlias: row.OperatorAlias,
				OperatorRank:  domain.Rank(row.OperatorRank),
				Notes:         row.Notes,
				CreatedAt:     createdAt,
				Lines:         make([]ports.CustodyHistoryLine, 0),
				HasCorrection: row.HasCorrection,
				EditCount:     int(row.EditCount),
			})

			entryIndex = len(entries) - 1
			entryIndexes[transactionID] = entryIndex
		}

		entries[entryIndex].Lines = append(entries[entryIndex].Lines, ports.CustodyHistoryLine{
			AssetID:   domain.AssetID(row.AssetID),
			AssetName: row.AssetName,
			Quantity:  int(row.Quantity),
		})
	}

	return entries, nil
}

// ListAssetCustodyHistory retrieves effective custody transaction history for an asset record.
func (r *CustodyRepository) ListAssetCustodyHistory(
	ctx context.Context,
	assetID domain.AssetID,
) ([]ports.AssetCustodyHistoryItem, error) {
	rows, err := r.queries.ListAssetCustodyHistoryRows(ctx, string(assetID))
	if err != nil {
		return nil, err
	}

	itemsByID := make(map[domain.CustodyTransactionID]*ports.AssetCustodyHistoryItem)
	itemOrder := make([]domain.CustodyTransactionID, 0)

	for _, row := range rows {
		transactionID := domain.CustodyTransactionID(row.ID)

		item, ok := itemsByID[transactionID]
		if !ok {
			createdAt, err := timestamptzToTime(row.CreatedAt)
			if err != nil {
				return nil, err
			}

			item = &ports.AssetCustodyHistoryItem{
				ID:                transactionID,
				Sequence:          int(row.SequenceNumber),
				Type:              domain.CustodyTransactionType(row.TransactionType),
				CreatedAt:         createdAt,
				PersonnelID:       domain.PersonnelID(row.PersonnelID),
				PersonnelRank:     domain.PersonnelRank(row.PersonnelRank),
				PersonnelAlias:    row.PersonnelAlias,
				PersonnelFullName: row.PersonnelFullName,
				OperatorID:        domain.OperatorID(row.OperatorID),
				OperatorRank:      domain.Rank(row.OperatorRank),
				OperatorAlias:     row.OperatorAlias,
				EditCount:         int(row.EditCount),
				Notes:             row.Notes,
				Lines:             make([]ports.AssetCustodyHistoryLine, 0),
			}

			itemsByID[transactionID] = item
			itemOrder = append(itemOrder, transactionID)
		}

		item.Lines = append(item.Lines, ports.AssetCustodyHistoryLine{
			AssetID:   domain.AssetID(row.AssetID),
			AssetName: row.AssetName,
			Quantity:  int(row.Quantity),
		})
	}

	items := make([]ports.AssetCustodyHistoryItem, 0, len(itemOrder))
	for _, id := range itemOrder {
		items = append(items, *itemsByID[id])
	}

	return items, nil
}

// ListTransactionLedgerPeriods returns available year/month periods for the custody ledger.
func (r *CustodyRepository) ListTransactionLedgerPeriods(
	ctx context.Context,
) ([]ports.CustodyTransactionLedgerPeriod, error) {
	rows, err := r.queries.ListCustodyTransactionLedgerPeriods(ctx)
	if err != nil {
		return nil, err
	}

	periods := make([]ports.CustodyTransactionLedgerPeriod, 0, len(rows))

	for _, row := range rows {
		periods = append(periods, ports.CustodyTransactionLedgerPeriod{
			Year:             int(row.Year),
			Month:            int(row.Month),
			TransactionCount: int(row.TransactionCount),
		})
	}

	return periods, nil
}

// ListTransactionSummaries returns custody transaction summaries for the ledger.
func (r *CustodyRepository) ListTransactionSummaries(
	ctx context.Context,
	filters ports.CustodyTransactionSummaryFilters,
) (ports.CustodyTransactionSummaryPage, error) {
	pageSize := filters.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}

	if pageSize > 100 {
		pageSize = 100
	}

	offset := filters.Offset
	if offset < 0 {
		offset = 0
	}

	transactionTypeFilter := filters.TransactionTypeFilter
	if transactionTypeFilter == "" {
		transactionTypeFilter = ports.CustodyTransactionTypeFilterAll
	}

	editStatusFilter := filters.EditStatusFilter
	if editStatusFilter == "" {
		editStatusFilter = ports.CustodyEditStatusFilterAll
	}

	searchPattern := ""
	searchQuery := strings.TrimSpace(filters.SearchQuery)
	if searchQuery != "" {
		searchPattern = "%" + escapeLikePattern(searchQuery) + "%"
	}

	rows, err := r.queries.ListCustodyTransactionSummaries(ctx, db.ListCustodyTransactionSummariesParams{
		PageSizePlusOne:       int32(pageSize + 1),
		OffsetCount:           int32(offset),
		SearchPattern:         searchPattern,
		TransactionTypeFilter: string(transactionTypeFilter),
		EditStatusFilter:      string(editStatusFilter),
		HasPeriod:             filters.HasPeriod,
		PeriodStart: pgtype.Timestamptz{
			Time:  filters.PeriodStart,
			Valid: filters.HasPeriod,
		},
		PeriodEnd: pgtype.Timestamptz{
			Time:  filters.PeriodEnd,
			Valid: filters.HasPeriod,
		},
	})
	if err != nil {
		return ports.CustodyTransactionSummaryPage{}, err
	}

	summariesByID := make(map[domain.CustodyTransactionID]*ports.CustodyTransactionSummary)
	orderedIDs := make([]domain.CustodyTransactionID, 0)

	for _, row := range rows {
		transactionID := domain.CustodyTransactionID(row.ID)

		summary, ok := summariesByID[transactionID]
		if !ok {
			createdAt, err := timestamptzToTime(row.CreatedAt)
			if err != nil {
				return ports.CustodyTransactionSummaryPage{}, err
			}

			summary = &ports.CustodyTransactionSummary{
				ID:                         transactionID,
				SequenceNumber:             int(row.SequenceNumber),
				TransactionType:            domain.CustodyTransactionType(row.TransactionType),
				OriginalPersonnelID:        domain.PersonnelID(row.OriginalPersonnelID),
				OriginalPersonnelFullName:  row.OriginalPersonnelFullName,
				OriginalPersonnelAlias:     row.OriginalPersonnelAlias,
				OriginalPersonnelRank:      domain.Rank(row.OriginalPersonnelRank),
				OriginalPersonnelActive:    row.OriginalPersonnelActive,
				EffectivePersonnelID:       domain.PersonnelID(row.EffectivePersonnelID),
				EffectivePersonnelFullName: row.EffectivePersonnelFullName,
				EffectivePersonnelAlias:    row.EffectivePersonnelAlias,
				EffectivePersonnelRank:     domain.Rank(row.EffectivePersonnelRank),
				EffectivePersonnelActive:   row.EffectivePersonnelActive,
				OperatorID:                 domain.OperatorID(row.OperatorID),
				OperatorAlias:              row.OperatorAlias,
				OperatorRank:               domain.Rank(row.OperatorRank),
				OperatorRole:               domain.OperatorRole(row.OperatorRole),
				OperatorActive:             row.OperatorActive,
				TotalQuantity:              int(row.TotalQuantity),
				CreatedAt:                  createdAt,
				HasCorrection:              row.HasCorrection,
				EditCount:                  int(row.EditCount),
				Lines:                      make([]ports.CustodyTransactionSummaryLine, 0),
			}

			summariesByID[transactionID] = summary
			orderedIDs = append(orderedIDs, transactionID)
		}

		if row.AssetID.Valid {
			summary.Lines = append(summary.Lines, ports.CustodyTransactionSummaryLine{
				AssetID:     domain.AssetID(row.AssetID.String),
				AssetName:   row.AssetName.String,
				AssetActive: row.AssetActive.Bool,
				Quantity:    int(row.Quantity.Int32),
			})
		}
	}

	summaries := make([]ports.CustodyTransactionSummary, 0, len(orderedIDs))

	for _, id := range orderedIDs {
		summaries = append(summaries, *summariesByID[id])
	}

	hasNextPage := len(summaries) > pageSize
	if hasNextPage {
		summaries = summaries[:pageSize]
	}

	return ports.CustodyTransactionSummaryPage{
		Items:       summaries,
		HasNextPage: hasNextPage,
	}, nil
}

// FindReceiptByID retrieves a complete custody transaction receipt.
func (r *CustodyRepository) FindReceiptByID(
	ctx context.Context,
	id domain.CustodyTransactionID,
) (ports.CustodyReceipt, error) {
	rows, err := r.queries.GetCustodyTransactionReceiptByID(ctx, string(id))
	if err != nil {
		return ports.CustodyReceipt{}, err
	}

	if len(rows) == 0 {
		return ports.CustodyReceipt{}, ports.ErrNotFound
	}

	first := rows[0]

	createdAt, err := timestamptzToTime(first.CreatedAt)
	if err != nil {
		return ports.CustodyReceipt{}, err
	}

	receipt := ports.CustodyReceipt{
		ID:                      domain.CustodyTransactionID(first.ID),
		TransactionType:         domain.CustodyTransactionType(first.TransactionType),
		PersonnelID:             domain.PersonnelID(first.PersonnelID),
		PersonnelFullName:       first.PersonnelFullName,
		PersonnelAlias:          first.PersonnelAlias,
		PersonnelRank:           domain.Rank(first.PersonnelRank),
		PersonnelRegistrationID: domain.RegistrationID(first.PersonnelRegistrationID),
		PersonnelActive:         first.PersonnelActive,
		OperatorID:              domain.OperatorID(first.OperatorID),
		OperatorAlias:           first.OperatorAlias,
		OperatorRank:            domain.Rank(first.OperatorRank),
		OperatorRole:            domain.OperatorRole(first.OperatorRole),
		OperatorActive:          first.OperatorActive,
		Notes:                   first.Notes,
		CreatedAt:               createdAt,
		Lines:                   make([]ports.CustodyReceiptLine, 0, len(rows)),
	}

	for _, row := range rows {
		receipt.Lines = append(receipt.Lines, ports.CustodyReceiptLine{
			AssetID:     domain.AssetID(row.AssetID),
			AssetName:   row.AssetName,
			AssetActive: row.AssetActive,
			Quantity:    int(row.Quantity),
		})
	}

	if err := r.attachLatestCorrection(ctx, &receipt); err != nil {
		return ports.CustodyReceipt{}, err
	}

	return receipt, nil
}

// ListCorrectionContextsByTransactionID returns all correction contexts for a transaction in chronological order.
func (r *CustodyRepository) ListCorrectionContextsByTransactionID(
	ctx context.Context,
	id domain.CustodyTransactionID,
) ([]ports.CustodyCorrectionContext, error) {
	rows, err := r.queries.GetCustodyCorrectionContextsByTransactionID(ctx, string(id))
	if err != nil {
		return nil, err
	}

	contextsByID := make(map[domain.CustodyCorrectionID]*ports.CustodyCorrectionContext)
	orderedIDs := make([]domain.CustodyCorrectionID, 0)

	for _, row := range rows {
		correctionID := domain.CustodyCorrectionID(row.ID)

		context, ok := contextsByID[correctionID]
		if !ok {
			createdAt, err := timestamptzToTime(row.CreatedAt)
			if err != nil {
				return nil, err
			}

			context = &ports.CustodyCorrectionContext{
				ID:                               correctionID,
				CorrectedTransactionID:           domain.CustodyTransactionID(row.CorrectedTransactionID),
				OperatorID:                       domain.OperatorID(row.OperatorID),
				OperatorAlias:                    row.OperatorAlias,
				OperatorRank:                     domain.Rank(row.OperatorRank),
				OperatorRole:                     domain.OperatorRole(row.OperatorRole),
				OperatorActive:                   row.OperatorActive,
				CorrectedPersonnelID:             domain.PersonnelID(row.CorrectedPersonnelID),
				CorrectedPersonnelFullName:       row.CorrectedPersonnelFullName,
				CorrectedPersonnelAlias:          row.CorrectedPersonnelAlias,
				CorrectedPersonnelRank:           domain.Rank(row.CorrectedPersonnelRank),
				CorrectedPersonnelRegistrationID: domain.RegistrationID(row.CorrectedPersonnelRegistrationID),
				CorrectedPersonnelActive:         row.CorrectedPersonnelActive,
				CorrectedNotes:                   row.CorrectedNotes,
				CreatedAt:                        createdAt,
				Lines:                            make([]ports.CustodyCorrectionContextLine, 0),
			}

			contextsByID[correctionID] = context
			orderedIDs = append(orderedIDs, correctionID)
		}

		context.Lines = append(context.Lines, ports.CustodyCorrectionContextLine{
			AssetID:     domain.AssetID(row.AssetID),
			AssetName:   row.AssetName,
			AssetActive: row.AssetActive,
			Quantity:    int(row.Quantity),
		})
	}

	contexts := make([]ports.CustodyCorrectionContext, 0, len(orderedIDs))

	for _, correctionID := range orderedIDs {
		contexts = append(contexts, *contextsByID[correctionID])
	}

	return contexts, nil
}

func (r *CustodyRepository) attachLatestCorrection(ctx context.Context, receipt *ports.CustodyReceipt) error {
	contexts, err := r.ListCorrectionContextsByTransactionID(ctx, receipt.ID)
	if err != nil {
		return err
	}

	if len(contexts) == 0 {
		return nil
	}

	receipt.HasCorrection = true
	receipt.Correction = contexts[len(contexts)-1]

	return nil
}

func timestamptzToTime(value pgtype.Timestamptz) (time.Time, error) {
	if !value.Valid {
		return time.Time{}, fmt.Errorf("invalid timestamptz value")
	}

	return value.Time, nil
}
