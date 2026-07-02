package app

import (
	"context"
	"testing"

	"cordell/internal/domain"
)

func TestRegisterCheckoutServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	asset := mustBuildAsset(t, "asset-1")

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}
	custodyRepository := &fakeCustodyRepository{}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		custodyRepository,
		idGenerator,
	)

	transaction, err := service.Execute(context.Background(), RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		Lines: []CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 2,
			},
		},
		Notes: "  Operational checkout  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transaction.ID() != "transaction-1" {
		t.Fatalf("expected transaction id transaction-1, got %s", transaction.ID())
	}

	if transaction.Type() != domain.CustodyTransactionTypeCheckout {
		t.Fatalf("expected checkout transaction, got %s", transaction.Type())
	}

	if transaction.PersonnelID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", transaction.PersonnelID())
	}

	if len(transaction.Lines()) != 1 {
		t.Fatalf("expected 1 transaction line, got %d", len(transaction.Lines()))
	}

	if transaction.Notes() != "Operational checkout" {
		t.Fatalf("expected trimmed notes, got %s", transaction.Notes())
	}

	if len(custodyRepository.saved) != 1 {
		t.Fatalf("expected 1 saved transaction, got %d", len(custodyRepository.saved))
	}
}

func TestRegisterCheckoutServiceRejectsInvalidQuantity(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	asset := mustBuildAsset(t, "asset-1")

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}
	custodyRepository := &fakeCustodyRepository{}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		custodyRepository,
		idGenerator,
	)

	_, err := service.Execute(context.Background(), RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		Lines: []CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 0,
			},
		},
	})
	if err != domain.ErrInvalidQuantity {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}

	if len(custodyRepository.saved) != 0 {
		t.Fatalf("expected no saved transactions, got %d", len(custodyRepository.saved))
	}
}

func TestRegisterReturnServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	asset := mustBuildAsset(t, "asset-1")

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}
	custodyRepository := &fakeCustodyRepository{
		currentQuantity: map[string]int{
			custodyBalanceKey("personnel-1", "asset-1"): 2,
		},
	}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		custodyRepository,
		idGenerator,
	)

	transaction, err := service.Execute(context.Background(), RegisterReturnCommand{
		PersonnelID: "personnel-1",
		Lines: []CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 1,
			},
		},
		Notes: "  Operational return  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transaction.Type() != domain.CustodyTransactionTypeReturn {
		t.Fatalf("expected return transaction, got %s", transaction.Type())
	}

	if transaction.Notes() != "Operational return" {
		t.Fatalf("expected trimmed notes, got %s", transaction.Notes())
	}

	if len(custodyRepository.saved) != 1 {
		t.Fatalf("expected 1 saved transaction, got %d", len(custodyRepository.saved))
	}
}

func TestRegisterReturnServiceRejectsInsufficientCustodyBalance(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	asset := mustBuildAsset(t, "asset-1")

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}
	custodyRepository := &fakeCustodyRepository{
		currentQuantity: map[string]int{
			custodyBalanceKey("personnel-1", "asset-1"): 1,
		},
	}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		custodyRepository,
		idGenerator,
	)

	_, err := service.Execute(context.Background(), RegisterReturnCommand{
		PersonnelID: "personnel-1",
		Lines: []CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 2,
			},
		},
	})
	if err != domain.ErrInsufficientCustodyBalance {
		t.Fatalf("expected ErrInsufficientCustodyBalance, got %v", err)
	}

	if len(custodyRepository.saved) != 0 {
		t.Fatalf("expected no saved transactions, got %d", len(custodyRepository.saved))
	}
}
