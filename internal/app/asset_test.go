package app

import (
	"context"
	"errors"
	"testing"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestCreateAssetServiceExecute(t *testing.T) {
	repository := &fakeAssetRepository{}
	idGenerator := fixedIDGenerator{id: "asset-1"}

	service := NewCreateAssetService(repository, idGenerator)

	asset, err := service.Execute(context.Background(), CreateAssetCommand{
		Name: "  Radio  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if asset.ID() != "asset-1" {
		t.Fatalf("expected asset id asset-1, got %s", asset.ID())
	}

	if asset.Name() != "Radio" {
		t.Fatalf("expected trimmed asset name Radio, got %s", asset.Name())
	}

	if len(repository.saved) != 1 {
		t.Fatalf("expected 1 saved asset, got %d", len(repository.saved))
	}
}

func TestCreateAssetServiceRejectsInvalidAsset(t *testing.T) {
	repository := &fakeAssetRepository{}
	idGenerator := fixedIDGenerator{id: "asset-1"}

	service := NewCreateAssetService(repository, idGenerator)

	_, err := service.Execute(context.Background(), CreateAssetCommand{
		Name: "   ",
	})
	if err != domain.ErrEmptyAssetName {
		t.Fatalf("expected ErrEmptyAssetName, got %v", err)
	}

	if len(repository.saved) != 0 {
		t.Fatalf("expected no saved asset, got %d", len(repository.saved))
	}
}

func TestCreateAssetServiceReturnsDuplicateAssetName(t *testing.T) {
	repository := &fakeAssetRepository{
		saveErr: domain.ErrDuplicateAssetName,
	}
	service := NewCreateAssetService(repository, fixedIDGenerator{id: "asset-1"})

	_, err := service.Execute(context.Background(), CreateAssetCommand{
		Name: "Radio",
	})
	if err != domain.ErrDuplicateAssetName {
		t.Fatalf("expected ErrDuplicateAssetName, got %v", err)
	}
}

func TestGetAssetServiceExecute(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	service := NewGetAssetService(repository)

	found, err := service.Execute(context.Background(), GetAssetCommand{
		ID: "asset-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID() != "asset-1" {
		t.Fatalf("expected asset id asset-1, got %s", found.ID())
	}

	if found.Name() != "Radio" {
		t.Fatalf("expected asset name Radio, got %s", found.Name())
	}
}

func TestGetAssetServiceRejectsEmptyID(t *testing.T) {
	repository := &fakeAssetRepository{}
	service := NewGetAssetService(repository)

	_, err := service.Execute(context.Background(), GetAssetCommand{
		ID: "",
	})
	if err != domain.ErrEmptyAssetID {
		t.Fatalf("expected ErrEmptyAssetID, got %v", err)
	}
}

func TestUpdateAssetServiceUpdatesAsset(t *testing.T) {
	asset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	service := NewUpdateAssetService(repository)

	updated, err := service.Execute(context.Background(), UpdateAssetCommand{
		ID:   asset.ID(),
		Name: "Updated Radio",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updated.Name() != "Updated Radio" {
		t.Fatalf("expected updated name, got %q", updated.Name())
	}
}

func TestUpdateAssetServiceReturnsNotFound(t *testing.T) {
	service := NewUpdateAssetService(&fakeAssetRepository{})

	_, err := service.Execute(context.Background(), UpdateAssetCommand{
		ID:   "missing",
		Name: "Radio",
	})
	if !errors.Is(err, ports.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestUpdateAssetServiceRejectsDuplicateName(t *testing.T) {
	first, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid first asset, got %v", err)
	}

	second, err := domain.NewAsset("asset-2", "Helmet")
	if err != nil {
		t.Fatalf("expected valid second asset, got %v", err)
	}

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			first.ID():  first,
			second.ID(): second,
		},
	}

	service := NewUpdateAssetService(repository)

	_, err = service.Execute(context.Background(), UpdateAssetCommand{
		ID:   second.ID(),
		Name: first.Name(),
	})
	if !errors.Is(err, domain.ErrDuplicateAssetName) {
		t.Fatalf("expected duplicate asset name, got %v", err)
	}
}

func TestListAssetsServiceExecute(t *testing.T) {
	firstAsset := mustBuildAsset(t, "asset-1")
	secondAsset := mustBuildAsset(t, "asset-2")

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			firstAsset.ID():  firstAsset,
			secondAsset.ID(): secondAsset,
		},
	}

	service := NewListAssetsService(repository)

	assets, err := service.Execute(context.Background(), ListAssetsCommand{
		Limit:        10,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(assets) != 2 {
		t.Fatalf("expected 2 asset records, got %d", len(assets))
	}
}

func TestListAssetsServiceAppliesDefaultLimit(t *testing.T) {
	repository := &fakeAssetRepository{}
	service := NewListAssetsService(repository)

	_, err := service.Execute(context.Background(), ListAssetsCommand{
		Limit:        0,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestListAssetsServiceCapsLimit(t *testing.T) {
	repository := &fakeAssetRepository{}
	service := NewListAssetsService(repository)

	_, err := service.Execute(context.Background(), ListAssetsCommand{
		Limit:        1000,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestListAssetsServiceDefaultsToActiveStatusFilter(t *testing.T) {
	repository := &fakeAssetRepository{}
	service := NewListAssetsService(repository)

	_, err := service.Execute(context.Background(), ListAssetsCommand{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.lastStatusFilter != ports.RecordStatusFilterActive {
		t.Fatalf("expected active status filter, got %s", repository.lastStatusFilter)
	}
}

func TestDeactivateAssetServiceExecute(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	service := NewDeactivateAssetService(repository)

	err := service.Execute(context.Background(), DeactivateAssetCommand{
		ID: asset.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := repository.FindByID(context.Background(), asset.ID())
	if err != nil {
		t.Fatalf("expected asset to exist, got %v", err)
	}

	if updated.Active() {
		t.Fatal("expected asset to be inactive")
	}
}

func TestDeactivateAssetServiceIsNoOpForInactiveAsset(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")
	inactiveAsset, err := domain.ReconstituteAsset(asset.ID(), asset.Name(), false)
	if err != nil {
		t.Fatalf("expected valid inactive asset, got %v", err)
	}

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			inactiveAsset.ID(): inactiveAsset,
		},
	}

	service := NewDeactivateAssetService(repository)

	err = service.Execute(context.Background(), DeactivateAssetCommand{
		ID: inactiveAsset.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReactivateAssetServiceExecute(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")
	inactiveAsset, err := domain.ReconstituteAsset(asset.ID(), asset.Name(), false)
	if err != nil {
		t.Fatalf("expected valid inactive asset, got %v", err)
	}

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			inactiveAsset.ID(): inactiveAsset,
		},
	}

	service := NewReactivateAssetService(repository)

	err = service.Execute(context.Background(), ReactivateAssetCommand{
		ID: inactiveAsset.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := repository.FindByID(context.Background(), inactiveAsset.ID())
	if err != nil {
		t.Fatalf("expected asset to exist, got %v", err)
	}

	if !updated.Active() {
		t.Fatal("expected asset to be active")
	}
}

func TestReactivateAssetServiceIsNoOpForActiveAsset(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	service := NewReactivateAssetService(repository)

	err := service.Execute(context.Background(), ReactivateAssetCommand{
		ID: asset.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSearchAssetsServiceExecute(t *testing.T) {
	firstAsset := mustBuildAsset(t, "asset-1")

	secondAsset, err := domain.NewAsset("asset-2", "Battery")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			firstAsset.ID():  firstAsset,
			secondAsset.ID(): secondAsset,
		},
	}

	service := NewSearchAssetsService(repository)

	assets, err := service.Execute(context.Background(), SearchAssetsCommand{
		Query: "battery",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(assets) != 1 {
		t.Fatalf("expected 1 asset record, got %d", len(assets))
	}

	if assets[0].Name() != "Battery" {
		t.Fatalf("expected asset name Battery, got %s", assets[0].Name())
	}
}

func TestSearchAssetsServiceFallsBackToListWhenQueryIsEmpty(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")

	repository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	service := NewSearchAssetsService(repository)

	result, err := service.Execute(context.Background(), SearchAssetsCommand{
		Query: "   ",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 asset record, got %d", len(result))
	}
}
