package web

import (
	"strings"
	"testing"

	"cordell/internal/domain"
)

func TestHumanizeCustodyCorrectionErrorExplainsInsufficientBalance(t *testing.T) {
	message := humanizeCustodyCorrectionError(domain.ErrInsufficientCustodyBalance)

	if message == "" {
		t.Fatal("expected message")
	}

	if message == "This edit cannot be applied because it would make a custody balance negative." {
		t.Fatal("expected expanded insufficient balance explanation")
	}
}

func TestNewCorrectionOptionsDoNotMarkInactiveRecords(t *testing.T) {
	personnel := mustBuildCorrectionTestPersonnel(t, "personnel-1", false)
	asset, err := domain.ReconstituteAsset("asset-1", "Radio", false)
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	personnelOptions := newCorrectionPersonnelOptions([]domain.Personnel{personnel}, string(personnel.ID()))
	if strings.Contains(personnelOptions[0].Label, "(Inactive)") {
		t.Fatalf("expected personnel option without inactive label, got %q", personnelOptions[0].Label)
	}

	assetOptions := newCustodyAssetOptions([]domain.Asset{asset})
	if strings.Contains(assetOptions[0].Label, "(Inactive)") {
		t.Fatalf("expected asset option without inactive label, got %q", assetOptions[0].Label)
	}
}

func TestCorrectionSelectedPersonnelIDClearsInactiveEffectivePersonnel(t *testing.T) {
	selectedID := correctionSelectedPersonnelID(
		"personnel-1",
		map[string]struct{}{"personnel-2": {}},
		"",
	)

	if selectedID != "" {
		t.Fatalf("expected blank selection for inactive effective personnel, got %q", selectedID)
	}
}

func TestCorrectionSelectedPersonnelIDKeepsActiveEffectivePersonnel(t *testing.T) {
	selectedID := correctionSelectedPersonnelID(
		"personnel-1",
		map[string]struct{}{"personnel-1": {}},
		"",
	)

	if selectedID != "personnel-1" {
		t.Fatalf("expected active effective personnel to remain selected, got %q", selectedID)
	}
}

func TestCorrectionSelectedPersonnelIDKeepsStateSelection(t *testing.T) {
	selectedID := correctionSelectedPersonnelID(
		"personnel-1",
		map[string]struct{}{"personnel-1": {}},
		"personnel-2",
	)

	if selectedID != "personnel-2" {
		t.Fatalf("expected state selection to be preserved, got %q", selectedID)
	}
}

func TestCorrectionLineRowsForFormClearsInactiveEffectiveAsset(t *testing.T) {
	rows := correctionLineRowsForForm(
		[]custodyLineFormRowView{
			{
				AssetID:           "asset-1",
				Quantity:          "1",
				CurrentAssetLabel: "Radio",
			},
		},
		map[string]struct{}{"asset-2": {}},
		nil,
	)

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	if rows[0].AssetID != "" {
		t.Fatalf("expected blank asset selection, got %q", rows[0].AssetID)
	}

	if !rows[0].NeedsReplacement {
		t.Fatal("expected inactive effective asset to need replacement")
	}

	if rows[0].CurrentAssetLabel != "Radio" {
		t.Fatalf("expected current asset label to be preserved, got %q", rows[0].CurrentAssetLabel)
	}
}

func TestCorrectionLineRowsForFormKeepsActiveEffectiveAsset(t *testing.T) {
	rows := correctionLineRowsForForm(
		[]custodyLineFormRowView{
			{
				AssetID:           "asset-1",
				Quantity:          "1",
				CurrentAssetLabel: "Radio",
			},
		},
		map[string]struct{}{"asset-1": {}},
		nil,
	)

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	if rows[0].AssetID != "asset-1" {
		t.Fatalf("expected active effective asset to remain selected, got %q", rows[0].AssetID)
	}

	if rows[0].NeedsReplacement {
		t.Fatal("expected active effective asset not to need replacement")
	}
}

func mustBuildCorrectionTestPersonnel(
	t *testing.T,
	id domain.PersonnelID,
	active bool,
) domain.Personnel {
	t.Helper()

	registrationID, err := domain.NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	personnel, err := domain.ReconstitutePersonnel(
		id,
		"John Doe",
		"Doe",
		domain.PersonnelRankSergeant,
		registrationID,
		domain.PersonnelSectionOperations,
		domain.OrganizationUnitDefault,
		active,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	return personnel
}
