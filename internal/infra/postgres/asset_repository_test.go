package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
	"cordell/internal/ports"
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

	assets, err := assetRepository.List(context.Background(), 10, ports.RecordStatusFilterActive)
	if err != nil {
		t.Fatalf("expected no error listing assets, got %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 asset records, got %d", len(assets))
	}
}

func TestPostgresAssetRepositoryListFiltersByStatus(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)
	repository := postgres.NewAssetRepository(queries)

	activeAsset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid active asset, got %v", err)
	}

	inactiveAsset, err := domain.ReconstituteAsset("asset-2", "Battery", false)
	if err != nil {
		t.Fatalf("expected valid inactive asset, got %v", err)
	}

	if err := repository.Save(context.Background(), activeAsset); err != nil {
		t.Fatalf("expected no error saving active asset, got %v", err)
	}

	if err := repository.Save(context.Background(), inactiveAsset); err != nil {
		t.Fatalf("expected no error saving inactive asset, got %v", err)
	}

	active, err := repository.List(context.Background(), 10, ports.RecordStatusFilterActive)
	if err != nil {
		t.Fatalf("expected no error listing active assets, got %v", err)
	}

	if len(active) != 1 {
		t.Fatalf("expected 1 active asset, got %d", len(active))
	}

	if !active[0].Active() {
		t.Fatal("expected listed asset to be active")
	}

	inactive, err := repository.List(context.Background(), 10, ports.RecordStatusFilterInactive)
	if err != nil {
		t.Fatalf("expected no error listing inactive assets, got %v", err)
	}

	if len(inactive) != 1 {
		t.Fatalf("expected 1 inactive asset, got %d", len(inactive))
	}

	if inactive[0].Active() {
		t.Fatal("expected listed asset to be inactive")
	}

	all, err := repository.List(context.Background(), 10, ports.RecordStatusFilterAll)
	if err != nil {
		t.Fatalf("expected no error listing all assets, got %v", err)
	}

	if len(all) != 2 {
		t.Fatalf("expected 2 assets, got %d", len(all))
	}
}

func TestPostgresAssetRepositorySearch(t *testing.T) {
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

	assets, err := assetRepository.Search(context.Background(), "battery", 10)
	if err != nil {
		t.Fatalf("expected no error searching assets, got %v", err)
	}

	if len(assets) != 1 {
		t.Fatalf("expected 1 asset record, got %d", len(assets))
	}

	if assets[0].Name() != "Battery" {
		t.Fatalf("expected asset name Battery, got %s", assets[0].Name())
	}
}

func TestPostgresAssetRepositorySearchByCombinedTerms(t *testing.T) {
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
		Name: "Radio Battery",
	})
	if err != nil {
		t.Fatalf("expected no error creating first asset, got %v", err)
	}

	_, err = secondService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Radio Cable",
	})
	if err != nil {
		t.Fatalf("expected no error creating second asset, got %v", err)
	}

	assets, err := assetRepository.Search(context.Background(), "radio battery", 10)
	if err != nil {
		t.Fatalf("expected no error searching assets, got %v", err)
	}

	if len(assets) != 1 {
		t.Fatalf("expected 1 asset record, got %d", len(assets))
	}

	if assets[0].ID() != "asset-1" {
		t.Fatalf("expected asset-1, got %s", assets[0].ID())
	}
}
