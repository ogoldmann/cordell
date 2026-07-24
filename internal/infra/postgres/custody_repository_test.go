package postgres_test

import (
	"context"
	"errors"
	"strconv"
	"testing"
	"time"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
	"cordell/internal/ports"
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

	result, err := registerCheckoutService.Execute(context.Background(), app.RegisterCheckoutCommand{
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

	transaction := result.Transaction

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

	result, err := registerReturnService.Execute(context.Background(), app.RegisterReturnCommand{
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

	transaction := result.Transaction

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

func TestPostgresCustodyRepositoryListHistoryByPersonnelUsesEffectiveCorrectedPersonnel(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	personnelA := mustNewTestPersonnel(t, "personnel-1", "John A", "alpha", domain.RankSergeant, "52998224725")
	personnelB := mustNewTestPersonnel(t, "personnel-2", "John B", "bravo", domain.RankCorporal, "11144477735")

	if err := personnelRepository.Save(context.Background(), personnelA); err != nil {
		t.Fatalf("expected no error saving personnel A, got %v", err)
	}

	if err := personnelRepository.Save(context.Background(), personnelB); err != nil {
		t.Fatalf("expected no error saving personnel B, got %v", err)
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
		"29109142088",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleOperator,
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnelA.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	if !created {
		t.Fatal("expected transaction to be created")
	}

	correctionLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid correction line, got %v", err)
	}

	correction, err := domain.NewCustodyCorrection(
		"correction-1",
		transaction.ID(),
		operator.ID(),
		personnelB.ID(),
		[]domain.CustodyLine{correctionLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid correction, got %v", err)
	}

	created, err = custodyRepository.SaveCorrection(
		context.Background(),
		correction,
		transaction.Type(),
		transaction.PersonnelID(),
		transaction.Lines(),
	)
	if err != nil {
		t.Fatalf("expected no error saving correction, got %v", err)
	}

	if !created {
		t.Fatal("expected correction to be created")
	}

	historyA, err := custodyRepository.ListHistoryByPersonnel(context.Background(), personnelA.ID(), 10)
	if err != nil {
		t.Fatalf("expected no error listing personnel A history, got %v", err)
	}

	if len(historyA) != 0 {
		t.Fatalf("expected personnel A history to be empty after correction, got %d", len(historyA))
	}

	historyB, err := custodyRepository.ListHistoryByPersonnel(context.Background(), personnelB.ID(), 10)
	if err != nil {
		t.Fatalf("expected no error listing personnel B history, got %v", err)
	}

	if len(historyB) != 1 {
		t.Fatalf("expected personnel B history to contain corrected transaction, got %d", len(historyB))
	}

	if !historyB[0].HasCorrection {
		t.Fatal("expected corrected history item to be marked as edited")
	}

	if historyB[0].EditCount != 1 {
		t.Fatalf("expected edit count 1, got %d", historyB[0].EditCount)
	}
}

func TestPostgresCustodyRepositoryListAssetCustodyHistoryUsesEffectiveCorrectedLines(t *testing.T) {
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

	operator := mustNewTestOperator(
		t,
		"operator-1",
		"29109142088",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleOperator,
	)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	radio, err := domain.NewAsset("asset-radio", "Radio")
	if err != nil {
		t.Fatalf("expected valid radio asset, got %v", err)
	}
	if err := assetRepository.Save(context.Background(), radio); err != nil {
		t.Fatalf("expected no error saving radio, got %v", err)
	}

	helmet, err := domain.NewAsset("asset-helmet", "Helmet")
	if err != nil {
		t.Fatalf("expected valid helmet asset, got %v", err)
	}
	if err := assetRepository.Save(context.Background(), helmet); err != nil {
		t.Fatalf("expected no error saving helmet, got %v", err)
	}

	radioLine, err := domain.NewCustodyLine(radio.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid radio line, got %v", err)
	}

	helmetLine, err := domain.NewCustodyLine(helmet.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid helmet line, got %v", err)
	}

	removedTransaction, err := domain.NewCustodyTransaction(
		"transaction-removed",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{radioLine, helmetLine},
		"Original radio and helmet.",
	)
	if err != nil {
		t.Fatalf("expected valid removed transaction, got %v", err)
	}
	if _, err := custodyRepository.SaveTransaction(context.Background(), removedTransaction); err != nil {
		t.Fatalf("expected no error saving removed transaction, got %v", err)
	}

	removedCorrection, err := domain.NewCustodyCorrection(
		"correction-removes-radio",
		removedTransaction.ID(),
		operator.ID(),
		personnel.ID(),
		[]domain.CustodyLine{helmetLine},
		"Radio removed from effective state.",
	)
	if err != nil {
		t.Fatalf("expected valid removed correction, got %v", err)
	}
	if _, err := custodyRepository.SaveCorrection(
		context.Background(),
		removedCorrection,
		removedTransaction.Type(),
		removedTransaction.PersonnelID(),
		removedTransaction.Lines(),
	); err != nil {
		t.Fatalf("expected no error saving removed correction, got %v", err)
	}

	addedTransaction, err := domain.NewCustodyTransaction(
		"transaction-added",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{helmetLine},
		"Original helmet only.",
	)
	if err != nil {
		t.Fatalf("expected valid added transaction, got %v", err)
	}
	if _, err := custodyRepository.SaveTransaction(context.Background(), addedTransaction); err != nil {
		t.Fatalf("expected no error saving added transaction, got %v", err)
	}

	addedCorrection, err := domain.NewCustodyCorrection(
		"correction-adds-radio",
		addedTransaction.ID(),
		operator.ID(),
		personnel.ID(),
		[]domain.CustodyLine{helmetLine, radioLine},
		"Radio added to effective state.",
	)
	if err != nil {
		t.Fatalf("expected valid added correction, got %v", err)
	}
	if _, err := custodyRepository.SaveCorrection(
		context.Background(),
		addedCorrection,
		addedTransaction.Type(),
		addedTransaction.PersonnelID(),
		addedTransaction.Lines(),
	); err != nil {
		t.Fatalf("expected no error saving added correction, got %v", err)
	}

	history, err := custodyRepository.ListAssetCustodyHistory(context.Background(), radio.ID())
	if err != nil {
		t.Fatalf("expected no error listing asset history, got %v", err)
	}

	if len(history) != 1 {
		t.Fatalf("expected only corrected transaction containing radio, got %d", len(history))
	}

	if history[0].ID != addedTransaction.ID() {
		t.Fatalf("expected added transaction, got %s", history[0].ID)
	}

	if history[0].EditCount != 1 {
		t.Fatalf("expected edit count 1, got %d", history[0].EditCount)
	}

	if len(history[0].Lines) != 2 {
		t.Fatalf("expected effective transaction lines, got %d", len(history[0].Lines))
	}

	foundRadio := false
	for _, line := range history[0].Lines {
		if line.AssetID == radio.ID() {
			foundRadio = true
		}
	}

	if !foundRadio {
		t.Fatal("expected radio line in effective asset history")
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

func TestPostgresCustodyRepositoryListCurrentByPersonnelIncludesAssetActiveState(t *testing.T) {
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

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	if _, err := custodyRepository.SaveTransaction(context.Background(), transaction); err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	_, err = pool.Exec(context.Background(), `
		UPDATE assets
		SET active = false
		WHERE id = $1
	`, string(asset.ID()))
	if err != nil {
		t.Fatalf("expected no error marking asset inactive, got %v", err)
	}

	items, err := custodyRepository.ListCurrentByPersonnel(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error listing current custody, got %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 current custody item, got %d", len(items))
	}

	if items[0].AssetActive {
		t.Fatal("expected current custody item to expose inactive asset state")
	}
}

func TestPostgresCustodyRepositoryListCurrentByAssetIncludesPersonnelActiveState(t *testing.T) {
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

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	if _, err := custodyRepository.SaveTransaction(context.Background(), transaction); err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	deactivated, err := personnelRepository.Deactivate(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating personnel, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected personnel to be deactivated")
	}

	holders, err := custodyRepository.ListCurrentByAsset(context.Background(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error listing asset holders, got %v", err)
	}

	if len(holders) != 1 {
		t.Fatalf("expected 1 current asset holder, got %d", len(holders))
	}

	if holders[0].PersonnelActive {
		t.Fatal("expected current asset holder to expose inactive personnel state")
	}

	if holders[0].PersonnelAlias != personnel.Alias() {
		t.Fatalf("expected personnel alias %s, got %s", personnel.Alias(), holders[0].PersonnelAlias)
	}

	if holders[0].PersonnelRank != personnel.Rank() {
		t.Fatalf("expected personnel rank %s, got %s", personnel.Rank(), holders[0].PersonnelRank)
	}
}

func TestPostgresCustodyRepositoryListPersonnelWithCurrentCustodyIncludesInactivePersonnel(t *testing.T) {
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
		domain.OperatorRoleOperator,
	)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	if !created {
		t.Fatal("expected transaction to be created")
	}

	deactivated, err := personnelRepository.Deactivate(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating personnel, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected personnel to be deactivated")
	}

	items, err := custodyRepository.ListPersonnelWithCurrentCustody(context.Background())
	if err != nil {
		t.Fatalf("expected no error listing personnel with current custody, got %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}

	if items[0].ID != personnel.ID() {
		t.Fatalf("expected personnel %s, got %s", personnel.ID(), items[0].ID)
	}

	if items[0].Active {
		t.Fatal("expected inactive personnel to be included with active=false")
	}

	if items[0].TotalQuantity != 1 {
		t.Fatalf("expected total quantity 1, got %d", items[0].TotalQuantity)
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

	if _, err := custodyRepository.SaveTransaction(context.Background(), transaction); err != nil {
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

func TestPostgresCustodyRepositoryFindReceiptIncludesCurrentRecordStatus(t *testing.T) {
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

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	if !created {
		t.Fatal("expected transaction to be created")
	}

	receipt, err := custodyRepository.FindReceiptByID(context.Background(), transaction.ID())
	if err != nil {
		t.Fatalf("expected no error finding receipt, got %v", err)
	}

	if !receipt.PersonnelActive {
		t.Fatal("expected personnel active status to be true")
	}

	if !receipt.OperatorActive {
		t.Fatal("expected operator active status to be true")
	}

	if len(receipt.Lines) != 1 {
		t.Fatalf("expected 1 receipt line, got %d", len(receipt.Lines))
	}

	if !receipt.Lines[0].AssetActive {
		t.Fatal("expected asset active status to be true")
	}
}

func TestPostgresCustodyRepositoryFindReceiptAfterPersonnelAndAssetDeactivation(t *testing.T) {
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

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	if !created {
		t.Fatal("expected transaction to be created")
	}

	deactivatedPersonnel, err := personnelRepository.Deactivate(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating personnel, got %v", err)
	}

	if !deactivatedPersonnel {
		t.Fatal("expected personnel to be deactivated")
	}

	deactivatedAsset, err := assetRepository.Deactivate(context.Background(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating asset, got %v", err)
	}

	if !deactivatedAsset {
		t.Fatal("expected asset to be deactivated")
	}

	receipt, err := custodyRepository.FindReceiptByID(context.Background(), transaction.ID())
	if err != nil {
		t.Fatalf("expected receipt to remain readable after deactivation, got %v", err)
	}

	if receipt.PersonnelActive {
		t.Fatal("expected receipt to expose current inactive personnel status")
	}

	if len(receipt.Lines) != 1 {
		t.Fatalf("expected 1 receipt line, got %d", len(receipt.Lines))
	}

	if receipt.Lines[0].AssetActive {
		t.Fatal("expected receipt to expose current inactive asset status")
	}

	if receipt.ID != transaction.ID() {
		t.Fatalf("expected receipt id %s, got %s", transaction.ID(), receipt.ID)
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

	if _, err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
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

	if _, err := custodyRepository.SaveTransaction(context.Background(), firstReturn); err != nil {
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

	if _, err := custodyRepository.SaveTransaction(context.Background(), finalReturn); err != nil {
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

	if _, err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
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

	_, err = custodyRepository.SaveTransaction(context.Background(), invalidReturn)
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

	if _, err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
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

func TestPostgresCustodyRepositoryConcurrentReturnsCannotOverdrawBalance(t *testing.T) {
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
		"transaction-checkout",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{checkoutLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid checkout, got %v", err)
	}

	if _, err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
		t.Fatalf("expected no error saving checkout, got %v", err)
	}

	errCh := make(chan error, 2)
	startCh := make(chan struct{})

	saveReturn := func(transactionID domain.CustodyTransactionID) {
		<-startCh

		returnLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
		if err != nil {
			errCh <- err
			return
		}

		returnTransaction, err := domain.NewCustodyTransaction(
			transactionID,
			domain.CustodyTransactionTypeReturn,
			personnel.ID(),
			operator.ID(),
			[]domain.CustodyLine{returnLine},
			"",
		)
		if err != nil {
			errCh <- err
			return
		}

		_, err = custodyRepository.SaveTransaction(context.Background(), returnTransaction)
		errCh <- err
	}

	go saveReturn("transaction-return-1")
	go saveReturn("transaction-return-2")

	close(startCh)

	firstErr := <-errCh
	secondErr := <-errCh

	successCount := 0
	insufficientCount := 0

	for _, err := range []error{firstErr, secondErr} {
		switch {
		case err == nil:
			successCount++
		case errors.Is(err, domain.ErrInsufficientCustodyBalance):
			insufficientCount++
		default:
			t.Fatalf("expected nil or ErrInsufficientCustodyBalance, got %v", err)
		}
	}

	if successCount != 1 {
		t.Fatalf("expected exactly 1 successful return, got %d", successCount)
	}

	if insufficientCount != 1 {
		t.Fatalf("expected exactly 1 insufficient balance error, got %d", insufficientCount)
	}

	quantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 0 {
		t.Fatalf("expected current quantity 0, got %d", quantity)
	}

	history, err := custodyRepository.ListHistoryByPersonnel(context.Background(), personnel.ID(), 10)
	if err != nil {
		t.Fatalf("expected no error listing history, got %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected checkout plus one successful return in history, got %d entries", len(history))
	}
}

func TestPostgresCustodyRepositoryConcurrentPartialReturnsRemainConsistent(t *testing.T) {
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
		"transaction-checkout",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{checkoutLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid checkout, got %v", err)
	}

	if _, err := custodyRepository.SaveTransaction(context.Background(), checkout); err != nil {
		t.Fatalf("expected no error saving checkout, got %v", err)
	}

	errCh := make(chan error, 2)
	startCh := make(chan struct{})

	saveReturn := func(transactionID domain.CustodyTransactionID, quantity int) {
		<-startCh

		returnLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(quantity))
		if err != nil {
			errCh <- err
			return
		}

		returnTransaction, err := domain.NewCustodyTransaction(
			transactionID,
			domain.CustodyTransactionTypeReturn,
			personnel.ID(),
			operator.ID(),
			[]domain.CustodyLine{returnLine},
			"",
		)
		if err != nil {
			errCh <- err
			return
		}

		_, err = custodyRepository.SaveTransaction(context.Background(), returnTransaction)
		errCh <- err
	}

	go saveReturn("transaction-return-1", 2)
	go saveReturn("transaction-return-2", 2)

	close(startCh)

	firstErr := <-errCh
	secondErr := <-errCh

	successCount := 0
	insufficientCount := 0

	for _, err := range []error{firstErr, secondErr} {
		switch {
		case err == nil:
			successCount++
		case errors.Is(err, domain.ErrInsufficientCustodyBalance):
			insufficientCount++
		default:
			t.Fatalf("expected nil or ErrInsufficientCustodyBalance, got %v", err)
		}
	}

	if successCount != 1 {
		t.Fatalf("expected exactly 1 successful return, got %d", successCount)
	}

	if insufficientCount != 1 {
		t.Fatalf("expected exactly 1 insufficient balance error, got %d", insufficientCount)
	}

	quantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 1 {
		t.Fatalf("expected current quantity 1, got %d", quantity)
	}
}

func TestPostgresCustodyRepositoryConcurrentCheckoutsAccumulateBalance(t *testing.T) {
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

	errCh := make(chan error, 2)
	startCh := make(chan struct{})

	saveCheckout := func(transactionID domain.CustodyTransactionID) {
		<-startCh

		line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
		if err != nil {
			errCh <- err
			return
		}

		transaction, err := domain.NewCustodyTransaction(
			transactionID,
			domain.CustodyTransactionTypeCheckout,
			personnel.ID(),
			operator.ID(),
			[]domain.CustodyLine{line},
			"",
		)
		if err != nil {
			errCh <- err
			return
		}

		_, err = custodyRepository.SaveTransaction(context.Background(), transaction)
		errCh <- err
	}

	go saveCheckout("transaction-checkout-1")
	go saveCheckout("transaction-checkout-2")

	close(startCh)

	firstErr := <-errCh
	secondErr := <-errCh

	if firstErr != nil {
		t.Fatalf("expected first checkout to succeed, got %v", firstErr)
	}

	if secondErr != nil {
		t.Fatalf("expected second checkout to succeed, got %v", secondErr)
	}

	quantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 2 {
		t.Fatalf("expected current quantity 2, got %d", quantity)
	}
}

func TestPostgresCustodyRepositorySaveTransactionIsIdempotentByTransactionID(t *testing.T) {
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

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	if !created {
		t.Fatal("expected first save to create transaction")
	}

	created, err = custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected duplicate save to be treated as idempotent, got %v", err)
	}

	if created {
		t.Fatal("expected second save to be idempotent duplicate")
	}

	quantity, err := custodyRepository.CurrentQuantity(context.Background(), personnel.ID(), asset.ID())
	if err != nil {
		t.Fatalf("expected no error reading current quantity, got %v", err)
	}

	if quantity != 1 {
		t.Fatalf("expected quantity to remain 1 after duplicate save, got %d", quantity)
	}

	history, err := custodyRepository.ListHistoryByPersonnel(context.Background(), personnel.ID(), 10)
	if err != nil {
		t.Fatalf("expected no error listing history, got %v", err)
	}

	if len(history) != 1 {
		t.Fatalf("expected only one transaction in history, got %d", len(history))
	}
}

func TestPostgresCustodyRepositoryListTransactionLedgerPeriods(t *testing.T) {
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

	operator := mustNewTestOperator(t, "operator-1", "93541134780", "silva", domain.RankSergeant, domain.OperatorRoleOperator)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transactions := []struct {
		id        domain.CustodyTransactionID
		createdAt time.Time
	}{
		{
			id:        "transaction-june",
			createdAt: time.Date(2026, time.June, 10, 9, 0, 0, 0, time.Local),
		},
		{
			id:        "transaction-july-a",
			createdAt: time.Date(2026, time.July, 3, 9, 0, 0, 0, time.Local),
		},
		{
			id:        "transaction-july-b",
			createdAt: time.Date(2026, time.July, 20, 9, 0, 0, 0, time.Local),
		},
	}

	for _, transactionData := range transactions {
		transaction, err := domain.NewCustodyTransaction(
			transactionData.id,
			domain.CustodyTransactionTypeCheckout,
			personnel.ID(),
			operator.ID(),
			[]domain.CustodyLine{line},
			"",
		)
		if err != nil {
			t.Fatalf("expected valid transaction, got %v", err)
		}

		created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
		if err != nil {
			t.Fatalf("expected no error saving transaction, got %v", err)
		}
		if !created {
			t.Fatal("expected transaction to be created")
		}

		_, err = pool.Exec(
			context.Background(),
			"UPDATE custody_transactions SET created_at = $1 WHERE id = $2",
			transactionData.createdAt,
			string(transactionData.id),
		)
		if err != nil {
			t.Fatalf("expected no error setting transaction period, got %v", err)
		}
	}

	periods, err := custodyRepository.ListTransactionLedgerPeriods(context.Background())
	if err != nil {
		t.Fatalf("expected no error listing ledger periods, got %v", err)
	}

	if len(periods) != 2 {
		t.Fatalf("expected 2 periods, got %d", len(periods))
	}

	if periods[0].Year != 2026 || periods[0].Month != 7 || periods[0].TransactionCount != 2 {
		t.Fatalf("expected July 2026 with 2 transactions, got %+v", periods[0])
	}

	if periods[1].Year != 2026 || periods[1].Month != 6 || periods[1].TransactionCount != 1 {
		t.Fatalf("expected June 2026 with 1 transaction, got %+v", periods[1])
	}
}

func TestPostgresCustodyRepositoryListTransactionSummariesFiltersByLedgerPeriod(t *testing.T) {
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

	operator := mustNewTestOperator(t, "operator-1", "93541134780", "silva", domain.RankSergeant, domain.OperatorRoleOperator)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transactions := []struct {
		id        domain.CustodyTransactionID
		createdAt time.Time
	}{
		{
			id:        "transaction-june",
			createdAt: time.Date(2026, time.June, 10, 9, 0, 0, 0, time.Local),
		},
		{
			id:        "transaction-july",
			createdAt: time.Date(2026, time.July, 3, 9, 0, 0, 0, time.Local),
		},
	}

	for _, transactionData := range transactions {
		transaction, err := domain.NewCustodyTransaction(
			transactionData.id,
			domain.CustodyTransactionTypeCheckout,
			personnel.ID(),
			operator.ID(),
			[]domain.CustodyLine{line},
			"",
		)
		if err != nil {
			t.Fatalf("expected valid transaction, got %v", err)
		}

		created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
		if err != nil {
			t.Fatalf("expected no error saving transaction, got %v", err)
		}
		if !created {
			t.Fatal("expected transaction to be created")
		}

		_, err = pool.Exec(
			context.Background(),
			"UPDATE custody_transactions SET created_at = $1 WHERE id = $2",
			transactionData.createdAt,
			string(transactionData.id),
		)
		if err != nil {
			t.Fatalf("expected no error setting transaction period, got %v", err)
		}
	}

	page, err := custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize:    10,
		PeriodStart: time.Date(2026, time.July, 1, 0, 0, 0, 0, time.Local),
		PeriodEnd:   time.Date(2026, time.August, 1, 0, 0, 0, 0, time.Local),
		HasPeriod:   true,
	})
	if err != nil {
		t.Fatalf("expected no error listing summaries, got %v", err)
	}

	summaries := page.Items

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}

	if summaries[0].ID != "transaction-july" {
		t.Fatalf("expected July transaction, got %s", summaries[0].ID)
	}

	if summaries[0].SequenceNumber != 2 {
		t.Fatalf("expected global sequence number 2, got %d", summaries[0].SequenceNumber)
	}
}

func TestPostgresCustodyRepositoryListTransactionSummariesDetectsNextPage(t *testing.T) {
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

	operator := mustNewTestOperator(t, "operator-1", "93541134780", "silva", domain.RankSergeant, domain.OperatorRoleOperator)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	for i := 1; i <= 3; i++ {
		transaction, err := domain.NewCustodyTransaction(
			domain.CustodyTransactionID("transaction-"+strconv.Itoa(i)),
			domain.CustodyTransactionTypeCheckout,
			personnel.ID(),
			operator.ID(),
			[]domain.CustodyLine{line},
			"",
		)
		if err != nil {
			t.Fatalf("expected valid transaction, got %v", err)
		}

		created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
		if err != nil {
			t.Fatalf("expected no error saving transaction, got %v", err)
		}

		if !created {
			t.Fatal("expected transaction to be created")
		}
	}

	page, err := custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize: 2,
	})
	if err != nil {
		t.Fatalf("expected no error listing transaction summaries, got %v", err)
	}

	if len(page.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(page.Items))
	}

	if !page.HasNextPage {
		t.Fatal("expected next page")
	}
}

func TestPostgresCustodyRepositoryListTransactionSummariesUsesEffectiveCorrectionState(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	personnelA := mustNewTestPersonnel(t, "personnel-1", "John A", "alpha", domain.RankSergeant, "52998224725")
	personnelB := mustNewTestPersonnel(t, "personnel-2", "John B", "bravo", domain.RankCorporal, "93541134780")

	if err := personnelRepository.Save(context.Background(), personnelA); err != nil {
		t.Fatalf("expected no error saving personnel A, got %v", err)
	}

	if err := personnelRepository.Save(context.Background(), personnelB); err != nil {
		t.Fatalf("expected no error saving personnel B, got %v", err)
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
		"29109142088",
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleOperator,
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnelA.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	if !created {
		t.Fatal("expected transaction to be created")
	}

	correctionLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(2))
	if err != nil {
		t.Fatalf("expected valid correction line, got %v", err)
	}

	correction, err := domain.NewCustodyCorrection(
		"correction-1",
		transaction.ID(),
		operator.ID(),
		personnelB.ID(),
		[]domain.CustodyLine{correctionLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid correction, got %v", err)
	}

	created, err = custodyRepository.SaveCorrection(
		context.Background(),
		correction,
		transaction.Type(),
		transaction.PersonnelID(),
		transaction.Lines(),
	)
	if err != nil {
		t.Fatalf("expected no error saving correction, got %v", err)
	}

	if !created {
		t.Fatal("expected correction to be created")
	}

	page, err := custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize: 10,
	})
	if err != nil {
		t.Fatalf("expected no error listing transaction summaries, got %v", err)
	}

	summaries := page.Items

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}

	summary := summaries[0]

	if summary.OriginalPersonnelID != personnelA.ID() {
		t.Fatalf("expected original personnel %s, got %s", personnelA.ID(), summary.OriginalPersonnelID)
	}

	if summary.EffectivePersonnelID != personnelB.ID() {
		t.Fatalf("expected effective personnel %s, got %s", personnelB.ID(), summary.EffectivePersonnelID)
	}

	if summary.TotalQuantity != 2 {
		t.Fatalf("expected total quantity 2, got %d", summary.TotalQuantity)
	}

	if !summary.HasCorrection {
		t.Fatal("expected summary to have correction")
	}

	if summary.EditCount != 1 {
		t.Fatalf("expected edit count 1, got %d", summary.EditCount)
	}

	if summary.SequenceNumber <= 0 {
		t.Fatalf("expected positive sequence number, got %d", summary.SequenceNumber)
	}

	if len(summary.Lines) != 1 {
		t.Fatalf("expected 1 effective line, got %d", len(summary.Lines))
	}

	if summary.Lines[0].AssetID != asset.ID() {
		t.Fatalf("expected asset %s, got %s", asset.ID(), summary.Lines[0].AssetID)
	}

	if summary.Lines[0].AssetName != asset.Name() {
		t.Fatalf("expected asset name %s, got %s", asset.Name(), summary.Lines[0].AssetName)
	}

	if summary.Lines[0].Quantity != 2 {
		t.Fatalf("expected effective line quantity 2, got %d", summary.Lines[0].Quantity)
	}
}

func TestPostgresCustodyRepositoryListTransactionSummariesFiltersByTransactionType(t *testing.T) {
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

	operator := mustNewTestOperator(t, "operator-1", "93541134780", "silva", domain.RankSergeant, domain.OperatorRoleOperator)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	checkout, err := domain.NewCustodyTransaction(
		"transaction-checkout",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid checkout, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), checkout)
	if err != nil {
		t.Fatalf("expected no error saving checkout, got %v", err)
	}
	if !created {
		t.Fatal("expected checkout to be created")
	}

	returnTransaction, err := domain.NewCustodyTransaction(
		"transaction-return",
		domain.CustodyTransactionTypeReturn,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid return, got %v", err)
	}

	created, err = custodyRepository.SaveTransaction(context.Background(), returnTransaction)
	if err != nil {
		t.Fatalf("expected no error saving return, got %v", err)
	}
	if !created {
		t.Fatal("expected return to be created")
	}

	page, err := custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize:              10,
		TransactionTypeFilter: ports.CustodyTransactionTypeFilterCheckout,
	})
	if err != nil {
		t.Fatalf("expected no error listing summaries, got %v", err)
	}

	summaries := page.Items

	if len(summaries) != 1 {
		t.Fatalf("expected 1 checkout summary, got %d", len(summaries))
	}

	if summaries[0].TransactionType != domain.CustodyTransactionTypeCheckout {
		t.Fatalf("expected checkout, got %s", summaries[0].TransactionType)
	}
}

func TestPostgresCustodyRepositoryListTransactionSummariesFiltersEditedTransactions(t *testing.T) {
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

	operator := mustNewTestOperator(t, "operator-1", "93541134780", "silva", domain.RankSergeant, domain.OperatorRoleOperator)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}
	if !created {
		t.Fatal("expected transaction to be created")
	}

	correctedLine, err := domain.NewCustodyLine(asset.ID(), domain.Quantity(2))
	if err != nil {
		t.Fatalf("expected valid corrected line, got %v", err)
	}

	correction, err := domain.NewCustodyCorrection(
		"correction-1",
		transaction.ID(),
		operator.ID(),
		personnel.ID(),
		[]domain.CustodyLine{correctedLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid correction, got %v", err)
	}

	created, err = custodyRepository.SaveCorrection(
		context.Background(),
		correction,
		transaction.Type(),
		transaction.PersonnelID(),
		transaction.Lines(),
	)
	if err != nil {
		t.Fatalf("expected no error saving correction, got %v", err)
	}
	if !created {
		t.Fatal("expected correction to be created")
	}

	page, err := custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize:         10,
		EditStatusFilter: ports.CustodyEditStatusFilterEdited,
	})
	if err != nil {
		t.Fatalf("expected no error listing edited summaries, got %v", err)
	}

	summaries := page.Items

	if len(summaries) != 1 {
		t.Fatalf("expected 1 edited summary, got %d", len(summaries))
	}

	if !summaries[0].HasCorrection {
		t.Fatal("expected summary to be edited")
	}
}

func TestPostgresCustodyRepositoryListTransactionSummariesSearchesEffectiveAssetName(t *testing.T) {
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

	operator := mustNewTestOperator(t, "operator-1", "93541134780", "silva", domain.RankSergeant, domain.OperatorRoleOperator)
	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	line, err := domain.NewCustodyLine(radio.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		"transaction-1",
		domain.CustodyTransactionTypeCheckout,
		personnel.ID(),
		operator.ID(),
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	created, err := custodyRepository.SaveTransaction(context.Background(), transaction)
	if err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}
	if !created {
		t.Fatal("expected transaction to be created")
	}

	correctedLine, err := domain.NewCustodyLine(helmet.ID(), domain.Quantity(1))
	if err != nil {
		t.Fatalf("expected valid corrected line, got %v", err)
	}

	correction, err := domain.NewCustodyCorrection(
		"correction-1",
		transaction.ID(),
		operator.ID(),
		personnel.ID(),
		[]domain.CustodyLine{correctedLine},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid correction, got %v", err)
	}

	created, err = custodyRepository.SaveCorrection(
		context.Background(),
		correction,
		transaction.Type(),
		transaction.PersonnelID(),
		transaction.Lines(),
	)
	if err != nil {
		t.Fatalf("expected no error saving correction, got %v", err)
	}
	if !created {
		t.Fatal("expected correction to be created")
	}

	page, err := custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize:    10,
		SearchQuery: "Helmet",
	})
	if err != nil {
		t.Fatalf("expected no error searching summaries, got %v", err)
	}

	summaries := page.Items

	if len(summaries) != 1 {
		t.Fatalf("expected 1 summary, got %d", len(summaries))
	}

	if len(summaries[0].Lines) != 1 {
		t.Fatalf("expected 1 effective line, got %d", len(summaries[0].Lines))
	}

	if summaries[0].Lines[0].AssetID != helmet.ID() {
		t.Fatalf("expected effective asset %s, got %s", helmet.ID(), summaries[0].Lines[0].AssetID)
	}

	page, err = custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize:    10,
		SearchQuery: "sergeant Helmet",
	})
	if err != nil {
		t.Fatalf("expected no error searching summaries by rank and asset, got %v", err)
	}

	if len(page.Items) != 1 {
		t.Fatalf("expected 1 summary by rank and asset, got %d", len(page.Items))
	}

	page, err = custodyRepository.ListTransactionSummaries(context.Background(), ports.CustodyTransactionSummaryFilters{
		PageSize:    10,
		SearchQuery: "logistics Helmet",
	})
	if err != nil {
		t.Fatalf("expected no error searching summaries by section and asset, got %v", err)
	}

	if len(page.Items) != 1 {
		t.Fatalf("expected 1 summary by section and asset, got %d", len(page.Items))
	}
}
