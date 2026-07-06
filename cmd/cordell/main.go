package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"cordell/internal/app"
	"cordell/internal/config"
	"cordell/internal/infra/ids"
	"cordell/internal/infra/postgres"
	postgresdb "cordell/internal/infra/postgres/db"
	"cordell/internal/web"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	pool, err := postgres.OpenPool(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("failed to open PostgreSQL pool", "error", err)
		os.Exit(1)

	}

	defer pool.Close()

	queries := postgresdb.New(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	idGenerator := ids.NewULIDGenerator()

	services := app.Services{
		CreatePersonnel: app.NewCreatePersonnelService(
			personnelRepository,
			idGenerator,
		),
		GetPersonnel: app.NewGetPersonnelService(
			personnelRepository,
		),
		ListPersonnel: app.NewListPersonnelService(
			personnelRepository,
		),
		CreateAsset: app.NewCreateAssetService(
			assetRepository,
			idGenerator,
		),
		RegisterCheckout: app.NewRegisterCheckoutService(
			personnelRepository,
			assetRepository,
			custodyRepository,
			idGenerator,
		),
		RegisterReturn: app.NewRegisterReturnService(
			personnelRepository,
			assetRepository,
			custodyRepository,
			idGenerator,
		),
	}

	server, err := web.NewServer(logger, services)
	if err != nil {
		logger.Error("failed to create web server", "error", err)
		os.Exit(1)
	}

	logger.Info("starting Cordell HTTP server", "address", cfg.HTTPAddress)

	if err := http.ListenAndServe(cfg.HTTPAddress, server.Routes()); err != nil {
		logger.Error("Cordell HTTP server stopped with error", "error", err)
		os.Exit(1)
	}
}
