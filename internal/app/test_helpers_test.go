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
	saved []domain.Personnel
	byID  map[domain.PersonnelID]domain.Personnel
}

func (r *fakePersonnelRepository) Save(_ context.Context, personnel domain.Personnel) error {
	r.saved = append(r.saved, personnel)

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

func (r *fakePersonnelRepository) List(_ context.Context, limit int) ([]domain.Personnel, error) {
	personnel := make([]domain.Personnel, 0, len(r.byID))

	for _, item := range r.byID {
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

func (r *fakePersonnelRepository) Search(_ context.Context, query string, limit int) ([]domain.Personnel, error) {
	query = strings.ToLower(strings.TrimSpace(query))

	personnel := make([]domain.Personnel, 0, len(r.byID))

	for _, item := range r.byID {
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
	saved []domain.Asset
	byID  map[domain.AssetID]domain.Asset
}

func (r *fakeAssetRepository) Save(_ context.Context, asset domain.Asset) error {
	r.saved = append(r.saved, asset)

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

func (r *fakeAssetRepository) List(_ context.Context, limit int) ([]domain.Asset, error) {
	assets := make([]domain.Asset, 0, len(r.byID))

	for _, item := range r.byID {
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

func (r *fakeAssetRepository) Search(_ context.Context, query string, limit int) ([]domain.Asset, error) {
	tokens := strings.Fields(strings.ToLower(strings.TrimSpace(query)))

	assets := make([]domain.Asset, 0, len(r.byID))

	for _, item := range r.byID {
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

type fakeCustodyRepository struct {
	saved           []domain.CustodyTransaction
	currentQuantity map[string]int
	currentByPerson map[domain.PersonnelID][]ports.CurrentCustodyItem
	currentByAsset  map[domain.AssetID][]ports.CurrentAssetHolder
	historyByPerson map[domain.PersonnelID][]ports.CustodyHistoryEntry
}

func (r *fakeCustodyRepository) SaveTransaction(_ context.Context, transaction domain.CustodyTransaction) error {
	r.saved = append(r.saved, transaction)

	return nil
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

func custodyBalanceKey(personnelID domain.PersonnelID, assetID domain.AssetID) string {
	return fmt.Sprintf("%s:%s", personnelID, assetID)
}

func mustBuildPersonnel(t *testing.T, id string) domain.Personnel {
	t.Helper()

	registrationID, err := domain.NewRegistrationID("52998224725")
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

type fakeOperatorRepository struct {
	byID       map[domain.OperatorID]domain.Operator
	byUsername map[string]domain.Operator
	summaries  []ports.OperatorSummary
	saveErr    error
}

func (r *fakeOperatorRepository) Save(_ context.Context, operator domain.Operator) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	if r.byID == nil {
		r.byID = map[domain.OperatorID]domain.Operator{}
	}

	if r.byUsername == nil {
		r.byUsername = map[string]domain.Operator{}
	}

	r.byID[operator.ID()] = operator
	r.byUsername[operator.Username()] = operator

	return nil
}

func (r *fakeOperatorRepository) FindByID(_ context.Context, id domain.OperatorID) (domain.Operator, error) {
	operator, ok := r.byID[id]
	if !ok {
		return domain.Operator{}, ports.ErrNotFound
	}

	return operator, nil
}

func (r *fakeOperatorRepository) FindByUsername(_ context.Context, username string) (domain.Operator, error) {
	operator, ok := r.byUsername[domain.NormalizeOperatorUsername(username)]
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

func (r *fakeOperatorRepository) List(_ context.Context, limit int) ([]ports.OperatorSummary, error) {
	if limit > len(r.summaries) {
		limit = len(r.summaries)
	}

	return r.summaries[:limit], nil
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
		operator.Username(),
		operator.Role(),
		operator.PasswordHash(),
		false,
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = deactivatedOperator
	r.byUsername[deactivatedOperator.Username()] = deactivatedOperator

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
		operator.Username(),
		role,
		operator.PasswordHash(),
		operator.Active(),
	)
	if err != nil {
		return false, err
	}

	r.byID[id] = changedOperator
	r.byUsername[changedOperator.Username()] = changedOperator

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
