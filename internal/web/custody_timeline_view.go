package web

import (
	"strconv"
	"strings"
	"time"

	"cordell/internal/domain"
)

type custodyTimelineView struct {
	Items      []custodyTimelineItemView
	EmptyState emptyStateView
}

type custodyTimelineItemView struct {
	ID                   string
	URL                  string
	SequenceLabel        string
	Type                 string
	TypeLabel            string
	TypeClass            string
	TypeTone             string
	ReceiptURL           string
	RegisteredBy         string
	RegisteredAt         string
	PersonnelDisplayName string
	PersonnelFullName    string
	DateLabel            string
	TimeLabel            string
	PersonnelLabel       string
	PersonnelURL         string
	OperatorLabel        string
	OperatorURL          string
	Edited               bool
	HasCorrections       bool
	CorrectionCount      int
	EditCountLabel       string
	Notes                string
	Lines                []custodyTimelineLineView
	PrimaryActionLabel   string
	PrimaryActionURL     string
	SecondaryActionLabel string
	SecondaryActionURL   string
}

type custodyTimelineLineView struct {
	AssetID     string
	AssetName   string
	AssetURL    string
	Quantity    string
	Highlighted bool
	Highlight   bool
}

func custodyTransactionTypeClass(transactionType string) string {
	switch transactionType {
	case string(domain.CustodyTransactionTypeCheckout):
		return "checkout"
	case string(domain.CustodyTransactionTypeReturn):
		return "return"
	default:
		return "neutral"
	}
}

func custodyTimelineTypeLabel(typeLabel string) string {
	if strings.TrimSpace(typeLabel) == "" {
		return "TRANSAÇÃO"
	}

	return strings.ToUpper(typeLabel)
}

func custodyTimelineTypeTone(typeLabel string) string {
	switch typeLabel {
	case checkoutLabel(), "CAUTELA":
		return "checkout"
	case returnLabel(), "DESCAUTELA":
		return "return"
	default:
		return "neutral"
	}
}

func custodyTimelineRegisteredAt(dateLabel string, timeLabel string) string {
	switch {
	case dateLabel != "" && timeLabel != "":
		return dateLabel + " " + timeLabel
	case dateLabel != "":
		return dateLabel
	default:
		return timeLabel
	}
}

func custodyEditCountLabel(editCount int) string {
	if editCount <= 0 {
		return ""
	}

	if editCount == 1 {
		return "1 edição"
	}

	return pluralizeInt(editCount, "edição", "edições")
}

func pluralizeInt(value int, singular string, plural string) string {
	if value == 1 {
		return "1 " + singular
	}

	return strconv.Itoa(value) + " " + plural
}

func formatTimelineDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Local().Format("02/01/2006")
}

func formatTimelineTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Local().Format("15:04")
}

func formatTimelineQuantity(quantity int) string {
	return strconv.Itoa(quantity)
}
