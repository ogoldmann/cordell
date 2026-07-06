package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/app"
	"cordell/internal/infra/postgres"
)

func TestPostgresAssetRepositoryCreateAndFind(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	assetRepository := postgres.NewAssetRepository(queries)
	idGenerator := fixedIDGenerator{id: "asset-1"}

	service := app.NewCreateAssetService(assetRepository, idGenerator)

	created, err := service.Execute(context.Background(), app.CreateAssetCommand{
		Name: "  Radio  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, err := assetRepository.FindByID(context.Background(), created.ID())
	if err != nil {
		t.Fatalf("expected no error finding asset, got %v", err)
	}

	if found.ID() != "asset-1" {
		t.Fatalf("expected asset id asset-1, got %s", found.ID())
	}

	if found.Name() != "Radio" {
		t.Fatalf("expected asset name Radio, got %s", found.Name())
	}

	if !found.Active() {
		t.Fatal("expected asset to be active")
	}
}
