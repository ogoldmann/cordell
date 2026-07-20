package web

import (
	"strings"

	"cordell/internal/app"
)

type custodyReceiptView struct {
	ID                       string
	TypeLabel                string
	EditURL                  string
	EditLabel                string
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
	HasCorrection            bool
	Correction               custodyReceiptCorrectionView
}

type custodyReceiptLineView struct {
	AssetID     string
	AssetName   string
	AssetActive bool
	StatusLabel string
	Quantity    int
}

type custodyReceiptCorrectionView struct {
	ID                               string
	CorrectedPersonnelID             string
	CorrectedPersonnelDisplay        string
	CorrectedPersonnelFullName       string
	CorrectedPersonnelRegistrationID string
	CorrectedPersonnelActive         bool
	CorrectedPersonnelStatusLabel    string
	OperatorID                       string
	OperatorDisplay                  string
	OperatorRoleLabel                string
	OperatorActive                   bool
	OperatorStatusLabel              string
	CorrectedNotes                   string
	HasCorrectedNotes                bool
	CreatedAt                        string
	Lines                            []custodyReceiptLineView
	TotalQuantity                    int
	HasInactiveRelatedRecord         bool
}

func newCustodyReceiptView(receipt app.CustodyReceipt) custodyReceiptView {
	typeLabel := custodyTransactionTypeLabel(receipt.TransactionType)

	view := custodyReceiptView{
		ID:                      string(receipt.ID),
		TypeLabel:               typeLabel,
		EditURL:                 "/custody/transactions/" + string(receipt.ID) + "/edit",
		EditLabel:               "Edit " + strings.ToLower(typeLabel),
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

	if receipt.HasCorrection {
		view.HasCorrection = true
		view.Correction = newCustodyReceiptCorrectionView(receipt.Correction)
		if view.Correction.HasInactiveRelatedRecord {
			view.HasInactiveRelatedRecord = true
		}
	}

	return view
}

func newCustodyReceiptCorrectionView(correction app.CustodyCorrectionContext) custodyReceiptCorrectionView {
	view := custodyReceiptCorrectionView{
		ID:                               string(correction.ID),
		CorrectedPersonnelID:             string(correction.CorrectedPersonnelID),
		CorrectedPersonnelDisplay:        militaryDisplayName(correction.CorrectedPersonnelRank, correction.CorrectedPersonnelAlias),
		CorrectedPersonnelFullName:       correction.CorrectedPersonnelFullName,
		CorrectedPersonnelRegistrationID: string(correction.CorrectedPersonnelRegistrationID),
		CorrectedPersonnelActive:         correction.CorrectedPersonnelActive,
		CorrectedPersonnelStatusLabel:    activeStatusLabel(correction.CorrectedPersonnelActive),
		OperatorID:                       string(correction.OperatorID),
		OperatorDisplay:                  militaryDisplayName(correction.OperatorRank, correction.OperatorAlias),
		OperatorRoleLabel:                correction.OperatorRole.Label(),
		OperatorActive:                   correction.OperatorActive,
		OperatorStatusLabel:              activeStatusLabel(correction.OperatorActive),
		CorrectedNotes:                   correction.CorrectedNotes,
		HasCorrectedNotes:                strings.TrimSpace(correction.CorrectedNotes) != "",
		CreatedAt:                        formatDateTime(correction.CreatedAt),
		Lines:                            make([]custodyReceiptLineView, 0, len(correction.Lines)),
	}

	for _, line := range correction.Lines {
		if !line.AssetActive {
			view.HasInactiveRelatedRecord = true
		}

		view.TotalQuantity += line.Quantity
		view.Lines = append(view.Lines, custodyReceiptLineView{
			AssetID:     string(line.AssetID),
			AssetName:   line.AssetName,
			AssetActive: line.AssetActive,
			StatusLabel: activeStatusLabel(line.AssetActive),
			Quantity:    line.Quantity,
		})
	}

	if !correction.CorrectedPersonnelActive || !correction.OperatorActive {
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
