package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/infra/postgres"
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

	personnel, err := personnelRepository.List(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error listing personnel, got %v", err)
	}

	if len(personnel) != 2 {
		t.Fatalf("expected 2 personnel records, got %d", len(personnel))
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
