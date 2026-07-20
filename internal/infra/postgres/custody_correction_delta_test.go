package postgres

import (
	"testing"

	"cordell/internal/domain"
)

func TestCustodyCorrectionBalanceDeltasForCheckoutPersonnelChange(t *testing.T) {
	previousLine := mustBuildPostgresCustodyLine(t, "asset-1", 1)
	correctedLine := mustBuildPostgresCustodyLine(t, "asset-1", 1)

	deltas := custodyCorrectionBalanceDeltas(
		domain.CustodyTransactionTypeCheckout,
		"personnel-a",
		[]domain.CustodyLine{previousLine},
		"personnel-b",
		[]domain.CustodyLine{correctedLine},
	)

	assertCustodyBalanceDeltas(t, deltas, []custodyBalanceDelta{
		{personnelID: "personnel-a", assetID: "asset-1", quantity: -1},
		{personnelID: "personnel-b", assetID: "asset-1", quantity: 1},
	})
}

func TestCustodyCorrectionBalanceDeltasForReturnQuantityChange(t *testing.T) {
	previousLine := mustBuildPostgresCustodyLine(t, "asset-1", 3)
	correctedLine := mustBuildPostgresCustodyLine(t, "asset-1", 1)

	deltas := custodyCorrectionBalanceDeltas(
		domain.CustodyTransactionTypeReturn,
		"personnel-a",
		[]domain.CustodyLine{previousLine},
		"personnel-a",
		[]domain.CustodyLine{correctedLine},
	)

	assertCustodyBalanceDeltas(t, deltas, []custodyBalanceDelta{
		{personnelID: "personnel-a", assetID: "asset-1", quantity: 2},
	})
}

func assertCustodyBalanceDeltas(t *testing.T, got []custodyBalanceDelta, want []custodyBalanceDelta) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("expected %d deltas, got %d: %#v", len(want), len(got), got)
	}

	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("expected delta %#v at index %d, got %#v", want[index], index, got[index])
		}
	}
}

func mustBuildPostgresCustodyLine(t *testing.T, assetID domain.AssetID, quantityValue int) domain.CustodyLine {
	t.Helper()

	quantity, err := domain.NewQuantity(quantityValue)
	if err != nil {
		t.Fatalf("expected valid quantity, got %v", err)
	}

	line, err := domain.NewCustodyLine(assetID, quantity)
	if err != nil {
		t.Fatalf("expected valid custody line, got %v", err)
	}

	return line
}
