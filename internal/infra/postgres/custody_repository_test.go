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
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

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
		custodyRepository,
		fixedIDGenerator{id: "transaction-1"},
	)

	_, err := createPersonnelService.Execute(context.Background(), app.CreatePersonnelCommand{
		FullName: "John Doe",
	})
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
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

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
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-1"},
	)
	registerReturnService := app.NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-return-1"},
	)

	_, err := createPersonnelService.Execute(context.Background(), app.CreatePersonnelCommand{
		FullName: "John Doe",
	})
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
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

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
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-1"},
	)
	registerReturnService := app.NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-return-1"},
	)

	_, err := createPersonnelService.Execute(context.Background(), app.CreatePersonnelCommand{
		FullName: "John Doe",
	})
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
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

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
		custodyRepository,
		fixedIDGenerator{id: "transaction-checkout-1"},
	)
	registerReturnService := app.NewRegisterReturnService(
		personnelRepository,
		assetRepository,
		custodyRepository,
		fixedIDGenerator{id: "transaction-return-1"},
	)

	_, err := createPersonnelService.Execute(context.Background(), app.CreatePersonnelCommand{
		FullName: "John Doe",
	})
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
