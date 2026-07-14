package domain

import (
	"testing"
	"time"
)

func TestNewOperatorSession(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)
	expiresAt := now.Add(24 * time.Hour)

	session, err := NewOperatorSession(
		"session-1",
		"operator-1",
		"token-hash",
		expiresAt,
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	if session.OperatorID() != "operator-1" {
		t.Fatalf("expected operator-1, got %s", session.OperatorID())
	}
}

func TestOperatorSessionExpired(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	session, err := NewOperatorSession(
		"session-1",
		"operator-1",
		"token-hash",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	if session.Expired(now.Add(30 * time.Minute)) {
		t.Fatal("expected session not to be expired")
	}

	if !session.Expired(now.Add(time.Hour)) {
		t.Fatal("expected session to be expired")
	}
}
