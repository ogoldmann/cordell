package web

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestHandleGlobalSearchRendersPartialResults(t *testing.T) {
	personnel := mustBuildSearchTestPersonnel(t)
	asset := mustBuildSearchTestAsset(t)

	server, err := NewServer(
		slog.Default(),
		app.Services{
			GlobalSearch: app.NewGlobalSearchService(
				fakeGlobalSearchPersonnelRepository{personnel: []domain.Personnel{personnel}},
				fakeGlobalSearchAssetRepository{assets: []domain.Asset{asset}},
			),
		},
		NewSessionCookieConfig(false),
		NewSecurityHeadersConfig(false),
	)
	if err != nil {
		t.Fatalf("expected server, got %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/search?q=radio", nil)
	request.Header.Set("X-Cordell-Partial", "1")
	response := httptest.NewRecorder()

	server.handleGlobalSearch(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	body := response.Body.String()
	if !strings.Contains(body, "Materiais") {
		t.Fatalf("expected asset result section, got body %q", body)
	}

	if strings.Contains(body, "Pesquisa global") {
		t.Fatalf("expected partial response without full search page header, got body %q", body)
	}
}

type fakeGlobalSearchPersonnelRepository struct {
	personnel []domain.Personnel
}

func (r fakeGlobalSearchPersonnelRepository) Save(context.Context, domain.Personnel) error {
	return nil
}

func (r fakeGlobalSearchPersonnelRepository) Update(context.Context, domain.Personnel) error {
	return nil
}

func (r fakeGlobalSearchPersonnelRepository) FindByID(
	context.Context,
	domain.PersonnelID,
) (domain.Personnel, error) {
	return domain.Personnel{}, ports.ErrNotFound
}

func (r fakeGlobalSearchPersonnelRepository) FindByRegistrationID(
	context.Context,
	domain.RegistrationID,
) (domain.Personnel, bool, error) {
	return domain.Personnel{}, false, nil
}

func (r fakeGlobalSearchPersonnelRepository) FindByRegistrationIDExcludingID(
	context.Context,
	domain.RegistrationID,
	domain.PersonnelID,
) (domain.Personnel, bool, error) {
	return domain.Personnel{}, false, nil
}

func (r fakeGlobalSearchPersonnelRepository) List(
	context.Context,
	int,
	ports.RecordStatusFilter,
) ([]domain.Personnel, error) {
	return r.personnel, nil
}

func (r fakeGlobalSearchPersonnelRepository) Search(
	context.Context,
	string,
	int,
	ports.RecordStatusFilter,
) ([]domain.Personnel, error) {
	return r.personnel, nil
}

func (r fakeGlobalSearchPersonnelRepository) Deactivate(context.Context, domain.PersonnelID) (bool, error) {
	return false, nil
}

func (r fakeGlobalSearchPersonnelRepository) Reactivate(context.Context, domain.PersonnelID) (bool, error) {
	return false, nil
}

type fakeGlobalSearchAssetRepository struct {
	assets []domain.Asset
}

func (r fakeGlobalSearchAssetRepository) Save(context.Context, domain.Asset) error {
	return nil
}

func (r fakeGlobalSearchAssetRepository) Update(context.Context, domain.Asset) error {
	return nil
}

func (r fakeGlobalSearchAssetRepository) FindByID(context.Context, domain.AssetID) (domain.Asset, error) {
	return domain.Asset{}, ports.ErrNotFound
}

func (r fakeGlobalSearchAssetRepository) FindByName(context.Context, string) (domain.Asset, bool, error) {
	return domain.Asset{}, false, nil
}

func (r fakeGlobalSearchAssetRepository) FindByNameExcludingID(
	context.Context,
	string,
	domain.AssetID,
) (domain.Asset, bool, error) {
	return domain.Asset{}, false, nil
}

func (r fakeGlobalSearchAssetRepository) List(
	context.Context,
	int,
	ports.RecordStatusFilter,
) ([]domain.Asset, error) {
	return r.assets, nil
}

func (r fakeGlobalSearchAssetRepository) Search(
	context.Context,
	string,
	int,
	ports.RecordStatusFilter,
) ([]domain.Asset, error) {
	return r.assets, nil
}

func (r fakeGlobalSearchAssetRepository) Deactivate(context.Context, domain.AssetID) (bool, error) {
	return false, nil
}

func (r fakeGlobalSearchAssetRepository) Reactivate(context.Context, domain.AssetID) (bool, error) {
	return false, nil
}

func mustBuildSearchTestPersonnel(t *testing.T) domain.Personnel {
	t.Helper()

	registrationID, err := domain.NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	personnel, err := domain.NewPersonnel(
		"personnel-1",
		"John Radio",
		"Radio",
		domain.RankSergeant,
		registrationID,
		domain.PersonnelSectionOperations,
		domain.OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	return personnel
}

func mustBuildSearchTestAsset(t *testing.T) domain.Asset {
	t.Helper()

	asset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	return asset
}
