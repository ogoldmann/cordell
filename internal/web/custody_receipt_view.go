package web

import (
	"sort"
	"strconv"
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
	EditCount                int
	HasEdits                 bool
	EditHistoryURL           string
	Current                  custodyReceiptCurrentView
	EditHistory              []custodyReceiptEditHistoryEntryView
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

type custodyReceiptCurrentView struct {
	PersonnelID             string
	PersonnelDisplay        string
	PersonnelFullName       string
	PersonnelRegistrationID string
	PersonnelActive         bool
	PersonnelStatusLabel    string
	Notes                   string
	HasNotes                bool
	Lines                   []custodyReceiptLineView
	TotalQuantity           int
}

type custodyReceiptEditHistoryEntryView struct {
	Kind           string
	Title          string
	CreatedAt      string
	OperatorLabel  string
	OperatorStatus string
	CorrectionID   string
	Changes        []custodyReceiptEditChangeView
	HasChanges     bool
}

type custodyReceiptEditChangeView struct {
	Label string
	From  string
	To    string
}

type custodyReceiptInterpretationView struct {
	PersonnelID             string
	PersonnelDisplay        string
	PersonnelFullName       string
	PersonnelRegistrationID string
	PersonnelActive         bool
	Notes                   string
	Lines                   []custodyReceiptLineView
}

func newCustodyReceiptView(receipt app.CustodyReceipt) custodyReceiptView {
	typeLabel := custodyTransactionTypeLabel(receipt.TransactionType)
	originalInterpretation := originalReceiptInterpretationView(receipt)
	effectiveInterpretation := originalInterpretation

	if receipt.HasCorrection {
		effectiveInterpretation = correctionInterpretationView(receipt.Correction)
	}

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
		EditCount:               receipt.EditCount,
		HasEdits:                receipt.EditCount > 0,
		EditHistoryURL:          "#edit-history",
		Current:                 newCustodyReceiptCurrentView(effectiveInterpretation),
		EditHistory:             newCustodyReceiptEditHistory(receipt),
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

func originalReceiptInterpretationView(receipt app.CustodyReceipt) custodyReceiptInterpretationView {
	lines := make([]custodyReceiptLineView, 0, len(receipt.Lines))

	for _, line := range receipt.Lines {
		lines = append(lines, custodyReceiptLineView{
			AssetID:     string(line.AssetID),
			AssetName:   line.AssetName,
			AssetActive: line.AssetActive,
			StatusLabel: activeStatusLabel(line.AssetActive),
			Quantity:    line.Quantity,
		})
	}

	return custodyReceiptInterpretationView{
		PersonnelID:             string(receipt.PersonnelID),
		PersonnelDisplay:        militaryDisplayName(receipt.PersonnelRank, receipt.PersonnelAlias),
		PersonnelFullName:       receipt.PersonnelFullName,
		PersonnelRegistrationID: string(receipt.PersonnelRegistrationID),
		PersonnelActive:         receipt.PersonnelActive,
		Notes:                   receipt.Notes,
		Lines:                   lines,
	}
}

func correctionInterpretationView(correction app.CustodyCorrectionContext) custodyReceiptInterpretationView {
	lines := make([]custodyReceiptLineView, 0, len(correction.Lines))

	for _, line := range correction.Lines {
		lines = append(lines, custodyReceiptLineView{
			AssetID:     string(line.AssetID),
			AssetName:   line.AssetName,
			AssetActive: line.AssetActive,
			StatusLabel: activeStatusLabel(line.AssetActive),
			Quantity:    line.Quantity,
		})
	}

	return custodyReceiptInterpretationView{
		PersonnelID:             string(correction.CorrectedPersonnelID),
		PersonnelDisplay:        militaryDisplayName(correction.CorrectedPersonnelRank, correction.CorrectedPersonnelAlias),
		PersonnelFullName:       correction.CorrectedPersonnelFullName,
		PersonnelRegistrationID: string(correction.CorrectedPersonnelRegistrationID),
		PersonnelActive:         correction.CorrectedPersonnelActive,
		Notes:                   correction.CorrectedNotes,
		Lines:                   lines,
	}
}

func newCustodyReceiptCurrentView(interpretation custodyReceiptInterpretationView) custodyReceiptCurrentView {
	current := custodyReceiptCurrentView{
		PersonnelID:             interpretation.PersonnelID,
		PersonnelDisplay:        interpretation.PersonnelDisplay,
		PersonnelFullName:       interpretation.PersonnelFullName,
		PersonnelRegistrationID: interpretation.PersonnelRegistrationID,
		PersonnelActive:         interpretation.PersonnelActive,
		PersonnelStatusLabel:    activeStatusLabel(interpretation.PersonnelActive),
		Notes:                   interpretation.Notes,
		HasNotes:                strings.TrimSpace(interpretation.Notes) != "",
		Lines:                   interpretation.Lines,
	}

	for _, line := range interpretation.Lines {
		current.TotalQuantity += line.Quantity
	}

	return current
}

func newCustodyReceiptEditHistory(receipt app.CustodyReceipt) []custodyReceiptEditHistoryEntryView {
	original := originalReceiptInterpretationView(receipt)

	history := []custodyReceiptEditHistoryEntryView{
		{
			Kind:           "original",
			Title:          "Original transaction",
			CreatedAt:      formatDateTime(receipt.CreatedAt),
			OperatorLabel:  militaryDisplayName(receipt.OperatorRank, receipt.OperatorAlias),
			OperatorStatus: activeStatusLabel(receipt.OperatorActive),
			Changes: []custodyReceiptEditChangeView{
				{
					Label: "Initial state",
					From:  "",
					To:    custodyReceiptInterpretationSummary(original),
				},
			},
			HasChanges: true,
		},
	}

	previous := original

	for index, correction := range receipt.Corrections {
		current := correctionInterpretationView(correction)
		changes := custodyReceiptChanges(previous, current)

		entry := custodyReceiptEditHistoryEntryView{
			Kind:           "edit",
			Title:          "Edit #" + strconv.Itoa(index+1),
			CreatedAt:      formatDateTime(correction.CreatedAt),
			OperatorLabel:  militaryDisplayName(correction.OperatorRank, correction.OperatorAlias),
			OperatorStatus: activeStatusLabel(correction.OperatorActive),
			CorrectionID:   string(correction.ID),
			Changes:        changes,
			HasChanges:     len(changes) > 0,
		}

		history = append(history, entry)
		previous = current
	}

	return history
}

func custodyReceiptChanges(
	previous custodyReceiptInterpretationView,
	current custodyReceiptInterpretationView,
) []custodyReceiptEditChangeView {
	changes := make([]custodyReceiptEditChangeView, 0)

	if previous.PersonnelID != current.PersonnelID {
		changes = append(changes, custodyReceiptEditChangeView{
			Label: "Personnel",
			From:  previous.PersonnelDisplay,
			To:    current.PersonnelDisplay,
		})
	}

	if strings.TrimSpace(previous.Notes) != strings.TrimSpace(current.Notes) {
		changes = append(changes, custodyReceiptEditChangeView{
			Label: "Notes",
			From:  custodyNotesChangeLabel(previous.Notes),
			To:    custodyNotesChangeLabel(current.Notes),
		})
	}

	changes = append(changes, custodyReceiptLineChanges(previous.Lines, current.Lines)...)

	return changes
}

func custodyNotesChangeLabel(notes string) string {
	if strings.TrimSpace(notes) == "" {
		return "No notes"
	}

	return "Notes recorded"
}

func custodyReceiptLineChanges(
	previous []custodyReceiptLineView,
	current []custodyReceiptLineView,
) []custodyReceiptEditChangeView {
	previousByAssetID := custodyReceiptLinesByAssetID(previous)
	currentByAssetID := custodyReceiptLinesByAssetID(current)

	assetIDs := make(map[string]struct{})

	for assetID := range previousByAssetID {
		assetIDs[assetID] = struct{}{}
	}

	for assetID := range currentByAssetID {
		assetIDs[assetID] = struct{}{}
	}

	orderedAssetIDs := make([]string, 0, len(assetIDs))
	for assetID := range assetIDs {
		orderedAssetIDs = append(orderedAssetIDs, assetID)
	}

	sort.Strings(orderedAssetIDs)

	changes := make([]custodyReceiptEditChangeView, 0)

	for _, assetID := range orderedAssetIDs {
		previousLine, hadPrevious := previousByAssetID[assetID]
		currentLine, hasCurrent := currentByAssetID[assetID]

		switch {
		case !hadPrevious && hasCurrent:
			changes = append(changes, custodyReceiptEditChangeView{
				Label: "Asset added",
				From:  "—",
				To:    custodyReceiptLineSummary(currentLine),
			})

		case hadPrevious && !hasCurrent:
			changes = append(changes, custodyReceiptEditChangeView{
				Label: "Asset removed",
				From:  custodyReceiptLineSummary(previousLine),
				To:    "—",
			})

		case hadPrevious && hasCurrent && previousLine.Quantity != currentLine.Quantity:
			changes = append(changes, custodyReceiptEditChangeView{
				Label: "Quantity changed",
				From:  custodyReceiptLineSummary(previousLine),
				To:    custodyReceiptLineSummary(currentLine),
			})
		}
	}

	return changes
}

func custodyReceiptLinesByAssetID(lines []custodyReceiptLineView) map[string]custodyReceiptLineView {
	linesByAssetID := make(map[string]custodyReceiptLineView)

	for _, line := range lines {
		existing, ok := linesByAssetID[line.AssetID]
		if ok {
			existing.Quantity += line.Quantity
			linesByAssetID[line.AssetID] = existing
			continue
		}

		linesByAssetID[line.AssetID] = line
	}

	return linesByAssetID
}

func custodyReceiptLineSummary(line custodyReceiptLineView) string {
	return line.AssetName + " ×" + strconv.Itoa(line.Quantity)
}

func custodyReceiptInterpretationSummary(interpretation custodyReceiptInterpretationView) string {
	return interpretation.PersonnelDisplay + " · " + strconv.Itoa(custodyReceiptTotalQuantity(interpretation.Lines)) + " item(s)"
}

func custodyReceiptTotalQuantity(lines []custodyReceiptLineView) int {
	total := 0

	for _, line := range lines {
		total += line.Quantity
	}

	return total
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
