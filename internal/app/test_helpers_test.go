package app

import (
	"context"
	"fmt"
	"testing"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) NewID() (string, error) {
	return g.id, nil
}

type fakePersonnelRepository struct {
	saved []domain.Personnel
	byID  map[domain.PersonnelID]domain.Personnel
}

func (r *fakePersonnelRepository) Save(_ context.Context, personnel domain.Personnel) error {
	r.saved = append(r.saved, personnel)

	if r.byID == nil {
		r.byID = make(map[domain.PersonnelID]domain.Personnel)
	}

	r.byID[personnel.ID()] = personnel

	return nil
}

func (r *fakePersonnelRepository) FindByID(_ context.Context, id domain.PersonnelID) (domain.Personnel, error) {
	personnel, ok := r.byID[id]
	if !ok {
		return domain.Personnel{}, ports.ErrNotFound
	}

	return personnel, nil
}

type fakeAssetRepository struct {
	saved []domain.Asset
	byID  map[domain.AssetID]domain.Asset
}

func (r *fakeAssetRepository) Save(_ context.Context, asset domain.Asset) error {
	r.saved = append(r.saved, asset)

	if r.byID == nil {
		r.byID = make(map[domain.AssetID]domain.Asset)
	}

	r.byID[asset.ID()] = asset

	return nil
}

func (r *fakeAssetRepository) FindByID(_ context.Context, id domain.AssetID) (domain.Asset, error) {
	asset, ok := r.byID[id]
	if !ok {
		return domain.Asset{}, ports.ErrNotFound
	}

	return asset, nil
}

type fakeCustodyRepository struct {
	saved           []domain.CustodyTransaction
	currentQuantity map[string]int
}

func (r *fakeCustodyRepository) SaveTransaction(_ context.Context, transaction domain.CustodyTransaction) error {
	r.saved = append(r.saved, transaction)

	return nil
}

func (r *fakeCustodyRepository) CurrentQuantity(
	_ context.Context,
	personnelID domain.PersonnelID,
	assetID domain.AssetID,
) (int, error) {
	if r.currentQuantity == nil {
		return 0, nil
	}

	return r.currentQuantity[custodyBalanceKey(personnelID, assetID)], nil
}

func custodyBalanceKey(personnelID domain.PersonnelID, assetID domain.AssetID) string {
	return fmt.Sprintf("%s:%s", personnelID, assetID)
}

func mustBuildPersonnel(t *testing.T, id domain.PersonnelID) domain.Personnel {
	t.Helper()

	personnel, err := domain.NewPersonnel(id, "John Doe")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	return personnel
}

func mustBuildAsset(t *testing.T, id domain.AssetID) domain.Asset {
	t.Helper()

	asset, err := domain.NewAsset(id, "Radio")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	return asset
}
