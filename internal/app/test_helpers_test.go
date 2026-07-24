package app

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) NewID() (string, error) {
	return g.id, nil
}

type fakePersonnelRepository struct {
	saved            []domain.Personnel
	byID             map[domain.PersonnelID]domain.Personnel
	lastStatusFilter ports.RecordStatusFilter
}

func (r *fakePersonnelRepository) Save(_ context.Context, personnel domain.Personnel) error {
	r.saved = append(r.saved, personnel)

	if r.byID == nil {
		r.byID = make(map[domain.PersonnelID]domain.Personnel)
	}

	r.byID[personnel.ID()] = personnel

	return nil
}

func (r *fakePersonnelRepository) Update(_ context.Context, personnel domain.Personnel) error {
	if r.byID == nil {
		r.byID = make(map[domain.PersonnelID]domain.Personnel)
	}

	r.byID[personnel.ID()] = personnel

	return nil
}

func (r *fakePersonnelRepository) FindByID(_ context.Context, id domain.PersonnelID) (domain.Personnel, error) {
	personnel, ok := r.byID[id]
	if !ok {
		return domain.Personnel{}, ports.ErrNotFound
	}

	return personnel, nil
}

func (r *fakePersonnelRepository) FindByRegistrationID(
	_ context.Context,
	registrationID domain.RegistrationID,
) (domain.Personnel, bool, error) {
	for _, personnel := range r.byID {
		if personnel.RegistrationID() == registrationID {
			return personnel, true, nil
		}
	}

	return domain.Personnel{}, false, nil
}

func (r *fakePersonnelRepository) FindByRegistrationIDExcludingID(
	_ context.Context,
	registrationID domain.RegistrationID,
	excludedID domain.PersonnelID,
) (domain.Personnel, bool, error) {
	for _, personnel := range r.byID {
		if personnel.ID() == excludedID {
			continue
		}

		if personnel.RegistrationID() == registrationID {
			return personnel, true, nil
		}
	}

	return domain.Personnel{}, false, nil
}

func (r *fakePersonnelRepository) List(
	_ context.Context,
	limit int,
	statusFilter ports.RecordStatusFilter,
) ([]domain.Personnel, error) {
	r.lastStatusFilter = statusFilter

	personnel := make([]domain.Personnel, 0, len(r.byID))

	for _, item := range r.byID {
		if !recordMatchesStatusFilter(item.Active(), statusFilter) {
			continue
		}

		personnel = append(personnel, item)
	}

	sort.Slice(personnel, func(i, j int) bool {
		return personnel[i].ID() < personnel[j].ID()
	})

	if limit > 0 && len(personnel) > limit {
		return personnel[:limit], nil
	}

	return personnel, nil
}

func (r *fakePersonnelRepository) Search(
	_ context.Context,
	query string,
	limit int,
	statusFilter ports.RecordStatusFilter,
) ([]domain.Personnel, error) {
	r.lastStatusFilter = statusFilter

	query = strings.ToLower(strings.TrimSpace(query))

	personnel := make([]domain.Personnel, 0, len(r.byID))

	for _, item := range r.byID {
		if !recordMatchesStatusFilter(item.Active(), statusFilter) {
			continue
		}

		if personnelMatchesQuery(item, query) {
			personnel = append(personnel, item)
		}
	}

	sort.Slice(personnel, func(i, j int) bool {
		return personnel[i].ID() < personnel[j].ID()
	})

	if limit > 0 && len(personnel) > limit {
		return personnel[:limit], nil
	}

	return personnel, nil
}

func (r *fakePersonnelRepository) Deactivate(_ context.Context, id domain.PersonnelID) (bool, error) {
	personnel, ok := r.byID[id]
	if !ok {
		return false, ports.ErrNotFound
	}

	deactivated, err := domain.ReconstitutePersonnel(
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
		return false, err
	}

	r.byID[id] = deactivated
	return true, nil
}

func (r *fakePersonnelRepository) Reactivate(_ context.Context, id domain.PersonnelID) (bool, error) {
	personnel, ok := r.byID[id]
	if !ok {
		return false, ports.ErrNotFound
	}

	reactivated, err := domain.ReconstitutePersonnel(
		personnel.ID(),
		personnel.FullName(),
		personnel.Alias(),
		personnel.Rank(),
		personnel.RegistrationID(),
		personnel.Section(),
		personnel.OrganizationUnit(),
		true,
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = reactivated
	return true, nil
}

func personnelMatchesQuery(personnel domain.Personnel, query string) bool {
	tokens := strings.Fields(strings.ToLower(strings.TrimSpace(query)))
	if len(tokens) == 0 {
		return true
	}

	values := []string{
		personnel.FullName(),
		personnel.Alias(),
		personnel.RegistrationID().String(),
		string(personnel.Rank()),
		string(personnel.Section()),
		string(personnel.OrganizationUnit()),
	}

	for _, token := range tokens {
		if !anyValueContainsToken(values, token) {
			return false
		}
	}

	return true
}

func anyValueContainsToken(values []string, token string) bool {
	for _, value := range values {
		if strings.Contains(strings.ToLower(value), token) {
			return true
		}
	}

	return false
}

type fakeAssetRepository struct {
	saved            []domain.Asset
	byID             map[domain.AssetID]domain.Asset
	lastStatusFilter ports.RecordStatusFilter
	saveErr          error
}

func (r *fakeAssetRepository) Save(_ context.Context, asset domain.Asset) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	r.saved = append(r.saved, asset)

	if r.byID == nil {
		r.byID = make(map[domain.AssetID]domain.Asset)
	}

	r.byID[asset.ID()] = asset

	return nil
}

func (r *fakeAssetRepository) Update(_ context.Context, asset domain.Asset) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	if r.byID == nil {
		r.byID = make(map[domain.AssetID]domain.Asset)
	}

	r.byID[asset.ID()] = asset

	return nil
}

func (r *fakeAssetRepository) FindByID(_ context.Context, id domain.AssetID) (domain.Asset, error) {
	asset, ok := r.byID[id]
	if !ok {
		return domain.Asset{}, ports.ErrNotFound
	}

	return asset, nil
}

func (r *fakeAssetRepository) FindByName(
	_ context.Context,
	name string,
) (domain.Asset, bool, error) {
	normalizedName := domain.NormalizeAssetName(name)

	for _, asset := range r.byID {
		if strings.EqualFold(asset.Name(), normalizedName) {
			return asset, true, nil
		}
	}

	return domain.Asset{}, false, nil
}

func (r *fakeAssetRepository) FindByNameExcludingID(
	_ context.Context,
	name string,
	excludedID domain.AssetID,
) (domain.Asset, bool, error) {
	normalizedName := domain.NormalizeAssetName(name)

	for _, asset := range r.byID {
		if asset.ID() == excludedID {
			continue
		}

		if strings.EqualFold(asset.Name(), normalizedName) {
			return asset, true, nil
		}
	}

	return domain.Asset{}, false, nil
}

func (r *fakeAssetRepository) List(
	_ context.Context,
	limit int,
	statusFilter ports.RecordStatusFilter,
) ([]domain.Asset, error) {
	r.lastStatusFilter = statusFilter

	assets := make([]domain.Asset, 0, len(r.byID))

	for _, item := range r.byID {
		if !recordMatchesStatusFilter(item.Active(), statusFilter) {
			continue
		}

		assets = append(assets, item)
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].ID() < assets[j].ID()
	})

	if limit > 0 && len(assets) > limit {
		return assets[:limit], nil
	}

	return assets, nil
}

func (r *fakeAssetRepository) Search(
	_ context.Context,
	query string,
	limit int,
	statusFilter ports.RecordStatusFilter,
) ([]domain.Asset, error) {
	r.lastStatusFilter = statusFilter

	tokens := strings.Fields(strings.ToLower(strings.TrimSpace(query)))

	assets := make([]domain.Asset, 0, len(r.byID))

	for _, item := range r.byID {
		if !recordMatchesStatusFilter(item.Active(), statusFilter) {
			continue
		}

		if assetMatchesQuery(item, tokens) {
			assets = append(assets, item)
		}
	}

	sort.Slice(assets, func(i, j int) bool {
		return assets[i].ID() < assets[j].ID()
	})

	if limit > 0 && len(assets) > limit {
		return assets[:limit], nil
	}

	return assets, nil
}

func (r *fakeAssetRepository) Deactivate(_ context.Context, id domain.AssetID) (bool, error) {
	asset, ok := r.byID[id]
	if !ok {
		return false, ports.ErrNotFound
	}

	deactivated, err := domain.ReconstituteAsset(
		asset.ID(),
		asset.Name(),
		false,
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = deactivated
	return true, nil
}

func (r *fakeAssetRepository) Reactivate(_ context.Context, id domain.AssetID) (bool, error) {
	asset, ok := r.byID[id]
	if !ok {
		return false, ports.ErrNotFound
	}

	reactivated, err := domain.ReconstituteAsset(
		asset.ID(),
		asset.Name(),
		true,
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = reactivated
	return true, nil
}

func assetMatchesQuery(asset domain.Asset, tokens []string) bool {
	if len(tokens) == 0 {
		return true
	}

	name := strings.ToLower(asset.Name())

	for _, token := range tokens {
		if !strings.Contains(name, token) {
			return false
		}
	}

	return true
}

func recordMatchesStatusFilter(active bool, statusFilter ports.RecordStatusFilter) bool {
	switch statusFilter {
	case ports.RecordStatusFilterAll:
		return true
	case ports.RecordStatusFilterInactive:
		return !active
	default:
		return active
	}
}

type fakeCustodyRepository struct {
	saved                       []domain.CustodyTransaction
	corrections                 []domain.CustodyCorrection
	saveErr                     error
	currentQuantity             map[string]int
	currentByPerson             map[domain.PersonnelID][]ports.CurrentCustodyItem
	currentByAsset              map[domain.AssetID][]ports.CurrentAssetHolder
	personnelWithCurrentCustody []ports.PersonnelWithCurrentCustody
	historyByPerson             map[domain.PersonnelID][]ports.CustodyHistoryEntry
	assetHistory                map[domain.AssetID][]ports.AssetCustodyHistoryItem
	transactionLedgerPeriods    []ports.CustodyTransactionLedgerPeriod
	transactionSummaries        []ports.CustodyTransactionSummary
	receipts                    map[domain.CustodyTransactionID]ports.CustodyReceipt
	correctionByTransactionID   map[domain.CustodyTransactionID]ports.CustodyCorrectionContext
	correctionsByTransactionID  map[domain.CustodyTransactionID][]ports.CustodyCorrectionContext
}

func (r *fakeCustodyRepository) SaveTransaction(_ context.Context, transaction domain.CustodyTransaction) (bool, error) {
	if r.saveErr != nil {
		return false, r.saveErr
	}

	for _, saved := range r.saved {
		if saved.ID() == transaction.ID() {
			return false, nil
		}
	}

	r.saved = append(r.saved, transaction)

	return true, nil
}

func (r *fakeCustodyRepository) SaveCorrection(
	_ context.Context,
	correction domain.CustodyCorrection,
	_ domain.CustodyTransactionType,
	_ domain.PersonnelID,
	_ []domain.CustodyLine,
) (bool, error) {
	if r.saveErr != nil {
		return false, r.saveErr
	}

	for _, saved := range r.corrections {
		if saved.ID() == correction.ID() {
			return false, nil
		}
	}

	r.corrections = append(r.corrections, correction)

	return true, nil
}

func (r *fakeCustodyRepository) CurrentQuantity(
	_ context.Context,
	personnelID domain.PersonnelID,
	assetID domain.AssetID,
) (int, error) {
	if r.currentQuantity == nil {
		return 0, nil
	}

	return r.currentQuantity[custodyBalanceKey(personnelID, assetID)], nil
}

func (r *fakeCustodyRepository) ListPersonnelWithCurrentCustody(
	_ context.Context,
) ([]ports.PersonnelWithCurrentCustody, error) {
	items := make([]ports.PersonnelWithCurrentCustody, len(r.personnelWithCurrentCustody))
	copy(items, r.personnelWithCurrentCustody)

	return items, nil
}

func (r *fakeCustodyRepository) ListCurrentByPersonnel(
	_ context.Context,
	personnelID domain.PersonnelID,
) ([]ports.CurrentCustodyItem, error) {
	if r.currentByPerson == nil {
		return nil, nil
	}

	items := r.currentByPerson[personnelID]
	copiedItems := make([]ports.CurrentCustodyItem, len(items))
	copy(copiedItems, items)

	return copiedItems, nil
}

func (r *fakeCustodyRepository) ListCurrentByAsset(
	_ context.Context,
	assetID domain.AssetID,
) ([]ports.CurrentAssetHolder, error) {
	if r.currentByAsset == nil {
		return nil, nil
	}

	holders := r.currentByAsset[assetID]
	copiedHolders := make([]ports.CurrentAssetHolder, len(holders))
	copy(copiedHolders, holders)

	return copiedHolders, nil
}

func (r *fakeCustodyRepository) ListHistoryByPersonnel(
	_ context.Context,
	personnelID domain.PersonnelID,
	limit int,
) ([]ports.CustodyHistoryEntry, error) {
	if r.historyByPerson == nil {
		return nil, nil
	}

	items := r.historyByPerson[personnelID]
	copiedItems := make([]ports.CustodyHistoryEntry, len(items))
	copy(copiedItems, items)

	if limit > 0 && len(copiedItems) > limit {
		return copiedItems[:limit], nil
	}

	return copiedItems, nil
}

func (r *fakeCustodyRepository) ListAssetCustodyHistory(
	_ context.Context,
	assetID domain.AssetID,
) ([]ports.AssetCustodyHistoryItem, error) {
	if r.assetHistory == nil {
		return nil, nil
	}

	items := r.assetHistory[assetID]
	copiedItems := make([]ports.AssetCustodyHistoryItem, len(items))
	copy(copiedItems, items)

	return copiedItems, nil
}

func (r *fakeCustodyRepository) ListTransactionSummaries(
	_ context.Context,
	_ ports.CustodyTransactionSummaryFilters,
) (ports.CustodyTransactionSummaryPage, error) {
	return ports.CustodyTransactionSummaryPage{
		Items:       r.transactionSummaries,
		HasNextPage: false,
	}, nil
}

func (r *fakeCustodyRepository) ListTransactionLedgerPeriods(
	_ context.Context,
) ([]ports.CustodyTransactionLedgerPeriod, error) {
	return r.transactionLedgerPeriods, nil
}

func (r *fakeCustodyRepository) FindReceiptByID(
	_ context.Context,
	id domain.CustodyTransactionID,
) (ports.CustodyReceipt, error) {
	if r.receipts == nil {
		return ports.CustodyReceipt{}, ports.ErrNotFound
	}

	receipt, ok := r.receipts[id]
	if !ok {
		return ports.CustodyReceipt{}, ports.ErrNotFound
	}

	return receipt, nil
}

func (r *fakeCustodyRepository) ListCorrectionContextsByTransactionID(
	_ context.Context,
	id domain.CustodyTransactionID,
) ([]ports.CustodyCorrectionContext, error) {
	if r.correctionsByTransactionID != nil {
		return r.correctionsByTransactionID[id], nil
	}

	if r.correctionByTransactionID != nil {
		correction, ok := r.correctionByTransactionID[id]
		if !ok {
			return nil, nil
		}

		return []ports.CustodyCorrectionContext{correction}, nil
	}

	receipt, ok := r.receipts[id]
	if ok && receipt.HasCorrection {
		return []ports.CustodyCorrectionContext{receipt.Correction}, nil
	}

	return nil, nil
}

func custodyBalanceKey(personnelID domain.PersonnelID, assetID domain.AssetID) string {
	return fmt.Sprintf("%s:%s", personnelID, assetID)
}

func mustBuildPersonnel(t *testing.T, id string) domain.Personnel {
	t.Helper()

	return mustBuildPersonnelWithRegistrationID(t, id, "52998224725")
}

func mustBuildPersonnelWithRegistrationID(t *testing.T, id string, registrationIDValue string) domain.Personnel {
	t.Helper()

	registrationID, err := domain.NewRegistrationID(registrationIDValue)
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	personnel, err := domain.NewPersonnel(
		domain.PersonnelID(id),
		"John Doe",
		"Doe",
		domain.PersonnelRankSergeant,
		registrationID,
		domain.PersonnelSectionOperations,
		domain.OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	return personnel
}

func validCreatePersonnelCommand(fullName string, alias string, registrationID string) CreatePersonnelCommand {
	return CreatePersonnelCommand{
		FullName:         fullName,
		Alias:            alias,
		Rank:             domain.PersonnelRankSergeant,
		RegistrationID:   registrationID,
		Section:          domain.PersonnelSectionOperations,
		OrganizationUnit: domain.OrganizationUnitDefault,
	}
}

func mustBuildAsset(t *testing.T, id domain.AssetID) domain.Asset {
	t.Helper()

	asset, err := domain.NewAsset(id, "Radio")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	return asset
}

func mustBuildOperator(
	t *testing.T,
	id domain.OperatorID,
	registrationIDValue string,
	alias string,
	rank domain.Rank,
	role domain.OperatorRole,
	passwordHash string,
) domain.Operator {
	t.Helper()

	operator, err := buildOperator(id, registrationIDValue, alias, rank, role, passwordHash)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	return operator
}

func mustNewTestOperator(
	t *testing.T,
	id domain.OperatorID,
	registrationID string,
	alias string,
	rank domain.Rank,
	role domain.OperatorRole,
) domain.Operator {
	t.Helper()

	validRegistrationID, err := domain.NewRegistrationID(registrationID)
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	operator, err := domain.NewOperator(
		id,
		validRegistrationID,
		alias,
		rank,
		role,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	return operator
}

func mustRegistrationID(t *testing.T, value string) domain.RegistrationID {
	t.Helper()

	registrationID, err := domain.NewRegistrationID(value)
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	return registrationID
}

func buildOperator(
	id domain.OperatorID,
	registrationIDValue string,
	alias string,
	rank domain.Rank,
	role domain.OperatorRole,
	passwordHash string,
) (domain.Operator, error) {
	registrationID, err := domain.NewRegistrationID(registrationIDValue)
	if err != nil {
		return domain.Operator{}, err
	}

	return domain.NewOperator(
		id,
		registrationID,
		alias,
		rank,
		role,
		passwordHash,
	)
}

func newFakeOperatorRepository(operators ...domain.Operator) *fakeOperatorRepository {
	repository := &fakeOperatorRepository{
		byID:             map[domain.OperatorID]domain.Operator{},
		byRegistrationID: map[domain.RegistrationID]domain.Operator{},
	}

	for _, operator := range operators {
		repository.byID[operator.ID()] = operator
		repository.byRegistrationID[operator.RegistrationID()] = operator
	}

	return repository
}

type fakeOperatorRepository struct {
	byID             map[domain.OperatorID]domain.Operator
	byRegistrationID map[domain.RegistrationID]domain.Operator
	summaries        []ports.OperatorSummary
	lastFilters      ports.OperatorFilters
	saveErr          error
}

func (r *fakeOperatorRepository) Save(_ context.Context, operator domain.Operator) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	if r.byID == nil {
		r.byID = map[domain.OperatorID]domain.Operator{}
	}

	if r.byRegistrationID == nil {
		r.byRegistrationID = map[domain.RegistrationID]domain.Operator{}
	}

	r.byID[operator.ID()] = operator
	r.byRegistrationID[operator.RegistrationID()] = operator

	return nil
}

func (r *fakeOperatorRepository) FindByID(_ context.Context, id domain.OperatorID) (domain.Operator, error) {
	operator, ok := r.byID[id]
	if !ok {
		return domain.Operator{}, ports.ErrNotFound
	}

	return operator, nil
}

func (r *fakeOperatorRepository) FindByRegistrationID(
	_ context.Context,
	registrationID domain.RegistrationID,
) (domain.Operator, error) {
	operator, ok := r.byRegistrationID[registrationID]
	if !ok {
		return domain.Operator{}, ports.ErrNotFound
	}

	return operator, nil
}

type fakePasswordHasher struct {
	hash string
	err  error
}

func (h fakePasswordHasher) Hash(_ string) (string, error) {
	if h.err != nil {
		return "", h.err
	}

	return h.hash, nil
}

func (h fakePasswordHasher) Verify(password string, encodedHash string) (bool, error) {
	return h.hash == encodedHash && password != "", nil
}

type fakeOperatorSessionRepository struct {
	byTokenHash map[string]domain.OperatorSession
}

func (r *fakeOperatorSessionRepository) Save(_ context.Context, session domain.OperatorSession) error {
	if r.byTokenHash == nil {
		r.byTokenHash = map[string]domain.OperatorSession{}
	}

	r.byTokenHash[session.TokenHash()] = session

	return nil
}

func (r *fakeOperatorSessionRepository) FindByTokenHash(_ context.Context, tokenHash string) (domain.OperatorSession, error) {
	session, ok := r.byTokenHash[tokenHash]
	if !ok {
		return domain.OperatorSession{}, ports.ErrNotFound
	}

	return session, nil
}

func (r *fakeOperatorSessionRepository) DeleteByTokenHash(_ context.Context, tokenHash string) error {
	delete(r.byTokenHash, tokenHash)
	return nil
}

func (r *fakeOperatorSessionRepository) DeleteExpired(_ context.Context, now time.Time) error {
	for tokenHash, session := range r.byTokenHash {
		if session.Expired(now) {
			delete(r.byTokenHash, tokenHash)
		}
	}

	return nil
}

func (r *fakeOperatorRepository) List(_ context.Context, filters ports.OperatorFilters) ([]ports.OperatorSummary, error) {
	r.lastFilters = filters

	summaries := r.summaries
	if summaries == nil {
		summaries = make([]ports.OperatorSummary, 0, len(r.byID))

		for _, operator := range r.byID {
			summaries = append(summaries, ports.OperatorSummary{
				ID:             operator.ID(),
				RegistrationID: operator.RegistrationID(),
				Alias:          operator.Alias(),
				Rank:           operator.Rank(),
				Role:           operator.Role(),
				Active:         operator.Active(),
			})
		}

		sort.Slice(summaries, func(i, j int) bool {
			return summaries[i].ID < summaries[j].ID
		})
	}

	query := strings.ToLower(strings.TrimSpace(filters.Query))
	filtered := make([]ports.OperatorSummary, 0, len(summaries))

	for _, summary := range summaries {
		if !recordMatchesStatusFilter(summary.Active, filters.Status) {
			continue
		}

		if query != "" && !operatorSummaryMatchesQuery(summary, query) {
			continue
		}

		filtered = append(filtered, summary)
	}

	if filters.Limit > 0 && filters.Limit < len(filtered) {
		return filtered[:filters.Limit], nil
	}

	return filtered, nil
}

func operatorSummaryMatchesQuery(summary ports.OperatorSummary, query string) bool {
	values := []string{
		summary.Alias,
		summary.RegistrationID.String(),
		summary.Rank.String(),
		summary.Role.String(),
	}

	return anyValueContainsToken(values, query)
}

func (r *fakeOperatorRepository) Deactivate(_ context.Context, id domain.OperatorID) (bool, error) {
	operator, ok := r.byID[id]
	if !ok {
		return false, nil
	}

	if !operator.Active() {
		return false, nil
	}

	deactivatedOperator, err := domain.ReconstituteOperator(
		operator.ID(),
		operator.RegistrationID(),
		operator.Alias(),
		operator.Rank(),
		operator.Role(),
		operator.PasswordHash(),
		false,
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = deactivatedOperator
	r.byRegistrationID[deactivatedOperator.RegistrationID()] = deactivatedOperator

	return true, nil
}

func (r *fakeOperatorRepository) ChangeRole(
	_ context.Context,
	id domain.OperatorID,
	role domain.OperatorRole,
) (bool, error) {
	operator, ok := r.byID[id]
	if !ok {
		return false, nil
	}

	changedOperator, err := domain.ReconstituteOperator(
		operator.ID(),
		operator.RegistrationID(),
		operator.Alias(),
		operator.Rank(),
		role,
		operator.PasswordHash(),
		operator.Active(),
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = changedOperator
	r.byRegistrationID[changedOperator.RegistrationID()] = changedOperator

	return true, nil
}

func (r *fakeOperatorRepository) CountActiveAdmins(_ context.Context) (int, error) {
	count := 0

	for _, operator := range r.byID {
		if operator.Active() && operator.Role() == domain.OperatorRoleAdmin {
			count++
		}
	}

	return count, nil
}

func (r *fakeOperatorSessionRepository) DeleteByOperatorID(_ context.Context, operatorID domain.OperatorID) error {
	for tokenHash, session := range r.byTokenHash {
		if session.OperatorID() == operatorID {
			delete(r.byTokenHash, tokenHash)
		}
	}

	return nil
}

func (r *fakeOperatorRepository) UpdatePasswordHash(
	_ context.Context,
	id domain.OperatorID,
	passwordHash string,
) (bool, error) {
	operator, ok := r.byID[id]
	if !ok {
		return false, nil
	}

	if !operator.Active() {
		return false, nil
	}

	updatedOperator, err := domain.ReconstituteOperator(
		operator.ID(),
		operator.RegistrationID(),
		operator.Alias(),
		operator.Rank(),
		operator.Role(),
		passwordHash,
		operator.Active(),
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = updatedOperator
	r.byRegistrationID[updatedOperator.RegistrationID()] = updatedOperator

	return true, nil
}

func (r *fakeOperatorRepository) FindSummaryByID(_ context.Context, id domain.OperatorID) (ports.OperatorSummary, error) {
	for _, summary := range r.summaries {
		if summary.ID == id {
			return summary, nil
		}
	}

	operator, ok := r.byID[id]
	if !ok {
		return ports.OperatorSummary{}, ports.ErrNotFound
	}

	return ports.OperatorSummary{
		ID:             operator.ID(),
		RegistrationID: operator.RegistrationID(),
		Alias:          operator.Alias(),
		Rank:           operator.Rank(),
		Role:           operator.Role(),
		Active:         operator.Active(),
	}, nil
}

func (r *fakeOperatorRepository) Reactivate(_ context.Context, id domain.OperatorID) (bool, error) {
	operator, ok := r.byID[id]
	if !ok {
		return false, nil
	}

	if operator.Active() {
		return false, nil
	}

	reactivatedOperator, err := domain.ReconstituteOperator(
		operator.ID(),
		operator.RegistrationID(),
		operator.Alias(),
		operator.Rank(),
		operator.Role(),
		operator.PasswordHash(),
		true,
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = reactivatedOperator
	r.byRegistrationID[reactivatedOperator.RegistrationID()] = reactivatedOperator

	return true, nil
}

type fixedSessionTokenGenerator struct {
	token string
	err   error
}

func (g fixedSessionTokenGenerator) NewToken() (string, error) {
	if g.err != nil {
		return "", g.err
	}

	return g.token, nil
}

type plainSessionTokenHasher struct{}

func (h plainSessionTokenHasher) Hash(token string) string {
	return "hash:" + token
}
