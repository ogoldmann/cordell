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

func TestRegisterCheckoutServiceCombinesDuplicateAssetLines(t *testing.T) {
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
		PersonnelID: personnel.ID(),
		OperatorID:  operator.ID(),
		Lines: []CustodyLineCommand{
			{AssetID: asset.ID(), Quantity: 1},
			{AssetID: asset.ID(), Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(transaction.Lines()) != 1 {
		t.Fatalf("expected duplicate asset lines to be combined into 1 line, got %d", len(transaction.Lines()))
	}

	if transaction.Lines()[0].Quantity().Int() != 3 {
		t.Fatalf("expected combined quantity 3, got %d", transaction.Lines()[0].Quantity().Int())
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

func TestRegisterReturnServiceCombinesDuplicateAssetLines(t *testing.T) {
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
			custodyBalanceKey(personnel.ID(), asset.ID()): 3,
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
		PersonnelID: personnel.ID(),
		OperatorID:  operator.ID(),
		Lines: []CustodyLineCommand{
			{AssetID: asset.ID(), Quantity: 1},
			{AssetID: asset.ID(), Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(transaction.Lines()) != 1 {
		t.Fatalf("expected duplicate asset lines to be combined into 1 line, got %d", len(transaction.Lines()))
	}

	if transaction.Lines()[0].Quantity().Int() != 3 {
		t.Fatalf("expected combined quantity 3, got %d", transaction.Lines()[0].Quantity().Int())
	}
}

func TestRegisterReturnServiceRejectsInsufficientCombinedDuplicateAssetLines(t *testing.T) {
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
			custodyBalanceKey(personnel.ID(), asset.ID()): 2,
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
		PersonnelID: personnel.ID(),
		OperatorID:  operator.ID(),
		Lines: []CustodyLineCommand{
			{AssetID: asset.ID(), Quantity: 1},
			{AssetID: asset.ID(), Quantity: 2},
		},
	})
	if err != domain.ErrInsufficientCustodyBalance {
		t.Fatalf("expected ErrInsufficientCustodyBalance, got %v", err)
	}

	if len(custodyRepository.saved) != 0 {
		t.Fatalf("expected no saved transactions, got %d", len(custodyRepository.saved))
	}
}

func TestNormalizeCustodyLineCommandsPreservesFirstAssetOrder(t *testing.T) {
	commands := []CustodyLineCommand{
		{AssetID: "asset-1", Quantity: 1},
		{AssetID: "asset-2", Quantity: 5},
		{AssetID: "asset-1", Quantity: 2},
		{AssetID: "asset-3", Quantity: 4},
		{AssetID: "asset-2", Quantity: 1},
	}

	normalized := normalizeCustodyLineCommands(commands)

	if len(normalized) != 3 {
		t.Fatalf("expected 3 normalized lines, got %d", len(normalized))
	}

	expected := []CustodyLineCommand{
		{AssetID: "asset-1", Quantity: 3},
		{AssetID: "asset-2", Quantity: 6},
		{AssetID: "asset-3", Quantity: 4},
	}

	for index, expectedLine := range expected {
		if normalized[index] != expectedLine {
			t.Fatalf("expected normalized line %+v at index %d, got %+v", expectedLine, index, normalized[index])
		}
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
					AssetActive: true,
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

func TestListCurrentCustodyServicePreservesAssetActiveState(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	custodyRepository := &fakeCustodyRepository{
		currentByPerson: map[domain.PersonnelID][]ports.CurrentCustodyItem{
			personnel.ID(): {
				{
					PersonnelID: personnel.ID(),
					AssetID:     "asset-1",
					AssetName:   "Radio",
					AssetActive: false,
					Quantity:    1,
				},
			},
		},
	}

	service := NewListCurrentCustodyService(personnelRepository, custodyRepository)

	items, err := service.Execute(context.Background(), ListCurrentCustodyCommand{
		PersonnelID: personnel.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if items[0].AssetActive {
		t.Fatal("expected inactive asset state to be preserved")
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

func TestGetCustodyReceiptServiceExecute(t *testing.T) {
	createdAt := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	repository := &fakeCustodyRepository{
		receipts: map[domain.CustodyTransactionID]ports.CustodyReceipt{
			"transaction-1": {
				ID:                        "transaction-1",
				Type:                      domain.CustodyTransactionTypeCheckout,
				PersonnelID:               "personnel-1",
				PersonnelFullName:         "João Silva",
				PersonnelAlias:            "silva",
				PersonnelRank:             domain.RankSergeant,
				PersonnelRegistrationID:   domain.RegistrationID("52998224725"),
				PersonnelSection:          domain.PersonnelSectionLogistics,
				PersonnelOrganizationUnit: domain.OrganizationUnitDefault,
				OperatorID:                "operator-1",
				OperatorRegistrationID:    domain.RegistrationID("93541134780"),
				OperatorAlias:             "costa",
				OperatorRank:              domain.RankCorporal,
				Notes:                     "Issued for field activity.",
				CreatedAt:                 createdAt,
				Lines: []ports.CustodyReceiptLine{
					{
						AssetID:   "asset-1",
						AssetName: "Radio",
						Quantity:  1,
					},
				},
			},
		},
	}

	service := NewGetCustodyReceiptService(repository)

	receipt, err := service.Execute(context.Background(), GetCustodyReceiptCommand{
		ID: "transaction-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if receipt.ID != "transaction-1" {
		t.Fatalf("expected transaction-1, got %s", receipt.ID)
	}

	if receipt.PersonnelFullName != "João Silva" {
		t.Fatalf("expected personnel full name João Silva, got %s", receipt.PersonnelFullName)
	}

	if receipt.OperatorAlias != "costa" {
		t.Fatalf("expected operator alias costa, got %s", receipt.OperatorAlias)
	}

	if len(receipt.Lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(receipt.Lines))
	}
}

func TestGetCustodyReceiptServiceRejectsEmptyID(t *testing.T) {
	service := NewGetCustodyReceiptService(&fakeCustodyRepository{})

	_, err := service.Execute(context.Background(), GetCustodyReceiptCommand{})
	if err != domain.ErrEmptyTransactionID {
		t.Fatalf("expected ErrEmptyTransactionID, got %v", err)
	}
}

func TestGetCustodyReceiptServiceReturnsNotFound(t *testing.T) {
	service := NewGetCustodyReceiptService(&fakeCustodyRepository{})

	_, err := service.Execute(context.Background(), GetCustodyReceiptCommand{
		ID: "missing",
	})
	if err != ports.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
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
					PersonnelAlias:    "doe",
					PersonnelRank:     domain.RankSergeant,
					PersonnelActive:   true,
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

func TestListCurrentAssetHoldersServicePreservesPersonnelActiveState(t *testing.T) {
	asset := mustBuildAsset(t, "asset-1")

	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			asset.ID(): asset,
		},
	}

	custodyRepository := &fakeCustodyRepository{
		currentByAsset: map[domain.AssetID][]ports.CurrentAssetHolder{
			asset.ID(): {
				{
					AssetID:           asset.ID(),
					PersonnelID:       "personnel-1",
					PersonnelFullName: "John Doe",
					PersonnelAlias:    "doe",
					PersonnelRank:     domain.RankSergeant,
					PersonnelActive:   false,
					Quantity:          1,
				},
			},
		},
	}

	service := NewListCurrentAssetHoldersService(assetRepository, custodyRepository)

	holders, err := service.Execute(context.Background(), ListCurrentAssetHoldersCommand{
		AssetID: asset.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(holders) != 1 {
		t.Fatalf("expected 1 holder, got %d", len(holders))
	}

	if holders[0].PersonnelActive {
		t.Fatal("expected inactive personnel state to be preserved")
	}
}

func TestRegisterCheckoutServiceRejectsInactivePersonnel(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	inactivePersonnel, err := domain.ReconstitutePersonnel(
		personnel.ID(),
		personnel.FullName(),
		personnel.Alias(),
		personnel.Rank(),
		personnel.RegistrationID(),
		personnel.Section(),
		personnel.OrganizationUnit(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive personnel, got %v", err)
	}

	asset := mustBuildAsset(t, "asset-1")
	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			inactivePersonnel.ID(): inactivePersonnel,
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

	_, err = service.Execute(context.Background(), RegisterCheckoutCommand{
		PersonnelID: inactivePersonnel.ID(),
		OperatorID:  operator.ID(),
		Lines: []CustodyLineCommand{
			{AssetID: asset.ID(), Quantity: 1},
		},
	})
	if err != domain.ErrInactivePersonnel {
		t.Fatalf("expected ErrInactivePersonnel, got %v", err)
	}

	if len(custodyRepository.saved) != 0 {
		t.Fatalf("expected no saved transactions, got %d", len(custodyRepository.saved))
	}
}

func TestRegisterCheckoutServiceRejectsInactiveAsset(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	asset := mustBuildAsset(t, "asset-1")
	inactiveAsset, err := domain.ReconstituteAsset(asset.ID(), asset.Name(), false)
	if err != nil {
		t.Fatalf("expected valid inactive asset, got %v", err)
	}

	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			inactiveAsset.ID(): inactiveAsset,
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

	_, err = service.Execute(context.Background(), RegisterCheckoutCommand{
		PersonnelID: personnel.ID(),
		OperatorID:  operator.ID(),
		Lines: []CustodyLineCommand{
			{AssetID: inactiveAsset.ID(), Quantity: 1},
		},
	})
	if err != domain.ErrInactiveAsset {
		t.Fatalf("expected ErrInactiveAsset, got %v", err)
	}

	if len(custodyRepository.saved) != 0 {
		t.Fatalf("expected no saved transactions, got %d", len(custodyRepository.saved))
	}
}

func TestRegisterReturnServiceAllowsInactiveAssetWhenCurrentlyCustodied(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	asset := mustBuildAsset(t, "asset-1")
	inactiveAsset, err := domain.ReconstituteAsset(asset.ID(), asset.Name(), false)
	if err != nil {
		t.Fatalf("expected valid inactive asset, got %v", err)
	}

	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}
	assetRepository := &fakeAssetRepository{
		byID: map[domain.AssetID]domain.Asset{
			inactiveAsset.ID(): inactiveAsset,
		},
	}
	operatorRepository := newFakeOperatorRepository(operator)
	custodyRepository := &fakeCustodyRepository{
		currentQuantity: map[string]int{
			custodyBalanceKey(personnel.ID(), inactiveAsset.ID()): 1,
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
		PersonnelID: personnel.ID(),
		OperatorID:  operator.ID(),
		Lines: []CustodyLineCommand{
			{AssetID: inactiveAsset.ID(), Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected inactive asset return to be allowed, got %v", err)
	}

	if transaction.ID() != "transaction-1" {
		t.Fatalf("expected transaction-1, got %s", transaction.ID())
	}
}

func TestRegisterReturnServiceAllowsInactivePersonnelWhenReturningCurrentCustody(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	inactivePersonnel, err := domain.ReconstitutePersonnel(
		personnel.ID(),
		personnel.FullName(),
		personnel.Alias(),
		personnel.Rank(),
		personnel.RegistrationID(),
		personnel.Section(),
		personnel.OrganizationUnit(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive personnel, got %v", err)
	}

	asset := mustBuildAsset(t, "asset-1")
	operator := mustNewTestOperator(t, "operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin)

	personnelRepository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			inactivePersonnel.ID(): inactivePersonnel,
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
			custodyBalanceKey(inactivePersonnel.ID(), asset.ID()): 1,
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
		PersonnelID: inactivePersonnel.ID(),
		OperatorID:  operator.ID(),
		Lines: []CustodyLineCommand{
			{AssetID: asset.ID(), Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected inactive personnel return to be allowed, got %v", err)
	}

	if transaction.ID() != "transaction-1" {
		t.Fatalf("expected transaction-1, got %s", transaction.ID())
	}
}
