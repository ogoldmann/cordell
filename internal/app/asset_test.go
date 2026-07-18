package app

import (
	"context"
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
