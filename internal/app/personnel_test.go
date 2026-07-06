package app

import (
	"context"
	"testing"

	"cordell/internal/domain"
)

func TestCreatePersonnelServiceExecute(t *testing.T) {
	repository := &fakePersonnelRepository{}
	idGenerator := fixedIDGenerator{id: "personnel-1"}

	service := NewCreatePersonnelService(repository, idGenerator)

	personnel, err := service.Execute(context.Background(), CreatePersonnelCommand{
		FullName: "  John Doe  ",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if personnel.ID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", personnel.ID())
	}

	if personnel.FullName() != "John Doe" {
		t.Fatalf("expected trimmed full name John Doe, got %s", personnel.FullName())
	}

	if len(repository.saved) != 1 {
		t.Fatalf("expected 1 saved personnel, got %d", len(repository.saved))
	}
}

func TestCreatePersonnelServiceRejectsInvalidPersonnel(t *testing.T) {
	repository := &fakePersonnelRepository{}
	idGenerator := fixedIDGenerator{id: "personnel-1"}

	service := NewCreatePersonnelService(repository, idGenerator)

	_, err := service.Execute(context.Background(), CreatePersonnelCommand{
		FullName: "   ",
	})
	if err != domain.ErrEmptyPersonnelName {
		t.Fatalf("expected ErrEmptyPersonnelName, got %v", err)
	}

	if len(repository.saved) != 0 {
		t.Fatalf("expected no saved personnel, got %d", len(repository.saved))
	}
}

func TestGetPersonnelServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	service := NewGetPersonnelService(repository)

	found, err := service.Execute(context.Background(), GetPersonnelCommand{
		ID: "personnel-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", found.ID())
	}

	if found.FullName() != "John Doe" {
		t.Fatalf("expected personnel name John Doe, got %s", found.FullName())
	}
}

func TestGetPersonnelServiceRejectsEmptyID(t *testing.T) {
	repository := &fakePersonnelRepository{}
	service := NewGetPersonnelService(repository)

	_, err := service.Execute(context.Background(), GetPersonnelCommand{
		ID: "",
	})
	if err != domain.ErrEmptyPersonnelID {
		t.Fatalf("expected ErrEmptyPersonnelID, got %v", err)
	}
}
