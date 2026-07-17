package domain

import "errors"

var (
	// ErrEmptyPersonnelID is returned when a personnel identifier is required but empty.
	ErrEmptyPersonnelID = errors.New("personnel id cannot be empty")
	// ErrEmptyPersonnelName is returned when a personnel name is required but empty.
	ErrEmptyPersonnelName = errors.New("personnel name cannot be empty")
	// ErrEmptyAssetName is returned when an asset name is required but empty.
	ErrEmptyAssetName = errors.New("asset name cannot be empty")
	// ErrEmptyAssetID is returned when an asset identifier is required but empty.
	ErrEmptyAssetID = errors.New("asset id cannot be empty")
	// ErrEmptyTransactionID is returned when a custody transaction identifier is required but empty.
	ErrEmptyTransactionID = errors.New("custody transaction id cannot be empty")
	// ErrInvalidTransactionType is returned when a custody transaction type is unsupported.
	ErrInvalidTransactionType = errors.New("invalid custody transaction type")
	// ErrEmptyTransactionLines is returned when a custody transaction has no lines.
	ErrEmptyTransactionLines = errors.New("custody transaction must have at least one line")
	// ErrInvalidQuantity is returned when a quantity is zero or negative.
	ErrInvalidQuantity = errors.New("quantity must be greater than zero")
	// ErrInsufficientCustodyBalance is returned when a return exceeds the current custody quantity.
	ErrInsufficientCustodyBalance = errors.New("insufficient custody balance")

	// ErrEmptyPersonnelAlias is returned when a personnel alias is required but empty.
	ErrEmptyPersonnelAlias = errors.New("empty personnel alias")
	// ErrEmptyRegistrationID is returned when a personnel registration identifier is required but empty.
	ErrEmptyRegistrationID = errors.New("empty registration id")
	// ErrInvalidRegistrationID is returned when a personnel registration identifier has an invalid format.
	ErrInvalidRegistrationID = errors.New("invalid registration id")
	// ErrDuplicateRegistrationID is returned when a personnel registration identifier is already in use.
	ErrDuplicateRegistrationID = errors.New("duplicate registration id")
	// ErrInvalidPersonnelRank is returned when a personnel rank is unsupported.
	ErrInvalidPersonnelRank = errors.New("invalid personnel rank")
	// ErrInvalidPersonnelSection is returned when a personnel section is unsupported.
	ErrInvalidPersonnelSection = errors.New("invalid personnel section")
	// ErrInvalidOrganizationUnit is returned when a personnel organization unit is unsupported.
	ErrInvalidOrganizationUnit = errors.New("invalid organization unit")

	// ErrEmptyOperatorID is returned when an operator identifier is required but empty.
	ErrEmptyOperatorID           = errors.New("empty operator id")
	ErrEmptyOperatorPasswordHash = errors.New("empty operator password hash")
	ErrEmptyOperatorPassword     = errors.New("empty operator password")
	// ErrWeakOperatorPassword is returned when an operator password does not satisfy password policy.
	ErrWeakOperatorPassword = errors.New("weak operator password")
	// ErrEmptyOperatorAlias is returned when an operator alias is required but empty.
	ErrEmptyOperatorAlias = errors.New("empty operator alias")
	// ErrInvalidOperatorRank is returned when an operator rank is unsupported.
	ErrInvalidOperatorRank = errors.New("invalid operator rank")
	ErrEmptyOperatorRole   = errors.New("empty operator role")
	ErrInvalidOperatorRole = errors.New("invalid operator role")

	// ErrEmptyOperatorSessionID is returned when an operator session identifier is required but empty.
	ErrEmptyOperatorSessionID = errors.New("empty operator session id")
	// ErrEmptyOperatorSessionTokenHash is returned when an operator session token hash is required but empty.
	ErrEmptyOperatorSessionTokenHash = errors.New("empty operator session token hash")
	// ErrExpiredOperatorSession is returned when an operator session is no longer valid.
	ErrExpiredOperatorSession = errors.New("expired operator session")
	// ErrInvalidCredentials is returned when authentication credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrEmptyOperatorSessionCSRFToken is returned when an operator session CSRF token is required but empty.
	ErrEmptyOperatorSessionCSRFToken = errors.New("empty operator session csrf token")
	// ErrInvalidCSRFToken is returned when a CSRF token is invalid.
	ErrInvalidCSRFToken = errors.New("invalid csrf token")

	ErrCannotDeactivateCurrentOperator = errors.New("cannot deactivate current operator")
	ErrCannotDeactivateLastAdmin       = errors.New("cannot deactivate last admin")

	ErrCannotChangeCurrentOperatorRole = errors.New("cannot change current operator role")
	ErrCannotDemoteLastAdmin           = errors.New("cannot demote last admin")

	ErrCannotResetCurrentOperatorPassword = errors.New("cannot reset current operator password")
)
