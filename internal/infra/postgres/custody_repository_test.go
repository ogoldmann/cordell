package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
)

func TestPostgresCustodyRepositoryRegisterCheckout(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"52998224725",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	createPersonnelService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	createAssetService := app.NewCreateAssetService(
		assetRepository,
		fixedIDGenerator{id: "asset-1"},
	)
	registerCheckoutService := app.NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-1"},
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	_, err := createPersonnelService.Execute(context.Background(), validCreatePersonnelCommand("John Doe", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error creating personnel, got %v", err)
	}

	_, err = createAssetService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Radio",
	})
	if err != nil {
		t.Fatalf("expected no error creating asset, got %v", err)
	}

	transaction, err := registerCheckoutService.Execute(context.Background(), app.RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 2,
			},
		},
		Notes: "  Operational checkout  ",
	})
	if err != nil {
		t.Fatalf("expected no error registering checkout, got %v", err)
	}

	if transaction.Type() != domain.CustodyTransactionTypeCheckout {
		t.Fatalf("expected checkout transaction, got %s", transaction.Type())
	}

	currentQuantity, err := custodyRepository.CurrentQuantity(
		context.Background(),
		"personnel-1",
		"asset-1",
	)
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if currentQuantity != 2 {
		t.Fatalf("expected current quantity 2, got %d", currentQuantity)
	}

	currentItems, err := custodyRepository.ListCurrentByPersonnel(
		context.Background(),
		"personnel-1",
	)
	if err != nil {
		t.Fatalf("expected no error listing current custody, got %v", err)
	}

	if len(currentItems) != 1 {
		t.Fatalf("expected 1 current custody item, got %d", len(currentItems))
	}

	if currentItems[0].AssetName != "Radio" {
		t.Fatalf("expected asset name Radio, got %s", currentItems[0].AssetName)
	}

	if currentItems[0].Quantity != 2 {
		t.Fatalf("expected quantity 2, got %d", currentItems[0].Quantity)
	}
}

func TestPostgresCustodyRepositoryRegisterReturn(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"52998224725",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	createPersonnelService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	createAssetService := app.NewCreateAssetService(
		assetRepository,
		fixedIDGenerator{id: "asset-1"},
	)
	registerCheckoutService := app.NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-1"},
	)
	registerReturnService := app.NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-return-1"},
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	_, err := createPersonnelService.Execute(context.Background(), validCreatePersonnelCommand("John Doe", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error creating personnel, got %v", err)
	}

	_, err = createAssetService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Radio",
	})
	if err != nil {
		t.Fatalf("expected no error creating asset, got %v", err)
	}

	_, err = registerCheckoutService.Execute(context.Background(), app.RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 2,
			},
		},
		Notes: "Initial checkout",
	})
	if err != nil {
		t.Fatalf("expected no error registering checkout, got %v", err)
	}

	transaction, err := registerReturnService.Execute(context.Background(), app.RegisterReturnCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 1,
			},
		},
		Notes: "  Operational return  ",
	})
	if err != nil {
		t.Fatalf("expected no error registering return, got %v", err)
	}

	if transaction.Type() != domain.CustodyTransactionTypeReturn {
		t.Fatalf("expected return transaction, got %s", transaction.Type())
	}

	currentQuantity, err := custodyRepository.CurrentQuantity(
		context.Background(),
		"personnel-1",
		"asset-1",
	)
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if currentQuantity != 1 {
		t.Fatalf("expected current quantity 1, got %d", currentQuantity)
	}
}

func TestPostgresCustodyRepositoryRejectsInsufficientBalance(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"52998224725",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	createPersonnelService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	createAssetService := app.NewCreateAssetService(
		assetRepository,
		fixedIDGenerator{id: "asset-1"},
	)
	registerCheckoutService := app.NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-1"},
	)
	registerReturnService := app.NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-return-1"},
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	_, err := createPersonnelService.Execute(context.Background(), validCreatePersonnelCommand("John Doe", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error creating personnel, got %v", err)
	}

	_, err = createAssetService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Radio",
	})
	if err != nil {
		t.Fatalf("expected no error creating asset, got %v", err)
	}

	_, err = registerCheckoutService.Execute(context.Background(), app.RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 1,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error registering checkout, got %v", err)
	}

	_, err = registerReturnService.Execute(context.Background(), app.RegisterReturnCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 2,
			},
		},
	})
	if err != domain.ErrInsufficientCustodyBalance {
		t.Fatalf("expected ErrInsufficientCustodyBalance, got %v", err)
	}

	currentQuantity, err := custodyRepository.CurrentQuantity(
		context.Background(),
		"personnel-1",
		"asset-1",
	)
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if currentQuantity != 1 {
		t.Fatalf("expected current quantity to remain 1, got %d", currentQuantity)
	}
}

func TestPostgresCustodyRepositoryListHistoryByPersonnel(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"52998224725",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	createPersonnelService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	createAssetService := app.NewCreateAssetService(
		assetRepository,
		fixedIDGenerator{id: "asset-1"},
	)
	registerCheckoutService := app.NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-1"},
	)
	registerReturnService := app.NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-return-1"},
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	_, err := createPersonnelService.Execute(context.Background(), validCreatePersonnelCommand("John Doe", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error creating personnel, got %v", err)
	}

	_, err = createAssetService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Radio",
	})
	if err != nil {
		t.Fatalf("expected no error creating asset, got %v", err)
	}

	_, err = registerCheckoutService.Execute(context.Background(), app.RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 2,
			},
		},
		Notes: "Initial checkout",
	})
	if err != nil {
		t.Fatalf("expected no error registering checkout, got %v", err)
	}

	_, err = registerReturnService.Execute(context.Background(), app.RegisterReturnCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 1,
			},
		},
		Notes: "Partial return",
	})
	if err != nil {
		t.Fatalf("expected no error registering return, got %v", err)
	}

	history, err := custodyRepository.ListHistoryByPersonnel(context.Background(), "personnel-1", 10)
	if err != nil {
		t.Fatalf("expected no error listing custody history, got %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 history entries, got %d", len(history))
	}

	if history[0].Type != domain.CustodyTransactionTypeReturn {
		t.Fatalf("expected most recent transaction to be return, got %s", history[0].Type)
	}

	if history[0].OperatorID != operator.ID() {
		t.Fatalf("expected operator id %s, got %s", operator.ID(), history[0].OperatorID)
	}

	if history[0].OperatorAlias != operator.Alias() {
		t.Fatalf("expected operator alias %s, got %s", operator.Alias(), history[0].OperatorAlias)
	}

	if history[0].OperatorRank != operator.Rank() {
		t.Fatalf("expected operator rank %s, got %s", operator.Rank(), history[0].OperatorRank)
	}

	if history[1].Type != domain.CustodyTransactionTypeCheckout {
		t.Fatalf("expected oldest transaction to be checkout, got %s", history[1].Type)
	}

	if len(history[0].Lines) != 1 {
		t.Fatalf("expected 1 return line, got %d", len(history[0].Lines))
	}

	if history[0].Lines[0].AssetName != "Radio" {
		t.Fatalf("expected asset name Radio, got %s", history[0].Lines[0].AssetName)
	}

	if history[0].Lines[0].Quantity != 1 {
		t.Fatalf("expected return quantity 1, got %d", history[0].Lines[0].Quantity)
	}
}

func TestPostgresCustodyRepositoryListCurrentByAsset(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"52998224725",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	createFirstPersonnelService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	createSecondPersonnelService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-2"},
	)
	createAssetService := app.NewCreateAssetService(
		assetRepository,
		fixedIDGenerator{id: "asset-1"},
	)

	_, err := createFirstPersonnelService.Execute(context.Background(), validCreatePersonnelCommand("John Doe", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error creating first personnel, got %v", err)
	}

	_, err = createSecondPersonnelService.Execute(context.Background(), validCreatePersonnelCommand("Jane Doe", "Jane", "11144477735"))
	if err != nil {
		t.Fatalf("expected no error creating second personnel, got %v", err)
	}

	_, err = createAssetService.Execute(context.Background(), app.CreateAssetCommand{
		Name: "Radio",
	})
	if err != nil {
		t.Fatalf("expected no error creating asset, got %v", err)
	}

	firstCheckoutService := app.NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-1"},
	)
	secondCheckoutService := app.NewRegisterCheckoutService(
		personnelRepository,
		assetRepository,
		operatorRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-2"},
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	_, err = firstCheckoutService.Execute(context.Background(), app.RegisterCheckoutCommand{
		PersonnelID: "personnel-1",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 2,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error registering first checkout, got %v", err)
	}

	_, err = secondCheckoutService.Execute(context.Background(), app.RegisterCheckoutCommand{
		PersonnelID: "personnel-2",
		OperatorID:  operator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  "asset-1",
				Quantity: 1,
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error registering second checkout, got %v", err)
	}

	holders, err := custodyRepository.ListCurrentByAsset(context.Background(), "asset-1")
	if err != nil {
		t.Fatalf("expected no error listing current asset holders, got %v", err)
	}

	if len(holders) != 2 {
		t.Fatalf("expected 2 current asset holders, got %d", len(holders))
	}

	if holders[0].PersonnelFullName == "" {
		t.Fatal("expected first holder personnel full name to be present")
	}

	if holders[0].Quantity <= 0 {
		t.Fatalf("expected positive holder quantity, got %d", holders[0].Quantity)
	}
}

func TestPostgresCustodyRepositoryFindReceiptByID(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	personnel := mustNewTestPersonnel(t, "personnel-1", "João Silva", "silva", domain.RankSergeant, "52998224725")
	if err := personnelRepository.Save(context.Background(), personnel); err != nil {
		t.Fatalf("expected no error saving personnel, got %v", err)
	}

	asset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if err := assetRepository.Save(context.Background(), asset); err != nil {
		t.Fatalf("expected no error saving asset, got %v", err)
	}

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"93541134780",
		"costa",
		domain.RankCorporal,
		domain.OperatorRoleAdmin,
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid custody line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"Issued for field activity.",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	if err := custodyRepository.SaveTransaction(context.Background(), transaction); err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	receipt, err := custodyRepository.FindReceiptByID(context.Background(), transaction.ID())
	if err != nil {
		t.Fatalf("expected no error finding receipt, got %v", err)
	}

	if receipt.ID != transaction.ID() {
		t.Fatalf("expected receipt id %s, got %s", transaction.ID(), receipt.ID)
	}

	if receipt.PersonnelID != personnel.ID() {
		t.Fatalf("expected personnel id %s, got %s", personnel.ID(), receipt.PersonnelID)
	}

	if receipt.OperatorID != operator.ID() {
		t.Fatalf("expected operator id %s, got %s", operator.ID(), receipt.OperatorID)
	}

	if receipt.OperatorAlias != operator.Alias() {
		t.Fatalf("expected operator alias %s, got %s", operator.Alias(), receipt.OperatorAlias)
	}

	if len(receipt.Lines) != 1 {
		t.Fatalf("expected 1 receipt line, got %d", len(receipt.Lines))
	}

	if receipt.Lines[0].AssetName != asset.Name() {
		t.Fatalf("expected asset name %s, got %s", asset.Name(), receipt.Lines[0].AssetName)
	}
}

func TestPostgresCustodyRepositoryCurrentStateSequence(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	personnel := mustNewTestPersonnel(t, "personnel-1", "John Doe", "doe", domain.RankSergeant, "52998224725")
	if err := personnelRepository.Save(context.Background(), personnel); err != nil {
		t.Fatalf("expected no error saving personnel, got %v", err)
	}

	asset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if err := assetRepository.Save(context.Background(), asset); err != nil {
		t.Fatalf("expected no error saving asset, got %v", err)
	}

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"93541134780",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	checkoutLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(3))
	if err != nil {
		t.Fatalf("expected valid checkout line, got %v", err)
	}

	checkout, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{checkoutLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid checkout, got %v", err)
	}

	if err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
		t.Fatalf("expected no error saving checkout, got %v", err)
	}

	quantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 3 {
		t.Fatalf("expected current quantity 3 after checkout, got %d", quantity)
	}

	returnLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid return line, got %v", err)
	}

	firstReturn, err := domain.NewCustodyTransaction(
		"transaction-2",
		domain.CustodyTransactionTypeReturn,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{returnLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid return, got %v", err)
	}

	if err := custodyRepository.SaveTransaction(context.Background(), firstReturn); err != nil {
		t.Fatalf("expected no error saving first return, got %v", err)
	}

	quantity, err = custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 2 {
		t.Fatalf("expected current quantity 2 after first return, got %d", quantity)
	}

	finalReturnLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(2))
	if err != nil {
		t.Fatalf("expected valid final return line, got %v", err)
	}

	finalReturn, err := domain.NewCustodyTransaction(
		"transaction-3",
		domain.CustodyTransactionTypeReturn,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{finalReturnLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid final return, got %v", err)
	}

	if err := custodyRepository.SaveTransaction(context.Background(), finalReturn); err != nil {
		t.Fatalf("expected no error saving final return, got %v", err)
	}

	quantity, err = custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 0 {
		t.Fatalf("expected current quantity 0 after final return, got %d", quantity)
	}

	currentItems, err := custodyRepository.ListCurrentByPersonnel(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error listing current custody, got %v", err)
	}

	if len(currentItems) != 0 {
		t.Fatalf("expected no current custody items after full return, got %d", len(currentItems))
	}
}

func TestPostgresCustodyRepositoryRollsBackInvalidReturn(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	personnel := mustNewTestPersonnel(t, "personnel-1", "John Doe", "doe", domain.RankSergeant, "52998224725")
	if err := personnelRepository.Save(context.Background(), personnel); err != nil {
		t.Fatalf("expected no error saving personnel, got %v", err)
	}

	asset, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid asset, got %v", err)
	}

	if err := assetRepository.Save(context.Background(), asset); err != nil {
		t.Fatalf("expected no error saving asset, got %v", err)
	}

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"93541134780",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	checkoutLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid checkout line, got %v", err)
	}

	checkout, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{checkoutLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid checkout, got %v", err)
	}

	if err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
		t.Fatalf("expected no error saving checkout, got %v", err)
	}

	invalidReturnLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(2))
	if err != nil {
		t.Fatalf("expected valid invalid-return line, got %v", err)
	}

	invalidReturn, err := domain.NewCustodyTransaction(
		"transaction-2",
		domain.CustodyTransactionTypeReturn,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{invalidReturnLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid return transaction object, got %v", err)
	}

	err = custodyRepository.SaveTransaction(context.Background(), invalidReturn)
	if err != domain.ErrInsufficientCustodyBalance {
		t.Fatalf("expected ErrInsufficientCustodyBalance, got %v", err)
	}

	quantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 1 {
		t.Fatalf("expected current quantity to remain 1 after rolled back invalid return, got %d", quantity)
	}

	history, err := custodyRepository.ListHistoryByPersonnel(context.Background(), personnel.ID(), 10)
	if err != nil {
		t.Fatalf("expected no error listing history, got %v", err)
	}

	if len(history) != 1 {
		t.Fatalf("expected only original checkout in history after rollback, got %d entries", len(history))
	}

	if history[0].ID != checkout.ID() {
		t.Fatalf("expected checkout transaction in history, got %s", history[0].ID)
	}
}

func TestPostgresCustodyRepositoryHandlesMultipleAssetsInOneTransaction(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	personnel := mustNewTestPersonnel(t, "personnel-1", "John Doe", "doe", domain.RankSergeant, "52998224725")
	if err := personnelRepository.Save(context.Background(), personnel); err != nil {
		t.Fatalf("expected no error saving personnel, got %v", err)
	}

	radio, err := domain.NewAsset("asset-1", "Radio")
	if err != nil {
		t.Fatalf("expected valid radio, got %v", err)
	}

	helmet, err := domain.NewAsset("asset-2", "Helmet")
	if err != nil {
		t.Fatalf("expected valid helmet, got %v", err)
	}

	if err := assetRepository.Save(context.Background(), radio); err != nil {
		t.Fatalf("expected no error saving radio, got %v", err)
	}

	if err := assetRepository.Save(context.Background(), helmet); err != nil {
		t.Fatalf("expected no error saving helmet, got %v", err)
	}

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"93541134780",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleAdmin,
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	radioLine, err := domain.NewCustodyLine(radio.ID(), domain.Quantity(2))
	if err != nil {
		t.Fatalf("expected valid radio line, got %v", err)
	}

	helmetLine, err := domain.NewCustodyLine(helmet.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid helmet line, got %v", err)
	}

	checkout, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{radioLine, helmetLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid checkout, got %v", err)
	}

	if err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
		t.Fatalf("expected no error saving checkout, got %v", err)
	}

	radioQuantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), radio.ID())
	if err != nil {
		t.Fatalf("expected no error reading radio quantity, got %v", err)
	}

	if radioQuantity != 2 {
		t.Fatalf("expected radio quantity 2, got %d", radioQuantity)
	}

	helmetQuantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), helmet.ID())
	if err != nil {
		t.Fatalf("expected no error reading helmet quantity, got %v", err)
	}

	if helmetQuantity != 1 {
		t.Fatalf("expected helmet quantity 1, got %d", helmetQuantity)
	}

	currentItems, err := custodyRepository.ListCurrentByPersonnel(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error listing current custody, got %v", err)
	}

	if len(currentItems) != 2 {
		t.Fatalf("expected 2 current custody items, got %d", len(currentItems))
	}
}
