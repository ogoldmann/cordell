package web

import (
	"net/http"

	"cordell/internal/app"
	"cordell/internal/domain"
)

type custodyTransactionLedgerPageData struct {
	privateLayoutData
	Title        string
	Transactions []custodyTransactionSummaryView
}

type custodyTransactionSummaryView struct {
	ID                         string
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
	CreatedAt                  string
	HasCorrection              bool
	EditLabel                  string
}

func (s *Server) handleCustodyTransactionLedger(w http.ResponseWriter, r *http.Request) {
	transactions, err := s.services.ListCustodyTransactionSummaries.Execute(r.Context(), app.ListCustodyTransactionSummariesCommand{
		Limit: 100,
	})
	if err != nil {
		s.logger.Error("failed to list custody transaction summaries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := custodyTransactionLedgerPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Custody transactions",
		Transactions:      newCustodyTransactionSummaryViews(transactions),
	}

	if err := s.renderer.Render(w, http.StatusOK, "custody_transactions_index.html", data); err != nil {
		s.handleRenderError(w, err)
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

		views = append(views, custodyTransactionSummaryView{
			ID:                         string(transaction.ID),
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
			CreatedAt:                  formatDateTime(transaction.CreatedAt),
			HasCorrection:              transaction.HasCorrection,
			EditLabel:                  editCountLabel(transaction.EditCount),
		})
	}

	return views
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
