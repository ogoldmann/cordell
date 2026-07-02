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
