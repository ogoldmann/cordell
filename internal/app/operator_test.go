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
		RegistrationID: "52998224725",
		Alias:          "silva",
		Rank:           domain.RankSergeant.String(),
		Role:           domain.OperatorRoleAdmin.String(),
		Password:       "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if operator.Role() != domain.OperatorRoleAdmin {
		t.Fatalf("expected admin role, got %s", operator.Role())
	}

	if operator.ID() != "operator-1" {
		t.Fatalf("expected operator-1, got %s", operator.ID())
	}

	if operator.RegistrationID().String() != "52998224725" {
		t.Fatalf("expected registration id 52998224725, got %s", operator.RegistrationID())
	}

	if operator.Alias() != "silva" {
		t.Fatalf("expected alias silva, got %s", operator.Alias())
	}

	if operator.Rank() != domain.RankSergeant {
		t.Fatalf("expected rank sergeant, got %s", operator.Rank())
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
		RegistrationID: "52998224725",
		Alias:          "silva",
		Rank:           domain.RankSergeant.String(),
		Role:           domain.OperatorRoleAdmin.String(),
		Password:       "short",
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
		RegistrationID: "52998224725",
		Alias:          "silva",
		Rank:           domain.RankSergeant.String(),
		Role:           domain.OperatorRoleAdmin.String(),
		Password:       "   ",
	})
	if err != domain.ErrEmptyOperatorPassword {
		t.Fatalf("expected ErrEmptyOperatorPassword, got %v", err)
	}
}

func TestCreateOperatorServiceRejectsInvalidRole(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewCreateOperatorService(
		repository,
		fixedIDGenerator{id: "operator-1"},
		fakePasswordHasher{hash: "$argon2id$hash"},
	)

	_, err := service.Execute(context.Background(), CreateOperatorCommand{
		RegistrationID: "52998224725",
		Alias:          "silva",
		Rank:           domain.RankSergeant.String(),
		Role:           "root",
		Password:       "correct horse battery staple",
	})
	if err != domain.ErrInvalidOperatorRole {
		t.Fatalf("expected ErrInvalidOperatorRole, got %v", err)
	}
}
