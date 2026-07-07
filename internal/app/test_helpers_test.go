package app

import (
	"context"
	"fmt"
	"sort"
	"strings"
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

func (r *fakePersonnelRepository) List(_ context.Context, limit int) ([]domain.Personnel, error) {
	personnel := make([]domain.Personnel, 0, len(r.byID))

	for _, item := range r.byID {
		personnel = append(personnel, item)
	}

	sort.Slice(personnel, func(i, j int) bool {
		return personnel[i].ID() < personnel[j].ID()
	})

	if limit > 0 && len(personnel) > limit {
		return personnel[:limit], nil
	}

	return personnel, nil
}

func (r *fakePersonnelRepository) Search(_ context.Context, query string, limit int) ([]domain.Personnel, error) {
	query = strings.ToLower(strings.TrimSpace(query))

	personnel := make([]domain.Personnel, 0, len(r.byID))

	for _, item := range r.byID {
		if personnelMatchesQuery(item, query) {
			personnel = append(personnel, item)
		}
	}

	sort.Slice(personnel, func(i, j int) bool {
		return personnel[i].ID() < personnel[j].ID()
	})

	if limit > 0 && len(personnel) > limit {
		return personnel[:limit], nil
	}

	return personnel, nil
}

func personnelMatchesQuery(personnel domain.Personnel, query string) bool {
	if query == "" {
		return true
	}

	values := []string{
		personnel.FullName(),
		personnel.Alias(),
		personnel.RegistrationID().String(),
		string(personnel.Rank()),
		string(personnel.Section()),
		string(personnel.OrganizationUnit()),
	}

	for _, value := range values {
		if strings.Contains(strings.ToLower(value), query) {
			return true
		}
	}

	return false
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

func (r *fakeAssetRepository) List(_ context.Context, limit int) ([]domain.Asset, error) {
	assets := make([]domain.Asset, 0, len(r.byID))

	for _, item := range r.byID {
		assets = append(assets, item)
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].ID() < assets[j].ID()
	})

	if limit > 0 && len(assets) > limit {
		return assets[:limit], nil
	}

	return assets, nil
}

func (r *fakeAssetRepository) Search(_ context.Context, query string, limit int) ([]domain.Asset, error) {
	query = strings.ToLower(strings.TrimSpace(query))

	assets := make([]domain.Asset, 0, len(r.byID))

	for _, item := range r.byID {
		if query == "" || strings.Contains(strings.ToLower(item.Name()), query) {
			assets = append(assets, item)
		}
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].ID() < assets[j].ID()
	})

	if limit > 0 && len(assets) > limit {
		return assets[:limit], nil
	}

	return assets, nil
}

type fakeCustodyRepository struct {
	saved           []domain.CustodyTransaction
	currentQuantity map[string]int
	currentByPerson map[domain.PersonnelID][]ports.CurrentCustodyItem
	currentByAsset  map[domain.AssetID][]ports.CurrentAssetHolder
	historyByPerson map[domain.PersonnelID][]ports.CustodyHistoryEntry
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

func (r *fakeCustodyRepository) ListCurrentByPersonnel(
	_ context.Context,
	personnelID domain.PersonnelID,
) ([]ports.CurrentCustodyItem, error) {
	if r.currentByPerson == nil {
		return nil, nil
	}

	items := r.currentByPerson[personnelID]
	copiedItems := make([]ports.CurrentCustodyItem, len(items))
	copy(copiedItems, items)

	return copiedItems, nil
}

func (r *fakeCustodyRepository) ListCurrentByAsset(
	_ context.Context,
	assetID domain.AssetID,
) ([]ports.CurrentAssetHolder, error) {
	if r.currentByAsset == nil {
		return nil, nil
	}

	holders := r.currentByAsset[assetID]
	copiedHolders := make([]ports.CurrentAssetHolder, len(holders))
	copy(copiedHolders, holders)

	return copiedHolders, nil
}

func (r *fakeCustodyRepository) ListHistoryByPersonnel(
	_ context.Context,
	personnelID domain.PersonnelID,
	limit int,
) ([]ports.CustodyHistoryEntry, error) {
	if r.historyByPerson == nil {
		return nil, nil
	}

	items := r.historyByPerson[personnelID]
	copiedItems := make([]ports.CustodyHistoryEntry, len(items))
	copy(copiedItems, items)

	if limit > 0 && len(copiedItems) > limit {
		return copiedItems[:limit], nil
	}

	return copiedItems, nil
}

func custodyBalanceKey(personnelID domain.PersonnelID, assetID domain.AssetID) string {
	return fmt.Sprintf("%s:%s", personnelID, assetID)
}

func mustBuildPersonnel(t *testing.T, id string) domain.Personnel {
	t.Helper()

	registrationID, err := domain.NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	personnel, err := domain.NewPersonnel(
		domain.PersonnelID(id),
		"John Doe",
		"Doe",
		domain.PersonnelRankSergeant,
		registrationID,
		domain.PersonnelSectionOperations,
		domain.OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	return personnel
}

func validCreatePersonnelCommand(fullName string, alias string, registrationID string) CreatePersonnelCommand {
	return CreatePersonnelCommand{
		FullName:         fullName,
		Alias:            alias,
		Rank:             domain.PersonnelRankSergeant,
		RegistrationID:   registrationID,
		Section:          domain.PersonnelSectionOperations,
		OrganizationUnit: domain.OrganizationUnitDefault,
	}
}

func mustBuildAsset(t *testing.T, id domain.AssetID) domain.Asset {
	t.Helper()

	asset, err := domain.NewAsset(id, "Radio")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	return asset
}
