package domain

import (
	"strings"
	"unicode"
)

// OperatorID identifies an operator account.
type OperatorID string

// Operator represents a system operator account.
type Operator struct {
	id           OperatorID
	username     string
	role         OperatorRole
	passwordHash string
	active       bool
}

// NewOperator creates an active Operator after validating required fields.
func NewOperator(id OperatorID, username string, role OperatorRole, passwordHash string) (Operator, error) {
	return buildOperator(id, username, role, passwordHash, true)
}

// ReconstituteOperator rebuilds an Operator from persisted state.
func ReconstituteOperator(
	id OperatorID,
	username string,
	role OperatorRole,
	passwordHash string,
	active bool,
) (Operator, error) {
	return buildOperator(id, username, role, passwordHash, active)
}

func buildOperator(
	id OperatorID,
	username string,
	role OperatorRole,
	passwordHash string,
	active bool,
) (Operator, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Operator{}, ErrEmptyOperatorID
	}

	username = NormalizeOperatorUsername(username)
	if username == "" {
		return Operator{}, ErrEmptyOperatorUsername
	}

	if !isValidOperatorUsername(username) {
		return Operator{}, ErrInvalidOperatorUsername
	}

	validRole, err := NewOperatorRole(role.String())
	if err != nil {
		return Operator{}, err
	}

	passwordHash = strings.TrimSpace(passwordHash)
	if passwordHash == "" {
		return Operator{}, ErrEmptyOperatorPasswordHash
	}

	return Operator{
		id:           id,
		username:     username,
		role:         validRole,
		passwordHash: passwordHash,
		active:       active,
	}, nil
}

// NormalizeOperatorUsername normalizes an operator username for storage and lookup.
func NormalizeOperatorUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func isValidOperatorUsername(username string) bool {
	if len(username) < 3 || len(username) > 64 {
		return false
	}

	for _, char := range username {
		if unicode.IsLower(char) || unicode.IsDigit(char) {
			continue
		}

		switch char {
		case '.', '_', '-':
			continue
		default:
			return false
		}
	}

	return true
}

// ID returns the operator identifier.
func (o Operator) ID() OperatorID {
	return o.id
}

// Username returns the operator username.
func (o Operator) Username() string {
	return o.username
}

// Role returns the operator role.
func (o Operator) Role() OperatorRole {
	return o.role
}

// PasswordHash returns the operator password hash.
func (o Operator) PasswordHash() string {
	return o.passwordHash
}

// Active reports whether the operator account is active.
func (o Operator) Active() bool {
	return o.active
}
