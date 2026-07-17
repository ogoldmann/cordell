package domain

import "strings"

// OperatorRole represents an operator authorization role.
type OperatorRole string

const (
	// OperatorRoleAdmin can administer the system.
	OperatorRoleAdmin OperatorRole = "admin"

	// OperatorRoleOperator can use regular custody workflows.
	OperatorRoleOperator OperatorRole = "operator"
)

// NewOperatorRole validates and normalizes an operator role.
func NewOperatorRole(value string) (OperatorRole, error) {
	normalized := OperatorRole(strings.ToLower(strings.TrimSpace(value)))

	if normalized == "" {
		return "", ErrEmptyOperatorRole
	}

	switch normalized {
	case OperatorRoleAdmin, OperatorRoleOperator:
		return normalized, nil
	default:
		return "", ErrInvalidOperatorRole
	}
}

// String returns the role string value.
func (r OperatorRole) String() string {
	return string(r)
}

// Label returns a human-readable role label.
func (r OperatorRole) Label() string {
	switch r {
	case OperatorRoleAdmin:
		return "Admin"
	case OperatorRoleOperator:
		return "Operator"
	default:
		return string(r)
	}
}

// CanManageOperators reports whether the role can manage operator accounts.
func (r OperatorRole) CanManageOperators() bool {
	return r == OperatorRoleAdmin
}

// OperatorRoleOptions returns the available operator roles.
func OperatorRoleOptions() []OperatorRole {
	return []OperatorRole{
		OperatorRoleAdmin,
		OperatorRoleOperator,
	}
}
