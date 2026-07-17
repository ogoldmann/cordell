package app

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestListOperatorsServiceExecute(t *testing.T) {
	createdAt := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	repository := &fakeOperatorRepository{
		summaries: []ports.OperatorSummary{
			{
				ID:        "operator-1",
				Username:  "admin",
				Role:      domain.OperatorRoleAdmin,
				Active:    true,
				CreatedAt: createdAt,
			},
			{
				ID:        "operator-2",
				Username:  "clerk",
				Role:      domain.OperatorRoleOperator,
				Active:    true,
				CreatedAt: createdAt,
			},
		},
	}

	service := NewListOperatorsService(repository)

	operators, err := service.Execute(context.Background(), ListOperatorsCommand{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operators) != 2 {
		t.Fatalf("expected 2 operators, got %d", len(operators))
	}

	if operators[0].Username != "admin" {
		t.Fatalf("expected admin, got %s", operators[0].Username)
	}
}

func TestListOperatorsServiceLimitsMaximum(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewListOperatorsService(repository)

	_, err := service.Execute(context.Background(), ListOperatorsCommand{
		Limit: maxListOperatorsLimit + 1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeactivateOperatorServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	admin, err := domain.NewOperator(
		"operator-1",
		"admin",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	clerk, err := domain.NewOperator(
		"operator-2",
		"clerk",
		domain.OperatorRoleOperator,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid clerk, got %v", err)
	}

	session, err := domain.NewOperatorSession(
		"session-1",
		clerk.ID(),
		"hash:token",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	operatorRepository := &fakeOperatorRepository{
		byID: map[domain.OperatorID]domain.Operator{
			admin.ID(): admin,
			clerk.ID(): clerk,
		},
		byUsername: map[string]domain.Operator{
			admin.Username(): admin,
			clerk.Username(): clerk,
		},
	}

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{
			session.TokenHash(): session,
		},
	}

	service := NewDeactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), DeactivateOperatorCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        clerk.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	deactivatedClerk := operatorRepository.byID[clerk.ID()]
	if deactivatedClerk.Active() {
		t.Fatal("expected clerk to be inactive")
	}

	if _, ok := sessionRepository.byTokenHash[session.TokenHash()]; ok {
		t.Fatal("expected clerk sessions to be deleted")
	}
}

func TestDeactivateOperatorServiceRejectsCurrentOperator(t *testing.T) {
	admin, err := domain.NewOperator(
		"operator-1",
		"admin",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	operatorRepository := &fakeOperatorRepository{
		byID: map[domain.OperatorID]domain.Operator{
			admin.ID(): admin,
		},
		byUsername: map[string]domain.Operator{
			admin.Username(): admin,
		},
	}

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewDeactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), DeactivateOperatorCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        admin.ID(),
	})
	if err != domain.ErrCannotDeactivateCurrentOperator {
		t.Fatalf("expected ErrCannotDeactivateCurrentOperator, got %v", err)
	}
}

func TestDeactivateOperatorServiceRejectsLastAdmin(t *testing.T) {
	currentAdmin, err := domain.NewOperator(
		"operator-1",
		"current-admin",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid current admin, got %v", err)
	}

	targetAdmin, err := domain.NewOperator(
		"operator-2",
		"target-admin",
		domain.OperatorRoleAdmin,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid target admin, got %v", err)
	}

	inactiveCurrentAdmin, err := domain.ReconstituteOperator(
		currentAdmin.ID(),
		currentAdmin.Username(),
		currentAdmin.Role(),
		currentAdmin.PasswordHash(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive current admin, got %v", err)
	}

	operatorRepository := &fakeOperatorRepository{
		byID: map[domain.OperatorID]domain.Operator{
			inactiveCurrentAdmin.ID(): inactiveCurrentAdmin,
			targetAdmin.ID():          targetAdmin,
		},
		byUsername: map[string]domain.Operator{
			inactiveCurrentAdmin.Username(): inactiveCurrentAdmin,
			targetAdmin.Username():          targetAdmin,
		},
	}

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewDeactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), DeactivateOperatorCommand{
		CurrentOperatorID: inactiveCurrentAdmin.ID(),
		OperatorID:        targetAdmin.ID(),
	})
	if err != domain.ErrCannotDeactivateLastAdmin {
		t.Fatalf("expected ErrCannotDeactivateLastAdmin, got %v", err)
	}
}
