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

func TestPostgresAssetRepositoryRejectsDuplicateName(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)
	repository := postgres.NewAssetRepository(queries)

	first, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	second, err := domain.NewAsset("asset-2", "radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if err := repository.Save(context.Background(), first); err != nil {
		t.Fatalf("expected no error saving first asset, got %v", err)
	}

	err = repository.Save(context.Background(), second)
	if err != domain.ErrDuplicateAssetName {
		t.Fatalf("expected ErrDuplicateAssetName, got %v", err)
	}
}

func TestPostgresAssetRepositoryRejectsDuplicateNormalizedName(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)
	repository := postgres.NewAssetRepository(queries)

	first, err := domain.NewAsset("asset-1", "Radio VHF")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	second, err := domain.NewAsset("asset-2", "  Radio   VHF  ")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if err := repository.Save(context.Background(), first); err != nil {
		t.Fatalf("expected no error saving first asset, got %v", err)
	}

	err = repository.Save(context.Background(), second)
	if err != domain.ErrDuplicateAssetName {
		t.Fatalf("expected ErrDuplicateAssetName, got %v", err)
	}
}

func TestPostgresAssetRepositoryDeactivateAndReactivate(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)
	repository := postgres.NewAssetRepository(queries)

	asset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if err := repository.Save(context.Background(), asset); err != nil {
		t.Fatalf("expected no error saving asset, got %v", err)
	}

	deactivated, err := repository.Deactivate(context.Background(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating asset, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected asset to be deactivated")
	}

	updated, err := repository.FindByID(context.Background(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error finding asset, got %v", err)
	}

	if updated.Active() {
		t.Fatal("expected asset to be inactive")
	}

	reactivated, err := repository.Reactivate(context.Background(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reactivating asset, got %v", err)
	}

	if !reactivated {
		t.Fatal("expected asset to be reactivated")
	}

	updated, err = repository.FindByID(context.Background(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error finding asset, got %v", err)
	}

	if !updated.Active() {
		t.Fatal("expected asset to be active")
	}
}

func TestPostgresAssetRepositoryDuplicateNameStillRejectedAfterDeactivation(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)
	repository := postgres.NewAssetRepository(queries)

	first, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	second, err := domain.NewAsset("asset-2", "radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if err := repository.Save(context.Background(), first); err != nil {
		t.Fatalf("expected no error saving first asset, got %v", err)
	}

	deactivated, err := repository.Deactivate(context.Background(), first.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating first asset, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected first asset to be deactivated")
	}

	err = repository.Save(context.Background(), second)
	if err != domain.ErrDuplicateAssetName {
		t.Fatalf("expected ErrDuplicateAssetName, got %v", err)
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

	assets, err := assetRepository.Search(context.Background(), "battery", 10, ports.RecordStatusFilterActive)
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

	assets, err := assetRepository.Search(context.Background(), "radio battery", 10, ports.RecordStatusFilterActive)
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

func TestPostgresAssetRepositorySearchFiltersByStatus(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	assetRepository := postgres.NewAssetRepository(queries)

	activeAsset, err := domain.NewAsset("asset-1", "Radio Active")
	if err != nil {
		t.Fatalf("expected valid active asset, got %v", err)
	}

	inactiveBase, err := domain.NewAsset("asset-2", "Radio Inactive")
	if err != nil {
		t.Fatalf("expected valid inactive base asset, got %v", err)
	}

	inactiveAsset, err := domain.ReconstituteAsset(inactiveBase.ID(), inactiveBase.Name(), false)
	if err != nil {
		t.Fatalf("expected valid inactive asset, got %v", err)
	}

	if err := assetRepository.Save(context.Background(), activeAsset); err != nil {
		t.Fatalf("expected no error saving active asset, got %v", err)
	}

	if err := assetRepository.Save(context.Background(), inactiveAsset); err != nil {
		t.Fatalf("expected no error saving inactive asset, got %v", err)
	}

	activeResults, err := assetRepository.Search(
		context.Background(),
		"Radio",
		10,
		ports.RecordStatusFilterActive,
	)
	if err != nil {
		t.Fatalf("expected no error searching active assets, got %v", err)
	}

	if len(activeResults) != 1 {
		t.Fatalf("expected 1 active result, got %d", len(activeResults))
	}

	if !activeResults[0].Active() {
		t.Fatal("expected active result")
	}

	inactiveResults, err := assetRepository.Search(
		context.Background(),
		"Radio",
		10,
		ports.RecordStatusFilterInactive,
	)
	if err != nil {
		t.Fatalf("expected no error searching inactive assets, got %v", err)
	}

	if len(inactiveResults) != 1 {
		t.Fatalf("expected 1 inactive result, got %d", len(inactiveResults))
	}

	if inactiveResults[0].Active() {
		t.Fatal("expected inactive result")
	}

	allResults, err := assetRepository.Search(
		context.Background(),
		"Radio",
		10,
		ports.RecordStatusFilterAll,
	)
	if err != nil {
		t.Fatalf("expected no error searching all assets, got %v", err)
	}

	if len(allResults) != 2 {
		t.Fatalf("expected 2 all results, got %d", len(allResults))
	}
}
