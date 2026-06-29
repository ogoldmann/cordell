package domain

import "strings"

// CustodyTransactionID uniquely identifies a custody transaction.
type CustodyTransactionID string

// CustodyTransactionType identifies the business meaning of a custody transaction.
type CustodyTransactionType string

const (
	// CustodyTransactionTypeCheckout represents an asset checkout to personnel.
	CustodyTransactionTypeCheckout CustodyTransactionType = "checkout"

	// CustodyTransactionTypeReturn represents an asset return from personnel.
	CustodyTransactionTypeReturn CustodyTransactionType = "return"
)

// CustodyLine represents one asset and quantity inside a custody transaction.
type CustodyLine struct {
	assetID  AssetID
	quantity Quantity
}

// NewCustodyLine creates a custody transaction line with validated fields.
func NewCustodyLine(assetID AssetID, quantity Quantity) (CustodyLine, error) {
	if assetID == "" {
		return CustodyLine{}, ErrEmptyAssetID
	}

	if quantity.Int() <= 0 {
		return CustodyLine{}, ErrInvalidQuantity
	}

	return CustodyLine{
		assetID:  assetID,
		quantity: quantity,
	}, nil
}

// AssetID returns the asset identifier for the custody line.
func (l CustodyLine) AssetID() AssetID {
	return l.assetID
}

// Quantity returns the quantity assigned to the custody line.
func (l CustodyLine) Quantity() Quantity {
	return l.quantity
}

// CustodyTransaction represents an immutable checkout or return business event.
type CustodyTransaction struct {
	id              CustodyTransactionID
	transactionType CustodyTransactionType
	personnelID     PersonnelID
	lines           []CustodyLine
	notes           string
}

// NewCustodyTransaction creates a custody transaction with validated required fields.
func NewCustodyTransaction(
	id CustodyTransactionID,
	transactionType CustodyTransactionType,
	personnelID PersonnelID,
	lines []CustodyLine,
	notes string,
) (CustodyTransaction, error) {
	if id == "" {
		return CustodyTransaction{}, ErrEmptyTransactionID
	}

	if !transactionType.Valid() {
		return CustodyTransaction{}, ErrInvalidTransactionType
	}

	if personnelID == "" {
		return CustodyTransaction{}, ErrEmptyPersonnelID
	}

	if len(lines) == 0 {
		return CustodyTransaction{}, ErrEmptyTransactionLines
	}

	copiedLines := make([]CustodyLine, len(lines))
	copy(copiedLines, lines)

	return CustodyTransaction{
		id:              id,
		transactionType: transactionType,
		personnelID:     personnelID,
		lines:           copiedLines,
		notes:           strings.TrimSpace(notes),
	}, nil
}

// Valid reports whether the transaction type is supported by the domain.
func (t CustodyTransactionType) Valid() bool {
	return t == CustodyTransactionTypeCheckout || t == CustodyTransactionTypeReturn
}

// ID returns the custody transaction identifier.
func (t CustodyTransaction) ID() CustodyTransactionID {
	return t.id
}

// Type returns the custody transaction type.
func (t CustodyTransaction) Type() CustodyTransactionType {
	return t.transactionType
}

// PersonnelID returns the personnel identifier associated with the transaction.
func (t CustodyTransaction) PersonnelID() PersonnelID {
	return t.personnelID
}

// Lines returns a defensive copy of the custody transaction lines.
func (t CustodyTransaction) Lines() []CustodyLine {
	copiedLines := make([]CustodyLine, len(t.lines))
	copy(copiedLines, t.lines)

	return copiedLines
}

// Notes returns the optional transaction notes.
func (t CustodyTransaction) Notes() string {
	return t.notes
}
