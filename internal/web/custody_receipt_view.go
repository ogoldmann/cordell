package web

import (
	"strings"

	"cordell/internal/app"
)

type custodyReceiptView struct {
	ID                       string
	TypeLabel                string
	PersonnelID              string
	PersonnelDisplay         string
	PersonnelFullName        string
	PersonnelRegistrationID  string
	PersonnelActive          bool
	PersonnelStatusLabel     string
	OperatorID               string
	OperatorDisplay          string
	OperatorRoleLabel        string
	OperatorActive           bool
	OperatorStatusLabel      string
	Notes                    string
	HasNotes                 bool
	CreatedAt                string
	Lines                    []custodyReceiptLineView
	TotalQuantity            int
	HasInactiveRelatedRecord bool
}

type custodyReceiptLineView struct {
	AssetID     string
	AssetName   string
	AssetActive bool
	StatusLabel string
	Quantity    int
}

func newCustodyReceiptView(receipt app.CustodyReceipt) custodyReceiptView {
	view := custodyReceiptView{
		ID:                      string(receipt.ID),
		TypeLabel:               custodyTransactionTypeLabel(receipt.TransactionType),
		PersonnelID:             string(receipt.PersonnelID),
		PersonnelDisplay:        militaryDisplayName(receipt.PersonnelRank, receipt.PersonnelAlias),
		PersonnelFullName:       receipt.PersonnelFullName,
		PersonnelRegistrationID: string(receipt.PersonnelRegistrationID),
		PersonnelActive:         receipt.PersonnelActive,
		PersonnelStatusLabel:    activeStatusLabel(receipt.PersonnelActive),
		OperatorID:              string(receipt.OperatorID),
		OperatorDisplay:         militaryDisplayName(receipt.OperatorRank, receipt.OperatorAlias),
		OperatorRoleLabel:       receipt.OperatorRole.Label(),
		OperatorActive:          receipt.OperatorActive,
		OperatorStatusLabel:     activeStatusLabel(receipt.OperatorActive),
		Notes:                   receipt.Notes,
		HasNotes:                strings.TrimSpace(receipt.Notes) != "",
		CreatedAt:               formatDateTime(receipt.CreatedAt),
		Lines:                   make([]custodyReceiptLineView, 0, len(receipt.Lines)),
	}

	for _, line := range receipt.Lines {
		statusLabel := activeStatusLabel(line.AssetActive)

		if !line.AssetActive {
			view.HasInactiveRelatedRecord = true
		}

		view.TotalQuantity += line.Quantity
		view.Lines = append(view.Lines, custodyReceiptLineView{
			AssetID:     string(line.AssetID),
			AssetName:   line.AssetName,
			AssetActive: line.AssetActive,
			StatusLabel: statusLabel,
			Quantity:    line.Quantity,
		})
	}

	if !receipt.PersonnelActive || !receipt.OperatorActive {
		view.HasInactiveRelatedRecord = true
	}

	return view
}

func activeStatusLabel(active bool) string {
	if active {
		return "Active"
	}

	return "Inactive"
}
