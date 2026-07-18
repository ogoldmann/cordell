package web

import (
	"time"

	"cordell/internal/app"
)

type custodyReceiptView struct {
	ID                        string
	Type                      string
	TypeLabel                 string
	PersonnelID               string
	PersonnelFullName         string
	PersonnelAlias            string
	PersonnelRankLabel        string
	PersonnelDisplay          string
	PersonnelRegistrationID   string
	PersonnelSection          string
	PersonnelOrganizationUnit string
	OperatorID                string
	OperatorDisplay           string
	OperatorRegistrationID    string
	Notes                     string
	CreatedAt                 string
	Lines                     []custodyReceiptLineView
	TotalQuantity             int
}

type custodyReceiptLineView struct {
	AssetID   string
	AssetName string
	Quantity  int
}

func newCustodyReceiptView(receipt app.CustodyReceipt) custodyReceiptView {
	view := custodyReceiptView{
		ID:                        string(receipt.ID),
		Type:                      string(receipt.Type),
		TypeLabel:                 custodyTransactionTypeLabel(receipt.Type),
		PersonnelID:               string(receipt.PersonnelID),
		PersonnelFullName:         receipt.PersonnelFullName,
		PersonnelAlias:            receipt.PersonnelAlias,
		PersonnelRankLabel:        receipt.PersonnelRank.Label(),
		PersonnelDisplay:          militaryDisplayName(receipt.PersonnelRank, receipt.PersonnelAlias),
		PersonnelRegistrationID:   receipt.PersonnelRegistrationID.String(),
		PersonnelSection:          receipt.PersonnelSection.Label(),
		PersonnelOrganizationUnit: receipt.PersonnelOrganizationUnit.Label(),
		OperatorID:                string(receipt.OperatorID),
		OperatorDisplay:           militaryDisplayName(receipt.OperatorRank, receipt.OperatorAlias),
		OperatorRegistrationID:    receipt.OperatorRegistrationID.String(),
		Notes:                     receipt.Notes,
		CreatedAt:                 formatReceiptTimestamp(receipt.CreatedAt),
		Lines:                     make([]custodyReceiptLineView, 0, len(receipt.Lines)),
	}

	for _, line := range receipt.Lines {
		view.TotalQuantity += line.Quantity
		view.Lines = append(view.Lines, custodyReceiptLineView{
			AssetID:   string(line.AssetID),
			AssetName: line.AssetName,
			Quantity:  line.Quantity,
		})
	}

	return view
}

func formatReceiptTimestamp(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Local().Format("2006-01-02 15:04")
}
