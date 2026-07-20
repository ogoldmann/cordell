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

	assetOptions := newCorrectionAssetOptions([]domain.Asset{asset})
	if strings.Contains(assetOptions[0].Label, "(Inactive)") {
		t.Fatalf("expected asset option without inactive label, got %q", assetOptions[0].Label)
	}
}

func TestCustodyCorrectionEditBlockedReasonRequiresActivePersonnel(t *testing.T) {
	reason := custodyCorrectionEditBlockedReason(
		"personnel-1",
		[]correctionLineRowView{{AssetID: "asset-1", Quantity: "1"}},
		map[string]struct{}{"personnel-2": {}},
		map[string]struct{}{"asset-1": {}},
	)

	expected := "This transaction currently references an inactive personnel. Reactivate the personnel before editing this transaction."
	if reason != expected {
		t.Fatalf("expected inactive personnel reason, got %q", reason)
	}
}

func TestCustodyCorrectionEditBlockedReasonRequiresActiveAssets(t *testing.T) {
	reason := custodyCorrectionEditBlockedReason(
		"personnel-1",
		[]correctionLineRowView{{AssetID: "asset-1", Quantity: "1"}},
		map[string]struct{}{"personnel-1": {}},
		map[string]struct{}{"asset-2": {}},
	)

	expected := "This transaction currently references an inactive asset. Reactivate the asset before editing this transaction."
	if reason != expected {
		t.Fatalf("expected inactive asset reason, got %q", reason)
	}
}

func TestCustodyCorrectionEditBlockedReasonAllowsActiveRecords(t *testing.T) {
	reason := custodyCorrectionEditBlockedReason(
		"personnel-1",
		[]correctionLineRowView{{AssetID: "asset-1", Quantity: "1"}},
		map[string]struct{}{"personnel-1": {}},
		map[string]struct{}{"asset-1": {}},
	)

	if reason != "" {
		t.Fatalf("expected no blocked reason, got %q", reason)
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
