package domain

import "strings"

// OperatorID identifies an operator account.
type OperatorID string

// Operator represents a system operator account.
type Operator struct {
	id             OperatorID
	registrationID RegistrationID
	alias          string
	rank           Rank
	role           OperatorRole
	passwordHash   string
	active         bool
}

// NewOperator creates an active Operator after validating required fields.
func NewOperator(
	id OperatorID,
	registrationID RegistrationID,
	alias string,
	rank Rank,
	role OperatorRole,
	passwordHash string,
) (Operator, error) {
	return buildOperator(id, registrationID, alias, rank, role, passwordHash, true)
}

// ReconstituteOperator rebuilds an Operator from persisted state.
func ReconstituteOperator(
	id OperatorID,
	registrationID RegistrationID,
	alias string,
	rank Rank,
	role OperatorRole,
	passwordHash string,
	active bool,
) (Operator, error) {
	return buildOperator(id, registrationID, alias, rank, role, passwordHash, active)
}

func buildOperator(
	id OperatorID,
	registrationID RegistrationID,
	alias string,
	rank Rank,
	role OperatorRole,
	passwordHash string,
	active bool,
) (Operator, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Operator{}, ErrEmptyOperatorID
	}

	validRegistrationID, err := NewRegistrationID(registrationID.String())
	if err != nil {
		return Operator{}, err
	}

	alias = strings.TrimSpace(alias)
	if alias == "" {
		return Operator{}, ErrEmptyOperatorAlias
	}

	rank = Rank(strings.ToLower(strings.TrimSpace(rank.String())))
	if !IsValidRank(rank) {
		return Operator{}, ErrInvalidOperatorRank
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
		id:             id,
		registrationID: validRegistrationID,
		alias:          alias,
		rank:           rank,
		role:           validRole,
		passwordHash:   passwordHash,
		active:         active,
	}, nil
}

// ID returns the operator identifier.
func (o Operator) ID() OperatorID {
	return o.id
}

// RegistrationID returns the operator registration identifier.
func (o Operator) RegistrationID() RegistrationID {
	return o.registrationID
}

// Alias returns the operator operational alias.
func (o Operator) Alias() string {
	return o.alias
}

// Rank returns the operator rank.
func (o Operator) Rank() Rank {
	return o.rank
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
