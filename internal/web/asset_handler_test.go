package web

import (
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestNewCustodyTimelineItemFromAssetHistoryItemHighlightsCurrentAsset(t *testing.T) {
	item := ports.AssetCustodyHistoryItem{
		ID:                "transaction-1",
		Type:              domain.CustodyTransactionTypeCheckout,
		CreatedAt:         time.Date(2026, 7, 24, 10, 30, 0, 0, time.Local),
		PersonnelID:       "personnel-1",
		PersonnelRank:     domain.PersonnelRankSoldier,
		PersonnelAlias:    "John",
		PersonnelFullName: "John Doe",
		OperatorID:        "operator-1",
		OperatorRank:      domain.RankThirdSergeant,
		OperatorAlias:     "Silva",
		Lines: []ports.AssetCustodyHistoryLine{
			{
				AssetID:   "asset-1",
				AssetName: "Radio",
				Quantity:  1,
			},
			{
				AssetID:   "asset-2",
				AssetName: "Helmet",
				Quantity:  1,
			},
		},
	}

	view := newCustodyTimelineItemFromAssetHistoryItem(item, "asset-1")

	if len(view.Lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(view.Lines))
	}

	if !view.Lines[0].Highlighted {
		t.Fatal("expected current asset line to be highlighted")
	}

	if view.Lines[1].Highlighted {
		t.Fatal("expected other asset line not to be highlighted")
	}
}
