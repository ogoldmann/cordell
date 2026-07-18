package postgres

import (
	"context"
	"errors"
	"fmt"
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
func (r *CustodyRepository) SaveTransaction(ctx context.Context, transaction domain.CustodyTransaction) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
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
		return err
	}

	for _, line := range transaction.Lines() {
		_, err := qtx.CreateCustodyLine(ctx, db.CreateCustodyLineParams{
			CustodyTransactionID: string(transaction.ID()),
			AssetID:              string(line.AssetID()),
			Quantity:             int32(line.Quantity().Int()),
		})
		if err != nil {
			return err
		}

		switch transaction.Type() {
		case domain.CustodyTransactionTypeCheckout:
			err = qtx.IncreaseCustodyBalanceForCheckout(ctx, db.IncreaseCustodyBalanceForCheckoutParams{
				PersonnelID: string(transaction.PersonnelID()),
				AssetID:     string(line.AssetID()),
				Quantity:    int32(line.Quantity().Int()),
			})
			if err != nil {
				return err
			}

		case domain.CustodyTransactionTypeReturn:
			rowsAffected, err := qtx.DecreaseCustodyBalanceForReturn(ctx, db.DecreaseCustodyBalanceForReturnParams{
				PersonnelID: string(transaction.PersonnelID()),
				AssetID:     string(line.AssetID()),
				Quantity:    int32(line.Quantity().Int()),
			})
			if err != nil {
				return err
			}

			if rowsAffected == 0 {
				return domain.ErrInsufficientCustodyBalance
			}
		}
	}

	return tx.Commit(ctx)
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

func timestamptzToTime(value pgtype.Timestamptz) (time.Time, error) {
	if !value.Valid {
		return time.Time{}, fmt.Errorf("invalid timestamptz value")
	}

	return value.Time, nil
}
