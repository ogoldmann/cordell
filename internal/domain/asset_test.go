package domain

import (
	"errors"
	"testing"
)

func TestNewAsset(t *testing.T) {
	asset, err := NewAsset("asset-1", "  Radio  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if asset.ID() != "asset-1" {
		t.Fatalf("expected asset id asset-1, got %s", asset.ID())
	}

	if asset.Name() != "Radio" {
		t.Fatalf("expected trimmed asset name Radio, got %s", asset.Name())
	}

	if !asset.Active() {
		t.Fatal("expected asset to be active")
	}
}

func TestNewAssetRejectsEmptyID(t *testing.T) {
	_, err := NewAsset("", "Radio")
	if err != ErrEmptyAssetID {
		t.Fatalf("expected ErrEmptyAssetID, got %v", err)
	}
}

func TestNewAssetRejectsEmptyName(t *testing.T) {
	_, err := NewAsset("asset-1", "   ")
	if err != ErrEmptyAssetName {
		t.Fatalf("expected ErrEmptyAssetName, got %v", err)
	}
}

func TestNewAssetNormalizesName(t *testing.T) {
	asset, err := NewAsset("asset-1", "  Radio   VHF  ")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if asset.Name() != "Radio VHF" {
		t.Fatalf("expected normalized name Radio VHF, got %q", asset.Name())
	}
}

func TestAssetUpdateDetails(t *testing.T) {
	asset, err := NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	err = asset.UpdateDetails("Updated Radio")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if asset.Name() != "Updated Radio" {
		t.Fatalf("expected updated name, got %q", asset.Name())
	}
}

func TestAssetUpdateDetailsRejectsEmptyName(t *testing.T) {
	asset, err := NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	err = asset.UpdateDetails("   ")
	if !errors.Is(err, ErrEmptyAssetName) {
		t.Fatalf("expected ErrEmptyAssetName, got %v", err)
	}
}
