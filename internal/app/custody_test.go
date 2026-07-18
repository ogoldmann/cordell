package app

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestRegisterCheckoutServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	asset := mustBuildAsset(t, "asset-1")
	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

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
	operatorRepository := newFakeOperatorRepository(operator)
	custodyRepository := &fakeCustodyRepository{}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		idGenerator,
	)

	transaction, err := service.Execute(context.Background(), RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
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

	if transaction.OperatorID() != operator.ID() {
		t.Fatalf("expected operator id %s, got %s", operator.ID(), transaction.OperatorID())
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
	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

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
	operatorRepository := newFakeOperatorRepository(operator)
	custodyRepository := &fakeCustodyRepository{}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		idGenerator,
	)

	_, err := service.Execute(context.Background(), RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
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

func TestRegisterCheckoutServiceRejectsEmptyOperatorID(t *testing.T) {
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
	operatorRepository := newFakeOperatorRepository()
	custodyRepository := &fakeCustodyRepository{}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		idGenerator,
	)

	_, err := service.Execute(context.Background(), RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  "",
		Lines: []CustodyLineCommand{
			{AssetID: "asset-1", Quantity: 1},
		},
	})

	if err != domain.ErrEmptyOperatorID {
		t.Fatalf("expected ErrEmptyOperatorID, got %v", err)
	}
}

func TestRegisterReturnServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	asset := mustBuildAsset(t, "asset-1")
	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

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
	operatorRepository := newFakeOperatorRepository(operator)
	custodyRepository := &fakeCustodyRepository{
		currentQuantity: map[string]int{
			custodyBalanceKey("personnel-1", "asset-1"): 2,
		},
	}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		idGenerator,
	)

	transaction, err := service.Execute(context.Background(), RegisterReturnCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
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
	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

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
	operatorRepository := newFakeOperatorRepository(operator)
	custodyRepository := &fakeCustodyRepository{
		currentQuantity: map[string]int{
			custodyBalanceKey("personnel-1", "asset-1"): 1,
		},
	}
	idGenerator := fixedIDGenerator{id: "transaction-1"}

	service := NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		idGenerator,
	)

	_, err := service.Execute(context.Background(), RegisterReturnCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
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

func TestListCurrentCustodyServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	custodyRepository := &fakeCustodyRepository{
		currentByPerson: map[domain.PersonnelID][]ports.CurrentCustodyItem{
			"personnel-1": {
				{
					PersonnelID: "personnel-1",
					AssetID:     "asset-1",
					AssetName:   "Radio",
					Quantity:    2,
				},
			},
		},
	}

	service := NewListCurrentCustodyService(personnelRepository, custodyRepository)

	items, err := service.Execute(context.Background(), ListCurrentCustodyCommand{
		PersonnelID: "personnel-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 current custody item, got %d", len(items))
	}

	if items[0].AssetName != "Radio" {
		t.Fatalf("expected asset name Radio, got %s", items[0].AssetName)
	}

	if items[0].Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", items[0].Quantity)
	}
}

func TestListCustodyHistoryServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	custodyRepository := &fakeCustodyRepository{
		historyByPerson: map[domain.PersonnelID][]ports.CustodyHistoryEntry{
			"personnel-1": {
				{
					ID:          "transaction-1",
					Type:        domain.CustodyTransactionTypeCheckout,
					PersonnelID: "personnel-1",
					Notes:       "Initial checkout",
					CreatedAt:   time.Date(2026, 7, 6, 10, 0, 0, 0, time.UTC),
					Lines: []ports.CustodyHistoryLine{
						{
							AssetID:   "asset-1",
							AssetName: "Radio",
							Quantity:  2,
						},
					},
				},
			},
		},
	}

	service := NewListCustodyHistoryService(personnelRepository, custodyRepository)

	entries, err := service.Execute(context.Background(), ListCustodyHistoryCommand{
		PersonnelID: "personnel-1",
		Limit:       10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(entries))
	}

	if entries[0].Type != domain.CustodyTransactionTypeCheckout {
		t.Fatalf("expected checkout transaction, got %s", entries[0].Type)
	}

	if len(entries[0].Lines) != 1 {
		t.Fatalf("expected 1 history line, got %d", len(entries[0].Lines))
	}

	if entries[0].Lines[0].AssetName != "Radio" {
		t.Fatalf("expected asset name Radio, got %s", entries[0].Lines[0].AssetName)
	}

	if entries[0].Lines[0].Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", entries[0].Lines[0].Quantity)
	}
}

func TestListCurrentAssetHoldersServiceExecute(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")

	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}
	custodyRepository := &fakeCustodyRepository{
		currentByAsset: map[domain.AssetID][]ports.CurrentAssetHolder{
			"asset-1": {
				{
					AssetID:           "asset-1",
					PersonnelID:       "personnel-1",
					PersonnelFullName: "John Doe",
					Quantity:          2,
				},
			},
		},
	}

	service := NewListCurrentAssetHoldersService(assetRepository, custodyRepository)

	holders, err := service.Execute(context.Background(), ListCurrentAssetHoldersCommand{
		AssetID: "asset-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(holders) != 1 {
		t.Fatalf("expected 1 current asset holder, got %d", len(holders))
	}

	if holders[0].PersonnelFullName != "John Doe" {
		t.Fatalf("expected personnel full name John Doe, got %s", holders[0].PersonnelFullName)
	}

	if holders[0].Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", holders[0].Quantity)
	}
}
