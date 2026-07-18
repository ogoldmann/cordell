package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
	"cordell/internal/ports"
)

func TestPostgresPersonnelRepositoryCreateAndFind(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	idGenerator := fixedIDGenerator{id: "personnel-1"}

	service := app.NewCreatePersonnelService(personnelRepository, idGenerator)

	created, err := service.Execute(context.Background(), validCreatePersonnelCommand("  John Doe  ", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found, err := personnelRepository.FindByID(context.Background(), created.ID())
	if err != nil {
		t.Fatalf("expected no error finding personnel, got %v", err)
	}

	if found.ID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", found.ID())
	}

	if found.FullName() != "John Doe" {
		t.Fatalf("expected full name John Doe, got %s", found.FullName())
	}

	if !found.Active() {
		t.Fatal("expected personnel to be active")
	}
}

func TestPostgresPersonnelRepositoryList(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)

	firstService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	secondService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-2"},
	)

	_, err := firstService.Execute(context.Background(), validCreatePersonnelCommand("John Doe", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error creating first personnel, got %v", err)
	}

	_, err = secondService.Execute(context.Background(), validCreatePersonnelCommand("Jane Doe", "Jane", "11144477735"))
	if err != nil {
		t.Fatalf("expected no error creating second personnel, got %v", err)
	}

	personnel, err := personnelRepository.List(context.Background(), 10, ports.RecordStatusFilterActive)
	if err != nil {
		t.Fatalf("expected no error listing personnel, got %v", err)
	}

	if len(personnel) != 2 {
		t.Fatalf("expected 2 personnel records, got %d", len(personnel))
	}
}

func TestPostgresPersonnelRepositoryListFiltersByStatus(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)
	repository := postgres.NewPersonnelRepository(queries)

	activePersonnel := mustNewTestPersonnel(t, "personnel-1", "John Active", "active", domain.RankSergeant, "52998224725")
	inactiveBase := mustNewTestPersonnel(t, "personnel-2", "John Inactive", "inactive", domain.RankCorporal, "93541134780")

	inactivePersonnel, err := domain.ReconstitutePersonnel(
		inactiveBase.ID(),
		inactiveBase.FullName(),
		inactiveBase.Alias(),
		inactiveBase.Rank(),
		inactiveBase.RegistrationID(),
		inactiveBase.Section(),
		inactiveBase.OrganizationUnit(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive personnel, got %v", err)
	}

	if err := repository.Save(context.Background(), activePersonnel); err != nil {
		t.Fatalf("expected no error saving active personnel, got %v", err)
	}

	if err := repository.Save(context.Background(), inactivePersonnel); err != nil {
		t.Fatalf("expected no error saving inactive personnel, got %v", err)
	}

	active, err := repository.List(context.Background(), 10, ports.RecordStatusFilterActive)
	if err != nil {
		t.Fatalf("expected no error listing active personnel, got %v", err)
	}

	if len(active) != 1 {
		t.Fatalf("expected 1 active personnel, got %d", len(active))
	}

	if !active[0].Active() {
		t.Fatal("expected listed personnel to be active")
	}

	inactive, err := repository.List(context.Background(), 10, ports.RecordStatusFilterInactive)
	if err != nil {
		t.Fatalf("expected no error listing inactive personnel, got %v", err)
	}

	if len(inactive) != 1 {
		t.Fatalf("expected 1 inactive personnel, got %d", len(inactive))
	}

	if inactive[0].Active() {
		t.Fatal("expected listed personnel to be inactive")
	}

	all, err := repository.List(context.Background(), 10, ports.RecordStatusFilterAll)
	if err != nil {
		t.Fatalf("expected no error listing all personnel, got %v", err)
	}

	if len(all) != 2 {
		t.Fatalf("expected 2 personnel, got %d", len(all))
	}
}

func TestPostgresPersonnelRepositoryDeactivateAndReactivate(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)
	repository := postgres.NewPersonnelRepository(queries)

	personnel := mustNewTestPersonnel(t, "personnel-1", "John Doe", "doe", domain.RankSergeant, "52998224725")

	if err := repository.Save(context.Background(), personnel); err != nil {
		t.Fatalf("expected no error saving personnel, got %v", err)
	}

	deactivated, err := repository.Deactivate(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error deactivating personnel, got %v", err)
	}

	if !deactivated {
		t.Fatal("expected personnel to be deactivated")
	}

	updated, err := repository.FindByID(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error finding personnel, got %v", err)
	}

	if updated.Active() {
		t.Fatal("expected personnel to be inactive")
	}

	reactivated, err := repository.Reactivate(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error reactivating personnel, got %v", err)
	}

	if !reactivated {
		t.Fatal("expected personnel to be reactivated")
	}

	updated, err = repository.FindByID(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected no error finding personnel, got %v", err)
	}

	if !updated.Active() {
		t.Fatal("expected personnel to be active")
	}
}

func TestPostgresPersonnelRepositoryRejectsDuplicateRegistrationID(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)

	firstService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	secondService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-2"},
	)

	_, err := firstService.Execute(
		context.Background(),
		validCreatePersonnelCommand("John Doe", "Doe", "52998224725"),
	)
	if err != nil {
		t.Fatalf("expected no error creating first personnel, got %v", err)
	}

	_, err = secondService.Execute(
		context.Background(),
		validCreatePersonnelCommand("Jane Doe", "Jane", "52998224725"),
	)
	if err != domain.ErrDuplicateRegistrationID {
		t.Fatalf("expected ErrDuplicateRegistrationID, got %v", err)
	}
}

func TestPostgresPersonnelRepositorySearch(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)

	firstService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	secondService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-2"},
	)

	_, err := firstService.Execute(
		context.Background(),
		validCreatePersonnelCommand("John Doe", "Doe", "52998224725"),
	)
	if err != nil {
		t.Fatalf("expected no error creating first personnel, got %v", err)
	}

	_, err = secondService.Execute(
		context.Background(),
		validCreatePersonnelCommand("Jane Smith", "Smith", "11144477735"),
	)
	if err != nil {
		t.Fatalf("expected no error creating second personnel, got %v", err)
	}

	personnel, err := personnelRepository.Search(context.Background(), "smith", 10)
	if err != nil {
		t.Fatalf("expected no error searching personnel, got %v", err)
	}

	if len(personnel) != 1 {
		t.Fatalf("expected 1 personnel record, got %d", len(personnel))
	}

	if personnel[0].Alias() != "Smith" {
		t.Fatalf("expected alias Smith, got %s", personnel[0].Alias())
	}
}

func TestPostgresPersonnelRepositorySearchByFormattedRegistrationID(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)

	service := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)

	_, err := service.Execute(
		context.Background(),
		validCreatePersonnelCommand("John Doe", "Doe", "52998224725"),
	)
	if err != nil {
		t.Fatalf("expected no error creating personnel, got %v", err)
	}

	personnel, err := personnelRepository.Search(context.Background(), "529.982", 10)
	if err != nil {
		t.Fatalf("expected no error searching personnel by registration id, got %v", err)
	}

	if len(personnel) != 1 {
		t.Fatalf("expected 1 personnel record, got %d", len(personnel))
	}

	if personnel[0].RegistrationID().String() != "52998224725" {
		t.Fatalf("expected registration id 52998224725, got %s", personnel[0].RegistrationID().String())
	}
}

func TestPostgresPersonnelRepositorySearchByCombinedTerms(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)

	firstService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-1"},
	)
	secondService := app.NewCreatePersonnelService(
		personnelRepository,
		fixedIDGenerator{id: "personnel-2"},
	)

	_, err := firstService.Execute(context.Background(), app.CreatePersonnelCommand{
		FullName:         "John Alpha",
		Alias:            "Doe",
		Rank:             domain.PersonnelRankSergeant,
		RegistrationID:   "52998224725",
		Section:          domain.PersonnelSectionOperations,
		OrganizationUnit: domain.OrganizationUnitDefault,
	})
	if err != nil {
		t.Fatalf("expected no error creating first personnel, got %v", err)
	}

	_, err = secondService.Execute(context.Background(), app.CreatePersonnelCommand{
		FullName:         "Jane Doe",
		Alias:            "Smith",
		Rank:             domain.PersonnelRankCorporal,
		RegistrationID:   "11144477735",
		Section:          domain.PersonnelSectionLogistics,
		OrganizationUnit: domain.OrganizationUnitDefault,
	})
	if err != nil {
		t.Fatalf("expected no error creating second personnel, got %v", err)
	}

	personnel, err := personnelRepository.Search(context.Background(), "sergeant doe", 10)
	if err != nil {
		t.Fatalf("expected no error searching personnel, got %v", err)
	}

	if len(personnel) != 1 {
		t.Fatalf("expected 1 personnel record, got %d", len(personnel))
	}

	if personnel[0].ID() != "personnel-1" {
		t.Fatalf("expected personnel-1, got %s", personnel[0].ID())
	}
}
