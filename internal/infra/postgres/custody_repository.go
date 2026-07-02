package postgres

import (
	"context"
	"errors"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5"
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
