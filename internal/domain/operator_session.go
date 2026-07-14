package domain

import (
	"strings"
	"time"
)

// OperatorSessionID identifies an operator session.
type OperatorSessionID string

// OperatorSession represents a server-side operator session.
type OperatorSession struct {
	id         OperatorSessionID
	operatorID OperatorID
	tokenHash  string
	expiresAt  time.Time
	createdAt  time.Time
}

// NewOperatorSession creates an OperatorSession.
func NewOperatorSession(
	id OperatorSessionID,
	operatorID OperatorID,
	tokenHash string,
	expiresAt time.Time,
	createdAt time.Time,
) (OperatorSession, error) {
	if strings.TrimSpace(string(id)) == "" {
		return OperatorSession{}, ErrEmptyOperatorSessionID
	}

	if strings.TrimSpace(string(operatorID)) == "" {
		return OperatorSession{}, ErrEmptyOperatorID
	}

	tokenHash = strings.TrimSpace(tokenHash)
	if tokenHash == "" {
		return OperatorSession{}, ErrEmptyOperatorSessionTokenHash
	}

	return OperatorSession{
		id:         id,
		operatorID: operatorID,
		tokenHash:  tokenHash,
		expiresAt:  expiresAt.UTC(),
		createdAt:  createdAt.UTC(),
	}, nil
}

// ID returns the session identifier.
func (s OperatorSession) ID() OperatorSessionID {
	return s.id
}

// OperatorID returns the authenticated operator identifier.
func (s OperatorSession) OperatorID() OperatorID {
	return s.operatorID
}

// TokenHash returns the session token hash.
func (s OperatorSession) TokenHash() string {
	return s.tokenHash
}

// ExpiresAt returns the session expiration time.
func (s OperatorSession) ExpiresAt() time.Time {
	return s.expiresAt
}

// CreatedAt returns the session creation time.
func (s OperatorSession) CreatedAt() time.Time {
	return s.createdAt
}

// Expired reports whether the session is expired at the given time.
func (s OperatorSession) Expired(now time.Time) bool {
	return !now.UTC().Before(s.expiresAt)
}
