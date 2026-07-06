package postgres_test

import (
	"context"
	"testing"

	"cordell/internal/app"
	"cordell/internal/infra/postgres"
)

func TestPostgresPersonnelRepositoryCreateAndFind(t *testing.T) {
	pool := openTestPool(t)
	queries := newTestQueries(pool)

	personnelRepository := postgres.NewPersonnelRepository(queries)
	idGenerator := fixedIDGenerator{id: "personnel-1"}

	service := app.NewCreatePersonnelService(personnelRepository, idGenerator)

	created, err := service.Execute(context.Background(), app.CreatePersonnelCommand{
		FullName: "  John Doe  ",
	})
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
