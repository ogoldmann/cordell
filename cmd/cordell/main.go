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
	"cordell/internal/ports"
	"cordell/internal/security"
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

	operatorRepository := postgres.NewOperatorRepository(queries)
	passwordHasher := security.NewDefaultArgon2idPasswordHasher()

	auditLogRepository := postgres.NewAuditLogRepository(queries)

	sessionRepository := postgres.NewOperatorSessionRepository(queries)
	sessionTokenGenerator := security.NewDefaultRandomSessionTokenGenerator()
	sessionTokenHasher := security.NewSHA256SessionTokenHasher()
	csrfTokenGenerator := security.NewDefaultRandomSessionTokenGenerator()

	idGenerator := ids.NewULIDGenerator()

	services := newAppServices(
		personnelRepository,
		assetRepository,
		custodyRepository,
		operatorRepository,
		passwordHasher,
		auditLogRepository,
		sessionRepository,
		sessionTokenGenerator,
		sessionTokenHasher,
		csrfTokenGenerator,
		idGenerator,
	)

	server, err := web.NewServer(
		logger,
		services,
		web.NewSessionCookieConfig(cfg.SessionCookieSecure),
		web.NewSecurityHeadersConfig(cfg.EnableHSTS),
	)
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

func newAppServices(
	personnelRepository ports.PersonnelRepository,
	assetRepository ports.AssetRepository,
	custodyRepository ports.CustodyRepository,
	operatorRepository ports.OperatorRepository,
	passwordHasher ports.PasswordHasher,
	auditLogRepository ports.AuditLogRepository,
	sessionRepository ports.OperatorSessionRepository,
	sessionTokenGenerator ports.SessionTokenGenerator,
	sessionTokenHasher ports.SessionTokenHasher,
	csrfTokenGenerator ports.SessionTokenGenerator,
	idGenerator ports.IDGenerator,
) app.Services {
	return app.Services{
		CreateOperator: app.NewCreateOperatorService(
			operatorRepository,
			idGenerator,
			passwordHasher,
		),
		GetOperatorAdmin: app.NewGetOperatorAdminService(operatorRepository),
		ListOperators:    app.NewListOperatorsService(operatorRepository),
		DeactivateOperator: app.NewDeactivateOperatorService(
			operatorRepository,
			sessionRepository,
		),
		ReactivateOperator: app.NewReactivateOperatorService(
			operatorRepository,
			sessionRepository,
		),
		ChangeOperatorRole: app.NewChangeOperatorRoleService(
			operatorRepository,
			sessionRepository,
		),
		ResetOperatorPassword: app.NewResetOperatorPasswordService(
			operatorRepository,
			sessionRepository,
			passwordHasher,
		),
		AuthenticateOperator: app.NewAuthenticateOperatorService(
			operatorRepository,
			passwordHasher,
		),
		CreateOperatorSession: app.NewCreateOperatorSessionService(
			sessionRepository,
			idGenerator,
			sessionTokenGenerator,
			csrfTokenGenerator,
			sessionTokenHasher,
		),
		GetOperatorBySessionToken: app.NewGetOperatorBySessionTokenService(
			sessionRepository,
			operatorRepository,
			sessionTokenHasher,
		),
		DeleteOperatorSession: app.NewDeleteOperatorSessionService(
			sessionRepository,
			sessionTokenHasher,
		),
		DeleteExpiredOperatorSessions: app.NewDeleteExpiredOperatorSessionsService(
			sessionRepository,
		),
		CreatePersonnel: app.NewCreatePersonnelService(
			personnelRepository,
			idGenerator,
		),
		GetPersonnel: app.NewGetPersonnelService(
			personnelRepository,
		),
		UpdatePersonnel: app.NewUpdatePersonnelService(
			personnelRepository,
		),
		ListPersonnel: app.NewListPersonnelService(
			personnelRepository,
		),
		DeactivatePersonnel: app.NewDeactivatePersonnelService(
			personnelRepository,
		),
		ReactivatePersonnel: app.NewReactivatePersonnelService(
			personnelRepository,
		),
		SearchPersonnel: app.NewSearchPersonnelService(
			personnelRepository,
		),
		CreateAsset: app.NewCreateAssetService(
			assetRepository,
			idGenerator,
		),
		GetAsset: app.NewGetAssetService(
			assetRepository,
		),
		UpdateAsset: app.NewUpdateAssetService(
			assetRepository,
		),
		ListAssets: app.NewListAssetsService(
			assetRepository,
		),
		DeactivateAsset: app.NewDeactivateAssetService(
			assetRepository,
		),
		ReactivateAsset: app.NewReactivateAssetService(
			assetRepository,
		),
		SearchAssets: app.NewSearchAssetsService(
			assetRepository,
		),
		GlobalSearch: app.NewGlobalSearchService(
			personnelRepository,
			assetRepository,
		),
		RegisterCheckout: app.NewRegisterCheckoutService(
			personnelRepository,
			assetRepository,
			operatorRepository,
			custodyRepository,
			idGenerator,
		),
		RegisterReturn: app.NewRegisterReturnService(
			personnelRepository,
			assetRepository,
			operatorRepository,
			custodyRepository,
			idGenerator,
		),
		RegisterCustodyCorrection: app.NewRegisterCustodyCorrectionService(
			personnelRepository,
			assetRepository,
			operatorRepository,
			custodyRepository,
		),
		GetCustodyReceipt: app.NewGetCustodyReceiptService(custodyRepository),
		ListPersonnelWithCurrentCustody: app.NewListPersonnelWithCurrentCustodyService(
			custodyRepository,
		),
		ListCurrentCustody: app.NewListCurrentCustodyService(
			personnelRepository,
			custodyRepository,
		),
		ListCurrentCustodySummary: app.NewListCurrentCustodySummaryService(
			custodyRepository,
		),
		ListCurrentCustodianSummary: app.NewListCurrentCustodianSummaryService(
			custodyRepository,
		),
		ListCurrentAssetHolders: app.NewListCurrentAssetHoldersService(
			assetRepository,
			custodyRepository,
		),
		ListCustodyHistory: app.NewListCustodyHistoryService(
			personnelRepository,
			custodyRepository,
		),
		ListAssetCustodyHistory: app.NewListAssetCustodyHistoryService(
			assetRepository,
			custodyRepository,
		),
		ListCustodyTransactionLedgerPeriods: app.NewListCustodyTransactionLedgerPeriodsService(
			custodyRepository,
		),
		ListCustodyTransactionSummaries: app.NewListCustodyTransactionSummariesService(
			custodyRepository,
		),
		RecordAuditEvent: app.NewRecordAuditEventService(
			auditLogRepository,
			idGenerator,
		),
		ListAuditEvents: app.NewListAuditEventsService(auditLogRepository),
	}
}
