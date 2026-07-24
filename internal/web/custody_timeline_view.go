package web

import (
	"strconv"
	"time"
)

type custodyTimelineView struct {
	Items      []custodyTimelineItemView
	EmptyState emptyStateView
}

type custodyTimelineItemView struct {
	ID                   string
	URL                  string
	SequenceLabel        string
	TypeLabel            string
	TypeTone             string
	DateLabel            string
	TimeLabel            string
	PersonnelLabel       string
	PersonnelURL         string
	OperatorLabel        string
	OperatorURL          string
	Edited               bool
	EditCountLabel       string
	Notes                string
	Lines                []custodyTimelineLineView
	PrimaryActionLabel   string
	PrimaryActionURL     string
	SecondaryActionLabel string
	SecondaryActionURL   string
}

type custodyTimelineLineView struct {
	AssetID   string
	AssetName string
	AssetURL  string
	Quantity  string
}

func custodyTimelineTypeTone(typeLabel string) string {
	switch typeLabel {
	case checkoutLabel():
		return "checkout"
	case returnLabel():
		return "return"
	default:
		return "neutral"
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
