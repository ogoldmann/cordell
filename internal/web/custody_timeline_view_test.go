package web

import "testing"

func TestCustodyTimelineTypeTone(t *testing.T) {
	if got := custodyTimelineTypeTone(checkoutLabel()); got != "checkout" {
		t.Fatalf("expected checkout tone, got %q", got)
	}

	if got := custodyTimelineTypeTone(returnLabel()); got != "return" {
		t.Fatalf("expected return tone, got %q", got)
	}

	if got := custodyTimelineTypeTone("Outro"); got != "neutral" {
		t.Fatalf("expected neutral tone, got %q", got)
	}
}

func TestCustodyEditCountLabel(t *testing.T) {
	if got := custodyEditCountLabel(0); got != "" {
		t.Fatalf("expected empty label, got %q", got)
	}

	if got := custodyEditCountLabel(1); got != "1 edição" {
		t.Fatalf("expected 1 edição, got %q", got)
	}

	if got := custodyEditCountLabel(2); got != "2 edições" {
		t.Fatalf("expected 2 edições, got %q", got)
	}
}

func TestFormatTimelineQuantity(t *testing.T) {
	if got := formatTimelineQuantity(3); got != "3" {
		t.Fatalf("expected 3, got %q", got)
	}
}

func TestNewCustodyTimelineItemFromLedgerItem(t *testing.T) {
	item := newCustodyTimelineItemFromLedgerItem(custodyTransactionSummaryView{
		ID:                        "transaction-1",
		SequenceLabel:             "#7",
		DateLabel:                 "24/07/2026",
		TimeLabel:                 "10:30",
		ReceiptURL:                "/custody/transactions/transaction-1",
		TypeLabel:                 checkoutLabel(),
		EffectivePersonnelURL:     "/personnel/personnel-1",
		EffectivePersonnelDisplay: "Sd Silva",
		OperatorDisplay:           "Sgt Costa",
		HasCorrection:             true,
		EditCountLabel:            "2 edições",
		Lines: []custodyTransactionSummaryLineView{
			{
				AssetID:   "asset-1",
				AssetURL:  "/assets/asset-1",
				AssetName: "Radio",
				Quantity:  3,
			},
		},
	})

	if item.TypeTone != "checkout" {
		t.Fatalf("expected checkout tone, got %q", item.TypeTone)
	}

	if item.TypeLabel != "CAUTELA" {
		t.Fatalf("expected uppercase checkout label, got %q", item.TypeLabel)
	}

	if item.ReceiptURL != "/custody/transactions/transaction-1" {
		t.Fatalf("expected receipt URL, got %q", item.ReceiptURL)
	}

	if item.RegisteredBy != "Sgt Costa" {
		t.Fatalf("expected registered by label, got %q", item.RegisteredBy)
	}

	if item.RegisteredAt != "24/07/2026 10:30" {
		t.Fatalf("expected registered at label, got %q", item.RegisteredAt)
	}

	if !item.Edited {
		t.Fatal("expected edited item")
	}

	if item.SecondaryActionURL != "/custody/transactions/transaction-1/edit" {
		t.Fatalf("expected edit action URL, got %q", item.SecondaryActionURL)
	}

	if len(item.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(item.Lines))
	}

	if item.Lines[0].Quantity != "3" {
		t.Fatalf("expected quantity 3, got %q", item.Lines[0].Quantity)
	}
}

func TestNewCustodyTimelineItemFromPersonnelHistoryItem(t *testing.T) {
	item := newCustodyTimelineItemFromPersonnelHistoryItem(custodyHistoryView{
		ID:              "transaction-1",
		TypeLabel:       returnLabel(),
		DateLabel:       "24/07/2026",
		TimeLabel:       "10:30",
		PersonnelLabel:  "Sd Silva",
		PersonnelURL:    "/personnel/personnel-1",
		OperatorDisplay: "Sgt Costa",
		Notes:           "Returned after inspection.",
		Lines: []custodyHistoryLineView{
			{
				AssetID:   "asset-1",
				AssetName: "Radio",
				Quantity:  1,
			},
		},
	})

	if item.TypeTone != "return" {
		t.Fatalf("expected return tone, got %q", item.TypeTone)
	}

	if item.TypeLabel != "DESCAUTELA" {
		t.Fatalf("expected uppercase return label, got %q", item.TypeLabel)
	}

	if item.Notes != "Returned after inspection." {
		t.Fatalf("expected notes, got %q", item.Notes)
	}

	if item.PersonnelURL != "/personnel/personnel-1" {
		t.Fatalf("expected personnel URL, got %q", item.PersonnelURL)
	}

	if item.Lines[0].AssetURL != "/assets/asset-1" {
		t.Fatalf("expected asset URL, got %q", item.Lines[0].AssetURL)
	}
}
