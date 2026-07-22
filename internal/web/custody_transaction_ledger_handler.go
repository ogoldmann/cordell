package web

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"cordell/internal/app"
	"cordell/internal/domain"
)

type custodyTransactionLedgerPageData struct {
	privateLayoutData
	Title               string
	Query               string
	TypeFilter          string
	EditStatusFilter    string
	HasFilters          bool
	SelectedYear        int
	SelectedMonth       int
	SelectedPeriodLabel string
	PeriodValue         string
	Periods             []custodyLedgerPeriodView
	Page                int
	PageSize            int
	HasPreviousPage     bool
	HasNextPage         bool
	PreviousPageURL     string
	NextPageURL         string
	Transactions        []custodyTransactionSummaryView
}

type custodyLedgerPeriodView struct {
	Year             int
	Month            int
	Label            string
	URL              string
	Selected         bool
	TransactionCount int
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

	periods, err := s.services.ListCustodyTransactionLedgerPeriods.Execute(r.Context())
	if err != nil {
		s.logger.Error("failed to list custody transaction ledger periods", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageNumber := parsePositiveLedgerInt(r.URL.Query().Get("page"))
	if pageNumber <= 0 {
		pageNumber = 1
	}

	selectedYear, selectedMonth, periodValue := parseLedgerPeriodValue(r.URL.Query().Get("period"), periods)

	transactionPage, err := s.services.ListCustodyTransactionSummaries.Execute(r.Context(), app.ListCustodyTransactionSummariesCommand{
		Page:                  pageNumber,
		SearchQuery:           query,
		TransactionTypeFilter: typeFilter,
		EditStatusFilter:      editStatusFilter,
		Year:                  selectedYear,
		Month:                 selectedMonth,
	})
	if err != nil {
		s.logger.Error("failed to list custody transaction summaries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	previousPageURL := ""
	if transactionPage.Page > 1 {
		previousPageURL = custodyLedgerURL(query, typeFilter, editStatusFilter, periodValue, transactionPage.Page-1)
	}

	nextPageURL := ""
	if transactionPage.HasNextPage {
		nextPageURL = custodyLedgerURL(query, typeFilter, editStatusFilter, periodValue, transactionPage.Page+1)
	}

	data := custodyTransactionLedgerPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Custody transactions",
		Query:             query,
		TypeFilter:        typeFilter,
		EditStatusFilter:  editStatusFilter,
		HasFilters: query != "" ||
			typeFilter != "all" ||
			editStatusFilter != "all" ||
			periodValue == "all" ||
			transactionPage.Page > 1,
		SelectedYear:        selectedYear,
		SelectedMonth:       selectedMonth,
		SelectedPeriodLabel: ledgerPeriodLabel(selectedYear, selectedMonth),
		PeriodValue:         periodValue,
		Periods: newCustodyLedgerPeriodViews(
			periods,
			periodValue,
			query,
			typeFilter,
			editStatusFilter,
		),
		Page:            transactionPage.Page,
		PageSize:        transactionPage.PageSize,
		HasPreviousPage: transactionPage.Page > 1,
		HasNextPage:     transactionPage.HasNextPage,
		PreviousPageURL: previousPageURL,
		NextPageURL:     nextPageURL,
		Transactions:    newCustodyTransactionSummaryViews(transactionPage.Items),
	}

	if err := s.renderer.Render(w, http.StatusOK, "custody_transactions_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func parsePositiveLedgerInt(value string) int {
	parsed, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || parsed <= 0 {
		return 0
	}

	return parsed
}

func parseLedgerPeriodValue(value string, periods []app.CustodyTransactionLedgerPeriod) (int, int, string) {
	value = strings.TrimSpace(value)

	if value == "all" {
		return 0, 0, "all"
	}

	if value != "" {
		parts := strings.Split(value, "-")
		if len(parts) == 2 {
			year, yearErr := strconv.Atoi(parts[0])
			month, monthErr := strconv.Atoi(parts[1])

			if yearErr == nil && monthErr == nil && month >= 1 && month <= 12 {
				return year, month, value
			}
		}
	}

	if len(periods) == 0 {
		return 0, 0, "all"
	}

	year := periods[0].Year
	month := periods[0].Month

	return year, month, ledgerPeriodValue(year, month)
}

func ledgerPeriodValue(year int, month int) string {
	if year <= 0 || month < 1 || month > 12 {
		return "all"
	}

	return strconv.Itoa(year) + "-" + twoDigitMonth(month)
}

func twoDigitMonth(month int) string {
	if month < 10 {
		return "0" + strconv.Itoa(month)
	}

	return strconv.Itoa(month)
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

func ledgerPeriodLabel(year int, month int) string {
	if year <= 0 || month < 1 || month > 12 {
		return "All periods"
	}

	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)

	return t.Format("January 2006")
}

func custodyLedgerURL(
	query string,
	typeFilter string,
	editStatusFilter string,
	periodValue string,
	page int,
) string {
	values := make(url.Values)

	if strings.TrimSpace(query) != "" {
		values.Set("q", strings.TrimSpace(query))
	}

	if typeFilter != "" && typeFilter != "all" {
		values.Set("type", typeFilter)
	}

	if editStatusFilter != "" && editStatusFilter != "all" {
		values.Set("edited", editStatusFilter)
	}

	if periodValue != "" && periodValue != "all" {
		values.Set("period", periodValue)
	}

	if periodValue == "all" {
		values.Set("period", "all")
	}

	if page > 1 {
		values.Set("page", strconv.Itoa(page))
	}

	encoded := values.Encode()
	if encoded == "" {
		return "/custody/transactions"
	}

	return "/custody/transactions?" + encoded
}

func newCustodyLedgerPeriodViews(
	periods []app.CustodyTransactionLedgerPeriod,
	selectedPeriodValue string,
	query string,
	typeFilter string,
	editStatusFilter string,
) []custodyLedgerPeriodView {
	views := make([]custodyLedgerPeriodView, 0, len(periods)+1)

	views = append(views, custodyLedgerPeriodView{
		Year:             0,
		Month:            0,
		Label:            "All periods",
		URL:              custodyLedgerURL(query, typeFilter, editStatusFilter, "all", 1),
		Selected:         selectedPeriodValue == "all",
		TransactionCount: 0,
	})

	for _, period := range periods {
		value := ledgerPeriodValue(period.Year, period.Month)

		views = append(views, custodyLedgerPeriodView{
			Year:             period.Year,
			Month:            period.Month,
			Label:            ledgerPeriodLabel(period.Year, period.Month),
			URL:              custodyLedgerURL(query, typeFilter, editStatusFilter, value, 1),
			Selected:         selectedPeriodValue == value,
			TransactionCount: period.TransactionCount,
		})
	}

	return views
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
