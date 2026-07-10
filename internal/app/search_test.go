package app

import (
	"context"
	"testing"

	"cordell/internal/domain"
)

func TestGlobalSearchServiceExecuteSearchesPersonnelAndAssets(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	asset, err := domain.NewAsset("asset-1", "Doe Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	service := NewGlobalSearchService(personnelRepository, assetRepository)

	result, err := service.Execute(context.Background(), GlobalSearchCommand{
		Query:        "doe",
		LimitPerType: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Personnel) != 1 {
		t.Fatalf("expected 1 personnel result, got %d", len(result.Personnel))
	}

	if len(result.Assets) != 1 {
		t.Fatalf("expected 1 asset result, got %d", len(result.Assets))
	}

	if result.Query != "doe" {
		t.Fatalf("expected query doe, got %s", result.Query)
	}
}

func TestGlobalSearchServiceExecuteReturnsEmptyResultForEmptyQuery(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	asset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	service := NewGlobalSearchService(personnelRepository, assetRepository)

	result, err := service.Execute(context.Background(), GlobalSearchCommand{
		Query:        "   ",
		LimitPerType: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Personnel) != 0 {
		t.Fatalf("expected 0 personnel results, got %d", len(result.Personnel))
	}

	if len(result.Assets) != 0 {
		t.Fatalf("expected 0 asset results, got %d", len(result.Assets))
	}

	if result.Query != "" {
		t.Fatalf("expected empty query, got %s", result.Query)
	}
}
