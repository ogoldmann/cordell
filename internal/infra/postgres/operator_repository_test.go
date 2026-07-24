package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
	"cordell/internal/ports"
)

func buildTestOperator(
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

	return domain.NewOperator(id, registrationID, alias, rank, role, passwordHash)
}

func TestPostgresOperatorRepositorySaveAndFindByRegistrationID(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	operator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	found, err := operatorRepository.FindByRegistrationID(context.Background(), operator.RegistrationID())
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if found.Role() != domain.OperatorRoleAdmin {
		t.Fatalf("expected admin role, got %s", found.Role())
	}

	if found.ID() != "operator-1" {
		t.Fatalf("expected operator-1, got %s", found.ID())
	}

	if found.RegistrationID() != operator.RegistrationID() {
		t.Fatalf("expected registration id %s, got %s", operator.RegistrationID(), found.RegistrationID())
	}

	if found.Alias() != "silva" {
		t.Fatalf("expected alias silva, got %s", found.Alias())
	}

}

func TestPostgresOperatorRepositoryRejectsDuplicateRegistrationID(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	firstOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid first operator, got %v", err)
	}

	secondOperator, err := buildTestOperator("operator-2", "52998224725", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid second operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), firstOperator); err != nil {
		t.Fatalf("expected no error saving first operator, got %v", err)
	}

	err = operatorRepository.Save(context.Background(), secondOperator)
	if err != domain.ErrDuplicateRegistrationID {
		t.Fatalf("expected ErrDuplicateRegistrationID, got %v", err)
	}
}

func TestPostgresOperatorRepositoryList(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	regularOperator, err := buildTestOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid regular operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), regularOperator); err != nil {
		t.Fatalf("expected no error saving regular operator, got %v", err)
	}

	operators, err := operatorRepository.List(context.Background(), ports.OperatorFilters{
		Status: ports.RecordStatusFilterAll,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("expected no error listing operators, got %v", err)
	}

	if len(operators) != 2 {
		t.Fatalf("expected 2 operators, got %d", len(operators))
	}

	for _, operator := range operators {
		if operator.RegistrationID.String() == "" {
			t.Fatal("expected registration id not to be empty")
		}

		if operator.Alias == "" {
			t.Fatal("expected alias not to be empty")
		}

		if operator.Rank == "" {
			t.Fatal("expected rank not to be empty")
		}

		if operator.Role == "" {
			t.Fatal("expected role not to be empty")
		}

		if operator.CreatedAt.IsZero() {
			t.Fatal("expected created_at not to be zero")
		}
	}
}

func TestPostgresOperatorRepositoryListFiltersByQueryAndStatus(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	regularOperator, err := buildTestOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid regular operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), regularOperator); err != nil {
		t.Fatalf("expected no error saving regular operator, got %v", err)
	}

	deactivated, err := operatorRepository.Deactivate(context.Background(), regularOperator.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating regular operator, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected regular operator to be deactivated")
	}

	activeOperators, err := operatorRepository.List(context.Background(), ports.OperatorFilters{
		Query:  "silva",
		Status: ports.RecordStatusFilterActive,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("expected no error listing active operators, got %v", err)
	}

	if len(activeOperators) != 1 {
		t.Fatalf("expected 1 active operator, got %d", len(activeOperators))
	}

	if activeOperators[0].ID != adminOperator.ID() {
		t.Fatalf("expected %s, got %s", adminOperator.ID(), activeOperators[0].ID)
	}

	inactiveOperators, err := operatorRepository.List(context.Background(), ports.OperatorFilters{
		Query:  "operator",
		Status: ports.RecordStatusFilterInactive,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("expected no error listing inactive operators, got %v", err)
	}

	if len(inactiveOperators) != 1 {
		t.Fatalf("expected 1 inactive operator, got %d", len(inactiveOperators))
	}

	if inactiveOperators[0].ID != regularOperator.ID() {
		t.Fatalf("expected %s, got %s", regularOperator.ID(), inactiveOperators[0].ID)
	}
}

func TestPostgresOperatorRepositoryDeactivate(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	regularOperator, err := buildTestOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid regular operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), regularOperator); err != nil {
		t.Fatalf("expected no error saving regular operator, got %v", err)
	}

	deactivated, err := operatorRepository.Deactivate(context.Background(), regularOperator.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating operator, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected operator to be deactivated")
	}

	found, err := operatorRepository.FindByID(context.Background(), regularOperator.ID())
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if found.Active() {
		t.Fatal("expected operator to be inactive")
	}
}

func TestPostgresOperatorRepositoryDoesNotDeactivateLastAdmin(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	deactivated, err := operatorRepository.Deactivate(context.Background(), adminOperator.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating operator, got %v", err)
	}

	if deactivated {
		t.Fatal("expected last admin not to be deactivated")
	}

	found, err := operatorRepository.FindByID(context.Background(), adminOperator.ID())
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if !found.Active() {
		t.Fatal("expected last admin to remain active")
	}
}

func TestPostgresOperatorRepositoryCountActiveAdmins(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	regularOperator, err := buildTestOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid regular operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), regularOperator); err != nil {
		t.Fatalf("expected no error saving regular operator, got %v", err)
	}

	count, err := operatorRepository.CountActiveAdmins(context.Background())
	if err != nil {
		t.Fatalf("expected no error counting active admins, got %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 active admin, got %d", count)
	}
}

func TestPostgresOperatorRepositoryChangeRole(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	regularOperator, err := buildTestOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid regular operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), regularOperator); err != nil {
		t.Fatalf("expected no error saving regular operator, got %v", err)
	}

	changed, err := operatorRepository.ChangeRole(
		context.Background(),
		regularOperator.ID(),
		domain.OperatorRoleAdmin,
	)
	if err != nil {
		t.Fatalf("expected no error changing role, got %v", err)
	}

	if !changed {
		t.Fatal("expected operator role to be changed")
	}

	found, err := operatorRepository.FindByID(context.Background(), regularOperator.ID())
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if found.Role() != domain.OperatorRoleAdmin {
		t.Fatalf("expected admin role, got %s", found.Role())
	}
}

func TestPostgresOperatorRepositoryDoesNotDemoteLastAdmin(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	adminOperator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), adminOperator); err != nil {
		t.Fatalf("expected no error saving admin operator, got %v", err)
	}

	changed, err := operatorRepository.ChangeRole(
		context.Background(),
		adminOperator.ID(),
		domain.OperatorRoleOperator,
	)
	if err != nil {
		t.Fatalf("expected no error changing role, got %v", err)
	}

	if changed {
		t.Fatal("expected last admin not to be demoted")
	}

	found, err := operatorRepository.FindByID(context.Background(), adminOperator.ID())
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if found.Role() != domain.OperatorRoleAdmin {
		t.Fatalf("expected admin role, got %s", found.Role())
	}
}

func TestPostgresOperatorRepositoryUpdatePasswordHash(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	operator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$old-hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	updated, err := operatorRepository.UpdatePasswordHash(
		context.Background(),
		operator.ID(),
		"$argon2id$new-hash",
	)
	if err != nil {
		t.Fatalf("expected no error updating password hash, got %v", err)
	}

	if !updated {
		t.Fatal("expected password hash to be updated")
	}

	found, err := operatorRepository.FindByID(context.Background(), operator.ID())
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if found.PasswordHash() != "$argon2id$new-hash" {
		t.Fatalf("expected updated password hash, got %s", found.PasswordHash())
	}
}

func TestPostgresOperatorRepositoryFindSummaryByID(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	operator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	summary, err := operatorRepository.FindSummaryByID(context.Background(), operator.ID())
	if err != nil {
		t.Fatalf("expected no error finding operator summary, got %v", err)
	}

	if summary.ID != operator.ID() {
		t.Fatalf("expected operator id %s, got %s", operator.ID(), summary.ID)
	}

	if summary.RegistrationID != operator.RegistrationID() {
		t.Fatalf("expected registration id %s, got %s", operator.RegistrationID(), summary.RegistrationID)
	}

	if summary.Alias != operator.Alias() {
		t.Fatalf("expected alias %s, got %s", operator.Alias(), summary.Alias)
	}

	if summary.Rank != operator.Rank() {
		t.Fatalf("expected rank %s, got %s", operator.Rank(), summary.Rank)
	}

	if summary.Role != operator.Role() {
		t.Fatalf("expected role %s, got %s", operator.Role(), summary.Role)
	}

	if summary.CreatedAt.IsZero() {
		t.Fatal("expected created_at not to be zero")
	}
}

func TestPostgresOperatorRepositoryReactivate(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)

	operator, err := buildTestOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	if err := operatorRepository.Save(context.Background(), operator); err != nil {
		t.Fatalf("expected no error saving operator, got %v", err)
	}

	deactivated, err := operatorRepository.Deactivate(context.Background(), operator.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating operator, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected operator to be deactivated")
	}

	reactivated, err := operatorRepository.Reactivate(context.Background(), operator.ID())
	if err != nil {
		t.Fatalf("expected no error reactivating operator, got %v", err)
	}

	if !reactivated {
		t.Fatal("expected operator to be reactivated")
	}

	found, err := operatorRepository.FindByID(context.Background(), operator.ID())
	if err != nil {
		t.Fatalf("expected no error finding operator, got %v", err)
	}

	if !found.Active() {
		t.Fatal("expected operator to be active")
	}
}
