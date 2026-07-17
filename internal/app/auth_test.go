package app

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
)

func TestAuthenticateOperatorServiceExecute(t *testing.T) {
	operator, err := domain.NewOperator("operator-1", "admin", domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	repository := &fakeOperatorRepository{
		byID: map[domain.OperatorID]domain.Operator{
			operator.ID(): operator,
		},
		byUsername: map[string]domain.Operator{
			operator.Username(): operator,
		},
	}

	service := NewAuthenticateOperatorService(
		repository,
		fakePasswordHasher{hash: "$argon2id$hash"},
	)

	result, err := service.Execute(context.Background(), AuthenticateOperatorCommand{
		Username: "ADMIN",
		Password: "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID() != "operator-1" {
		t.Fatalf("expected operator-1, got %s", result.ID())
	}
}

func TestAuthenticateOperatorServiceRejectsInvalidCredentials(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewAuthenticateOperatorService(
		repository,
		fakePasswordHasher{hash: "$argon2id$hash"},
	)

	_, err := service.Execute(context.Background(), AuthenticateOperatorCommand{
		Username: "missing",
		Password: "correct horse battery staple",
	})
	if err != domain.ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestCreateOperatorSessionServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	repository := &fakeOperatorSessionRepository{}

	service := NewCreateOperatorSessionService(
		repository,
		fixedIDGenerator{id: "session-1"},
		fixedSessionTokenGenerator{token: "raw-token"},
		fixedSessionTokenGenerator{token: "csrf-token"},
		plainSessionTokenHasher{},
	)
	result, err := service.Execute(context.Background(), CreateOperatorSessionCommand{
		OperatorID: "operator-1",
		Now:        now,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Token != "raw-token" {
		t.Fatalf("expected raw-token, got %s", result.Token)
	}

	if result.Session.CSRFToken() != "csrf-token" {
		t.Fatalf("expected csrf-token, got %s", result.Session.CSRFToken())
	}

	if result.Session.TokenHash() != "hash:raw-token" {
		t.Fatalf("expected hash:raw-token, got %s", result.Session.TokenHash())
	}
}

func TestGetOperatorBySessionTokenServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	operator, err := domain.NewOperator("operator-1", "admin", domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	session, err := domain.NewOperatorSession(
		"session-1",
		operator.ID(),
		"hash:raw-token",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{
			session.TokenHash(): session,
		},
	}

	operatorRepository := &fakeOperatorRepository{
		byID: map[domain.OperatorID]domain.Operator{
			operator.ID(): operator,
		},
		byUsername: map[string]domain.Operator{
			operator.Username(): operator,
		},
	}

	service := NewGetOperatorBySessionTokenService(
		sessionRepository,
		operatorRepository,
		plainSessionTokenHasher{},
	)

	result, err := service.Execute(context.Background(), GetOperatorBySessionTokenCommand{
		Token: "raw-token",
		Now:   now,
	})
	if result.Operator.ID() != operator.ID() {
		t.Fatalf("expected operator %s, got %s", operator.ID(), result.Operator.ID())
	}

	if result.Session.ID() != session.ID() {
		t.Fatalf("expected session %s, got %s", session.ID(), result.Session.ID())
	}
}

func TestDeleteOperatorSessionServiceExecute(t *testing.T) {
	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewDeleteOperatorSessionService(
		sessionRepository,
		plainSessionTokenHasher{},
	)

	if err := service.Execute(context.Background(), DeleteOperatorSessionCommand{
		Token: "raw-token",
	}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteExpiredOperatorSessionsServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	activeSession, err := domain.NewOperatorSession(
		"session-1",
		"operator-1",
		"hash:active-token",
		"csrf-token-1",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid active session, got %v", err)
	}

	expiredSession, err := domain.NewOperatorSession(
		"session-2",
		"operator-1",
		"hash:expired-token",
		"csrf-token-2",
		now.Add(-time.Hour),
		now.Add(-2*time.Hour),
	)
	if err != nil {
		t.Fatalf("expected valid expired session, got %v", err)
	}

	repository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{
			activeSession.TokenHash():  activeSession,
			expiredSession.TokenHash(): expiredSession,
		},
	}

	service := NewDeleteExpiredOperatorSessionsService(repository)

	if err := service.Execute(context.Background(), DeleteExpiredOperatorSessionsCommand{
		Now: now,
	}); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if _, ok := repository.byTokenHash[activeSession.TokenHash()]; !ok {
		t.Fatal("expected active session to remain")
	}

	if _, ok := repository.byTokenHash[expiredSession.TokenHash()]; ok {
		t.Fatal("expected expired session to be deleted")
	}
}
