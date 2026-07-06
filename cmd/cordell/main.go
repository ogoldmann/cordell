package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"cordell/internal/app"
	"cordell/internal/config"
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

	services := app.Services{
		CreatePersonnel: app.NewCreatePersonnelService(
			personnelRepository,
			staticIDGenerator{},
		),
		CreateAsset: app.NewCreateAssetService(
			assetRepository,
			staticIDGenerator{},
		),
		RegisterCheckout: app.NewRegisterCheckoutService(
			personnelRepository,
			assetRepository,
			custodyRepository,
			staticIDGenerator{},
		),
		RegisterReturn: app.NewRegisterReturnService(
			personnelRepository,
			assetRepository,
			custodyRepository,
			staticIDGenerator{},
		),
	}

	server := web.NewServer(logger, services)

	logger.Info("starting Cordell HTTP server", "address", cfg.HTTPAddress)

	if err := http.ListenAndServe(cfg.HTTPAddress, server.Routes()); err != nil {
		logger.Error("Cordell HTTP server stopped with error", "error", err)
		os.Exit(1)
	}
}

type staticIDGenerator struct{}

func (g staticIDGenerator) NewID() string {
	return "temporary-id"
}
