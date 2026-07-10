package app

import (
	"context"
	"testing"

	"cordell/internal/domain"
)

func TestCreateOperatorServiceExecute(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewCreateOperatorService(
		repository,
		fixedIDGenerator{id: "operator-1"},
		fakePasswordHasher{hash: "$argon2id$hash"},
	)

	operator, err := service.Execute(context.Background(), CreateOperatorCommand{
		Username: "Admin.User",
		Password: "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if operator.ID() != "operator-1" {
		t.Fatalf("expected operator-1, got %s", operator.ID())
	}

	if operator.Username() != "admin.user" {
		t.Fatalf("expected admin.user, got %s", operator.Username())
	}

	if operator.PasswordHash() != "$argon2id$hash" {
		t.Fatalf("expected stored password hash, got %s", operator.PasswordHash())
	}
}

func TestCreateOperatorServiceRejectsWeakPassword(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewCreateOperatorService(
		repository,
		fixedIDGenerator{id: "operator-1"},
		fakePasswordHasher{hash: "$argon2id$hash"},
	)

	_, err := service.Execute(context.Background(), CreateOperatorCommand{
		Username: "admin",
		Password: "short",
	})
	if err != domain.ErrWeakOperatorPassword {
		t.Fatalf("expected ErrWeakOperatorPassword, got %v", err)
	}
}

func TestCreateOperatorServiceRejectsEmptyPassword(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewCreateOperatorService(
		repository,
		fixedIDGenerator{id: "operator-1"},
		fakePasswordHasher{hash: "$argon2id$hash"},
	)

	_, err := service.Execute(context.Background(), CreateOperatorCommand{
		Username: "admin",
		Password: "   ",
	})
	if err != domain.ErrEmptyOperatorPassword {
		t.Fatalf("expected ErrEmptyOperatorPassword, got %v", err)
	}
}
