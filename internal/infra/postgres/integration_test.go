package postgres_test

import (
	"context"
	"os"
	"testing"

	"cordell/internal/infra/postgres"
	"cordell/internal/infra/postgres/db"

	"github.com/jackc/pgx/v5/pgxpool"
)

func openTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()

	if os.Getenv("CORDELL_INTEGRATION_TESTS") != "1" {
		t.Skip("skipping PostgreSQL integration test; set CORDELL_INTEGRATION_TESTS=1 to run")
	}

	databaseURL := os.Getenv("CORDELL_DATABASE_URL")
	if databaseURL == "" {
		t.Fatal("CORDELL_DATABASE_URL is required for integration tests")
	}

	ctx := context.Background()

	pool, err := postgres.OpenPool(ctx, databaseURL)
	if err != nil {
		t.Fatalf("failed to open PostgreSQL test pool: %v", err)
	}

	t.Cleanup(pool.Close)

	cleanDatabase(t, pool)

	return pool
}

func cleanDatabase(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	_, err := pool.Exec(
		context.Background(),
		`
		TRUNCATE TABLE
			custody_balances,
			custody_lines,
			custody_transactions,
			assets,
			personnel,
			operators
		RESTART IDENTITY CASCADE
		`,
	)
	if err != nil {
		t.Fatalf("failed to clean test database: %v", err)
	}
}

func newTestQueries(pool *pgxpool.Pool) *db.Queries {
	return db.New(pool)
}
