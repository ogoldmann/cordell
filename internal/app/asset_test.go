package app

import (
	"context"
	"testing"

	"cordell/internal/domain"
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
