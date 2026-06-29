package domain

import "testing"

func TestNewCustodyLine(t *testing.T) {
	quantity, err := NewQuantity(2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	line, err := NewCustodyLine("asset-1", quantity)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if line.AssetID() != "asset-1" {
		t.Fatalf("expected asset id asset-1, got %s", line.AssetID())
	}

	if line.Quantity().Int() != 2 {
		t.Fatalf("expected quantity 2, got %d", line.Quantity().Int())
	}
}

func TestNewCustodyLineRejectsEmptyAssetID(t *testing.T) {
	quantity, err := NewQuantity(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = NewCustodyLine("", quantity)
	if err != ErrEmptyAssetID {
		t.Fatalf("expected ErrEmptyAssetID, got %v", err)
	}
}

func TestCustodyTransactionTypeValid(t *testing.T) {
	validTypes := []CustodyTransactionType{
		CustodyTransactionTypeCheckout,
		CustodyTransactionTypeReturn,
	}

	for _, transactionType := range validTypes {
		if !transactionType.Valid() {
			t.Fatalf("expected transaction type %s to be valid", transactionType)
		}
	}

	if CustodyTransactionType("invalid").Valid() {
		t.Fatal("expected invalid transaction type to be invalid")
	}
}

func TestNewCustodyTransaction(t *testing.T) {
	quantity, err := NewQuantity(2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	line, err := NewCustodyLine("asset-1", quantity)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	transaction, err := NewCustodyTransaction(
		"transaction-1",
		CustodyTransactionTypeCheckout,
		"personnel-1",
		[]CustodyLine{line},
		"  Operational checkout  ",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if transaction.ID() != "transaction-1" {
		t.Fatalf("expected transaction id transaction-1, got %s", transaction.ID())
	}

	if transaction.Type() != CustodyTransactionTypeCheckout {
		t.Fatalf("expected checkout transaction, got %s", transaction.Type())
	}

	if transaction.PersonnelID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", transaction.PersonnelID())
	}

	if len(transaction.Lines()) != 1 {
		t.Fatalf("expected 1 line, got %d", len(transaction.Lines()))
	}

	if transaction.Notes() != "Operational checkout" {
		t.Fatalf("expected trimmed notes, got %s", transaction.Notes())
	}
}

func TestNewCustodyTransactionRejectsEmptyID(t *testing.T) {
	line := mustBuildCustodyLine(t)

	_, err := NewCustodyTransaction(
		"",
		CustodyTransactionTypeCheckout,
		"personnel-1",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyTransactionID {
		t.Fatalf("expected ErrEmptyTransactionID, got %v", err)
	}
}

func TestNewCustodyTransactionRejectsInvalidType(t *testing.T) {
	line := mustBuildCustodyLine(t)

	_, err := NewCustodyTransaction(
		"transaction-1",
		CustodyTransactionType("invalid"),
		"personnel-1",
		[]CustodyLine{line},
		"",
	)
	if err != ErrInvalidTransactionType {
		t.Fatalf("expected ErrInvalidTransactionType, got %v", err)
	}
}

func TestNewCustodyTransactionRejectsEmptyPersonnelID(t *testing.T) {
	line := mustBuildCustodyLine(t)

	_, err := NewCustodyTransaction(
		"transaction-1",
		CustodyTransactionTypeCheckout,
		"",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyPersonnelID {
		t.Fatalf("expected ErrEmptyPersonnelID, got %v", err)
	}
}

func TestNewCustodyTransactionRejectsEmptyLines(t *testing.T) {
	_, err := NewCustodyTransaction(
		"transaction-1",
		CustodyTransactionTypeCheckout,
		"personnel-1",
		nil,
		"",
	)
	if err != ErrEmptyTransactionLines {
		t.Fatalf("expected ErrEmptyTransactionLines, got %v", err)
	}
}

func TestCustodyTransactionLinesReturnsCopy(t *testing.T) {
	line := mustBuildCustodyLine(t)

	transaction, err := NewCustodyTransaction(
		"transaction-1",
		CustodyTransactionTypeCheckout,
		"personnel-1",
		[]CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	lines := transaction.Lines()
	lines[0] = CustodyLine{}

	if transaction.Lines()[0].AssetID() != "asset-1" {
		t.Fatal("expected transaction lines to be protected from external mutation")
	}
}

func mustBuildCustodyLine(t *testing.T) CustodyLine {
	t.Helper()

	quantity, err := NewQuantity(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	line, err := NewCustodyLine("asset-1", quantity)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	return line
}
