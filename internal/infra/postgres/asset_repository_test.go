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

func TestPostgresAssetRepositoryList(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	assetRepository := postgres.NewAssetRepository(queries)

	firstService := app.NewCreateAssetService(
		assetRepository,
		fixedIDGenerator{id: "asset-1"},
	)
	secondService := app.NewCreateAssetService(
		assetRepository,
		fixedIDGenerator{id: "asset-2"},
	)

	_, err := firstService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Radio",
	})
	if err != nil {
		t.Fatalf("expected no error creating first asset, got %v", err)
	}

	_, err = secondService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Battery",
	})
	if err != nil {
		t.Fatalf("expected no error creating second asset, got %v", err)
	}

	assets, err := assetRepository.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error listing assets, got %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 asset records, got %d", len(assets))
	}
}
