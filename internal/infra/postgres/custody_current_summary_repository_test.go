package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
)

func TestPostgresCustodyRepositoryListCurrentCustodySummariesForCheckout(t *testing.T) {
	custodyRepository := newCurrentSummaryFixture(t)

	custodyRepository.saveCheckout(t, "checkout-1", "personnel-1", "asset-1", 3)

	personnelSummaries := custodyRepository.personnelSummaries(t)
	if got := personnelSummaries["personnel-1"]; got != 3 {
		t.Fatalf("expected personnel-1 current custody quantity 3, got %d", got)
	}

	assetSummaries := custodyRepository.assetSummaries(t)
	if got := assetSummaries["asset-1"]; got != 1 {
		t.Fatalf("expected asset-1 current custodian count 1, got %d", got)
	}
}

func TestPostgresCustodyRepositoryListCurrentCustodySummariesForPartialReturn(t *testing.T) {
	custodyRepository := newCurrentSummaryFixture(t)

	custodyRepository.saveCheckout(t, "checkout-1", "personnel-1", "asset-1", 3)
	custodyRepository.saveReturn(t, "return-1", "personnel-1", "asset-1", 1)

	personnelSummaries := custodyRepository.personnelSummaries(t)
	if got := personnelSummaries["personnel-1"]; got != 2 {
		t.Fatalf("expected personnel-1 current custody quantity 2, got %d", got)
	}

	assetSummaries := custodyRepository.assetSummaries(t)
	if got := assetSummaries["asset-1"]; got != 1 {
		t.Fatalf("expected asset-1 current custodian count 1, got %d", got)
	}
}

func TestPostgresCustodyRepositoryListCurrentCustodySummariesOmitsZeroBalances(t *testing.T) {
	custodyRepository := newCurrentSummaryFixture(t)

	custodyRepository.saveCheckout(t, "checkout-1", "personnel-1", "asset-1", 3)
	custodyRepository.saveReturn(t, "return-1", "personnel-1", "asset-1", 3)

	personnelSummaries := custodyRepository.personnelSummaries(t)
	if got := personnelSummaries["personnel-1"]; got != 0 {
		t.Fatalf("expected personnel-1 current custody quantity 0, got %d", got)
	}

	assetSummaries := custodyRepository.assetSummaries(t)
	if got := assetSummaries["asset-1"]; got != 0 {
		t.Fatalf("expected asset-1 current custodian count 0, got %d", got)
	}
}

func TestPostgresCustodyRepositoryListCurrentCustodySummariesCountsDistinctPersonnelByAsset(t *testing.T) {
	custodyRepository := newCurrentSummaryFixture(t)

	custodyRepository.saveCheckout(t, "checkout-1", "personnel-1", "asset-1", 2)
	custodyRepository.saveCheckout(t, "checkout-2", "personnel-2", "asset-1", 1)

	assetSummaries := custodyRepository.assetSummaries(t)
	if got := assetSummaries["asset-1"]; got != 2 {
		t.Fatalf("expected asset-1 current custodian count 2, got %d", got)
	}
}

func TestPostgresCustodyRepositoryListCurrentCustodySummariesUsesCorrectedQuantity(t *testing.T) {
	custodyRepository := newCurrentSummaryFixture(t)

	transaction := custodyRepository.saveCheckout(t, "checkout-1", "personnel-1", "asset-1", 1)
	custodyRepository.saveCorrection(t, "correction-1", transaction, "personnel-1", "asset-1", 4)

	personnelSummaries := custodyRepository.personnelSummaries(t)
	if got := personnelSummaries["personnel-1"]; got != 4 {
		t.Fatalf("expected personnel-1 current custody quantity 4, got %d", got)
	}

	assetSummaries := custodyRepository.assetSummaries(t)
	if got := assetSummaries["asset-1"]; got != 1 {
		t.Fatalf("expected asset-1 current custodian count 1, got %d", got)
	}
}

func TestPostgresCustodyRepositoryListCurrentCustodySummariesUsesCorrectedPersonnel(t *testing.T) {
	custodyRepository := newCurrentSummaryFixture(t)

	transaction := custodyRepository.saveCheckout(t, "checkout-1", "personnel-1", "asset-1", 2)
	custodyRepository.saveCorrection(t, "correction-1", transaction, "personnel-2", "asset-1", 2)

	personnelSummaries := custodyRepository.personnelSummaries(t)
	if got := personnelSummaries["personnel-1"]; got != 0 {
		t.Fatalf("expected personnel-1 current custody quantity 0, got %d", got)
	}

	if got := personnelSummaries["personnel-2"]; got != 2 {
		t.Fatalf("expected personnel-2 current custody quantity 2, got %d", got)
	}
}

type currentSummaryFixture struct {
	repository *postgres.CustodyRepository
	operatorID domain.OperatorID
}

func newCurrentSummaryFixture(t *testing.T) currentSummaryFixture {
	t.Helper()

	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	assetRepository := postgres.NewAssetRepository(queries)
	operatorRepository := postgres.NewOperatorRepository(queries)
	custodyRepository := postgres.NewCustodyRepository(pool, queries)

	firstPersonnel := mustNewTestPersonnel(t, "personnel-1", "John Silva", "silva", domain.RankSergeant, "52998224725")
	secondPersonnel := mustNewTestPersonnel(t, "personnel-2", "Jane Santos", "santos", domain.RankCorporal, "93541134780")

	if err := personnelRepository.Save(context.Background(), firstPersonnel); err != nil {
		t.Fatalf("expected no error saving first personnel, got %v", err)
	}

	if err := personnelRepository.Save(context.Background(), secondPersonnel); err != nil {
		t.Fatalf("expected no error saving second personnel, got %v", err)
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
		"operator",
		domain.RankSergeant,
		domain.OperatorRoleOperator,
	)

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	return currentSummaryFixture{
		repository: custodyRepository,
		operatorID: operator.ID(),
	}
}

func (f currentSummaryFixture) saveCheckout(
	t *testing.T,
	id domain.CustodyTransactionID,
	personnelID domain.PersonnelID,
	assetID domain.AssetID,
	quantity domain.Quantity,
) domain.CustodyTransaction {
	t.Helper()

	return f.saveTransaction(t, id, domain.CustodyTransactionTypeCheckout, personnelID, assetID, quantity)
}

func (f currentSummaryFixture) saveReturn(
	t *testing.T,
	id domain.CustodyTransactionID,
	personnelID domain.PersonnelID,
	assetID domain.AssetID,
	quantity domain.Quantity,
) domain.CustodyTransaction {
	t.Helper()

	return f.saveTransaction(t, id, domain.CustodyTransactionTypeReturn, personnelID, assetID, quantity)
}

func (f currentSummaryFixture) saveTransaction(
	t *testing.T,
	id domain.CustodyTransactionID,
	transactionType domain.CustodyTransactionType,
	personnelID domain.PersonnelID,
	assetID domain.AssetID,
	quantity domain.Quantity,
) domain.CustodyTransaction {
	t.Helper()

	line, err := domain.NewCustodyLine(assetID, quantity)
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	transaction, err := domain.NewCustodyTransaction(
		id,
		transactionType,
		personnelID,
		f.operatorID,
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid transaction, got %v", err)
	}

	if _, err := f.repository.SaveTransaction(context.Background(), transaction); err != nil {
		t.Fatalf("expected no error saving transaction, got %v", err)
	}

	return transaction
}

func (f currentSummaryFixture) saveCorrection(
	t *testing.T,
	id domain.CustodyCorrectionID,
	transaction domain.CustodyTransaction,
	correctedPersonnelID domain.PersonnelID,
	assetID domain.AssetID,
	quantity domain.Quantity,
) {
	t.Helper()

	line, err := domain.NewCustodyLine(assetID, quantity)
	if err != nil {
		t.Fatalf("expected valid correction line, got %v", err)
	}

	correction, err := domain.NewCustodyCorrection(
		id,
		transaction.ID(),
		f.operatorID,
		correctedPersonnelID,
		[]domain.CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected valid correction, got %v", err)
	}

	if _, err := f.repository.SaveCorrection(
		context.Background(),
		correction,
		transaction.Type(),
		transaction.PersonnelID(),
		transaction.Lines(),
	); err != nil {
		t.Fatalf("expected no error saving correction, got %v", err)
	}
}

func (f currentSummaryFixture) personnelSummaries(t *testing.T) map[domain.PersonnelID]int64 {
	t.Helper()

	summaries, err := f.repository.ListCurrentCustodySummaryByPersonnel(context.Background())
	if err != nil {
		t.Fatalf("expected no error listing personnel summaries, got %v", err)
	}

	byPersonnelID := make(map[domain.PersonnelID]int64, len(summaries))
	for _, summary := range summaries {
		byPersonnelID[summary.PersonnelID] = summary.CurrentCustodyQuantity
	}

	return byPersonnelID
}

func (f currentSummaryFixture) assetSummaries(t *testing.T) map[domain.AssetID]int64 {
	t.Helper()

	summaries, err := f.repository.ListCurrentCustodySummaryByAsset(context.Background())
	if err != nil {
		t.Fatalf("expected no error listing asset summaries, got %v", err)
	}

	byAssetID := make(map[domain.AssetID]int64, len(summaries))
	for _, summary := range summaries {
		byAssetID[summary.AssetID] = summary.CurrentCustodianCount
	}

	return byAssetID
}
