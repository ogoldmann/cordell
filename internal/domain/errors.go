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

)