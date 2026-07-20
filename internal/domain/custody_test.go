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
		"operator-1",
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

	if transaction.OperatorID() != "operator-1" {
		t.Fatalf("expected operator id operator-1, got %s", transaction.OperatorID())
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
		"operator-1",
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
		"operator-1",
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
		"operator-1",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyPersonnelID {
		t.Fatalf("expected ErrEmptyPersonnelID, got %v", err)
	}
}

func TestNewCustodyTransactionRejectsEmptyOperatorID(t *testing.T) {
	line, err := NewCustodyLine("asset-1", Quantity(1))
	if err != nil {
		t.Fatalf("expected valid line, got %v", err)
	}

	_, err = NewCustodyTransaction(
		"transaction-1",
		CustodyTransactionTypeCheckout,
		"personnel-1",
		"",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyOperatorID {
		t.Fatalf("expected ErrEmptyOperatorID, got %v", err)
	}
}

func TestNewCustodyTransactionRejectsEmptyLines(t *testing.T) {
	_, err := NewCustodyTransaction(
		"transaction-1",
		CustodyTransactionTypeCheckout,
		"personnel-1",
		"operator-1",
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
		"operator-1",
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

func TestNewCustodyCorrection(t *testing.T) {
	line := mustBuildCustodyLine(t)

	correction, err := NewCustodyCorrection(
		"correction-1",
		"transaction-1",
		"operator-1",
		"personnel-2",
		[]CustodyLine{line},
		"  Corrected notes  ",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if correction.ID() != "correction-1" {
		t.Fatalf("expected correction id correction-1, got %s", correction.ID())
	}

	if correction.CorrectedTransactionID() != "transaction-1" {
		t.Fatalf("expected corrected transaction id transaction-1, got %s", correction.CorrectedTransactionID())
	}

	if correction.OperatorID() != "operator-1" {
		t.Fatalf("expected operator id operator-1, got %s", correction.OperatorID())
	}

	if correction.CorrectedPersonnelID() != "personnel-2" {
		t.Fatalf("expected corrected personnel id personnel-2, got %s", correction.CorrectedPersonnelID())
	}

	if correction.CorrectedNotes() != "Corrected notes" {
		t.Fatalf("expected trimmed corrected notes, got %s", correction.CorrectedNotes())
	}
}

func TestNewCustodyCorrectionRejectsEmptyID(t *testing.T) {
	line := mustBuildCustodyLine(t)

	_, err := NewCustodyCorrection(
		"",
		"transaction-1",
		"operator-1",
		"personnel-2",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyCustodyCorrectionID {
		t.Fatalf("expected ErrEmptyCustodyCorrectionID, got %v", err)
	}
}

func TestNewCustodyCorrectionRejectsEmptyTransactionID(t *testing.T) {
	line := mustBuildCustodyLine(t)

	_, err := NewCustodyCorrection(
		"correction-1",
		"",
		"operator-1",
		"personnel-2",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyTransactionID {
		t.Fatalf("expected ErrEmptyTransactionID, got %v", err)
	}
}

func TestNewCustodyCorrectionRejectsEmptyOperatorID(t *testing.T) {
	line := mustBuildCustodyLine(t)

	_, err := NewCustodyCorrection(
		"correction-1",
		"transaction-1",
		"",
		"personnel-2",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyOperatorID {
		t.Fatalf("expected ErrEmptyOperatorID, got %v", err)
	}
}

func TestNewCustodyCorrectionRejectsEmptyPersonnelID(t *testing.T) {
	line := mustBuildCustodyLine(t)

	_, err := NewCustodyCorrection(
		"correction-1",
		"transaction-1",
		"operator-1",
		"",
		[]CustodyLine{line},
		"",
	)
	if err != ErrEmptyPersonnelID {
		t.Fatalf("expected ErrEmptyPersonnelID, got %v", err)
	}
}

func TestNewCustodyCorrectionRejectsEmptyLines(t *testing.T) {
	_, err := NewCustodyCorrection(
		"correction-1",
		"transaction-1",
		"operator-1",
		"personnel-2",
		nil,
		"",
	)
	if err != ErrEmptyTransactionLines {
		t.Fatalf("expected ErrEmptyTransactionLines, got %v", err)
	}
}

func TestCustodyCorrectionLinesReturnsCopy(t *testing.T) {
	line := mustBuildCustodyLine(t)

	correction, err := NewCustodyCorrection(
		"correction-1",
		"transaction-1",
		"operator-1",
		"personnel-2",
		[]CustodyLine{line},
		"",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	lines := correction.Lines()
	lines[0] = CustodyLine{}

	if correction.Lines()[0].AssetID() != "asset-1" {
		t.Fatal("expected correction lines to be protected from external mutation")
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
