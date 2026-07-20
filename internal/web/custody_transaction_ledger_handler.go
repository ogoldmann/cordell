package web

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"cordell/internal/app"
	"cordell/internal/domain"
)

type custodyTransactionLedgerPageData struct {
	privateLayoutData
	Title            string
	Query            string
	TypeFilter       string
	EditStatusFilter string
	HasFilters       bool
	Transactions     []custodyTransactionSummaryView
}

type custodyTransactionSummaryView struct {
	ID                         string
	SequenceLabel              string
	DateLabel                  string
	TimeLabel                  string
	ReceiptURL                 string
	TypeLabel                  string
	TypeBadgeClass             string
	EffectivePersonnelURL      string
	EffectivePersonnelDisplay  string
	EffectivePersonnelFullName string
	EffectivePersonnelActive   bool
	EffectivePersonnelStatus   string
	OriginalPersonnelDisplay   string
	ShowOriginalPersonnel      bool
	OperatorDisplay            string
	OperatorStatus             string
	TotalQuantity              int
	Lines                      []custodyTransactionSummaryLineView
	CreatedAt                  string
	HasCorrection              bool
	EditLabel                  string
}

type custodyTransactionSummaryLineView struct {
	AssetID     string
	AssetURL    string
	AssetName   string
	AssetActive bool
	StatusLabel string
	Quantity    int
}

func (s *Server) handleCustodyTransactionLedger(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	typeFilter := normalizedLedgerTypeFilter(r.URL.Query().Get("type"))
	editStatusFilter := normalizedLedgerEditStatusFilter(r.URL.Query().Get("edited"))

	transactions, err := s.services.ListCustodyTransactionSummaries.Execute(r.Context(), app.ListCustodyTransactionSummariesCommand{
		Limit:                 100,
		SearchQuery:           query,
		TransactionTypeFilter: typeFilter,
		EditStatusFilter:      editStatusFilter,
	})
	if err != nil {
		s.logger.Error("failed to list custody transaction summaries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := custodyTransactionLedgerPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Custody transactions",
		Query:             query,
		TypeFilter:        typeFilter,
		EditStatusFilter:  editStatusFilter,
		HasFilters:        query != "" || typeFilter != "all" || editStatusFilter != "all",
		Transactions:      newCustodyTransactionSummaryViews(transactions),
	}

	if err := s.renderer.Render(w, http.StatusOK, "custody_transactions_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func normalizedLedgerTypeFilter(value string) string {
	switch value {
	case "checkout", "return":
		return value
	default:
		return "all"
	}
}

func normalizedLedgerEditStatusFilter(value string) string {
	switch value {
	case "edited", "unedited":
		return value
	default:
		return "all"
	}
}

func newCustodyTransactionSummaryViews(
	transactions []app.CustodyTransactionSummary,
) []custodyTransactionSummaryView {
	views := make([]custodyTransactionSummaryView, 0, len(transactions))

	for _, transaction := range transactions {
		effectivePersonnelDisplay := militaryDisplayName(
			transaction.EffectivePersonnelRank,
			transaction.EffectivePersonnelAlias,
		)

		originalPersonnelDisplay := militaryDisplayName(
			transaction.OriginalPersonnelRank,
			transaction.OriginalPersonnelAlias,
		)

		showOriginalPersonnel := transaction.HasCorrection &&
			transaction.OriginalPersonnelID != transaction.EffectivePersonnelID

		lineViews := make([]custodyTransactionSummaryLineView, 0, len(transaction.Lines))

		for _, line := range transaction.Lines {
			lineViews = append(lineViews, custodyTransactionSummaryLineView{
				AssetID:     string(line.AssetID),
				AssetURL:    "/assets/" + string(line.AssetID),
				AssetName:   line.AssetName,
				AssetActive: line.AssetActive,
				StatusLabel: activeStatusLabel(line.AssetActive),
				Quantity:    line.Quantity,
			})
		}

		views = append(views, custodyTransactionSummaryView{
			ID:                         string(transaction.ID),
			SequenceLabel:              ledgerSequenceLabel(transaction.SequenceNumber),
			DateLabel:                  ledgerDateLabel(transaction.CreatedAt),
			TimeLabel:                  ledgerTimeLabel(transaction.CreatedAt),
			ReceiptURL:                 "/custody/transactions/" + string(transaction.ID),
			TypeLabel:                  custodyTransactionTypeLabel(transaction.TransactionType),
			TypeBadgeClass:             custodyTransactionTypeBadgeClass(transaction.TransactionType),
			EffectivePersonnelURL:      "/personnel/" + string(transaction.EffectivePersonnelID),
			EffectivePersonnelDisplay:  effectivePersonnelDisplay,
			EffectivePersonnelFullName: transaction.EffectivePersonnelFullName,
			EffectivePersonnelActive:   transaction.EffectivePersonnelActive,
			EffectivePersonnelStatus:   activeStatusLabel(transaction.EffectivePersonnelActive),
			OriginalPersonnelDisplay:   originalPersonnelDisplay,
			ShowOriginalPersonnel:      showOriginalPersonnel,
			OperatorDisplay:            militaryDisplayName(transaction.OperatorRank, transaction.OperatorAlias),
			OperatorStatus:             activeStatusLabel(transaction.OperatorActive),
			TotalQuantity:              transaction.TotalQuantity,
			Lines:                      lineViews,
			CreatedAt:                  formatDateTime(transaction.CreatedAt),
			HasCorrection:              transaction.HasCorrection,
			EditLabel:                  editCountLabel(transaction.EditCount),
		})
	}

	return views
}

func ledgerDateLabel(t time.Time) string {
	return t.Format("02 Jan 2006")
}

func ledgerTimeLabel(t time.Time) string {
	return t.Format("15:04")
}

func ledgerSequenceLabel(sequenceNumber int) string {
	if sequenceNumber <= 0 {
		return "#—"
	}

	return "#" + strconv.Itoa(sequenceNumber)
}

func custodyTransactionTypeBadgeClass(transactionType domain.CustodyTransactionType) string {
	switch transactionType {
	case domain.CustodyTransactionTypeCheckout:
		return "bg-emerald-50 text-emerald-700 border-emerald-200 dark:bg-emerald-950/40 dark:text-emerald-200 dark:border-emerald-900/60"
	case domain.CustodyTransactionTypeReturn:
		return "bg-sky-50 text-sky-700 border-sky-200 dark:bg-sky-950/40 dark:text-sky-200 dark:border-sky-900/60"
	default:
		return "bg-slate-100 text-slate-600 border-slate-200 dark:bg-slate-800 dark:text-slate-300 dark:border-slate-700"
	}
}
