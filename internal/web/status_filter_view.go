package web

import (
	"net/url"
	"strings"

	"cordell/internal/ports"
)

type statusFilterTabView struct {
	Label  string
	Value  string
	URL    string
	Active bool
}

func newStatusFilterTabs(basePath string, selected ports.RecordStatusFilter) []statusFilterTabView {
	return []statusFilterTabView{
		{
			Label:  activeStatusLabel(true),
			Value:  string(ports.RecordStatusFilterActive),
			URL:    basePath + "?status=active",
			Active: selected == ports.RecordStatusFilterActive,
		},
		{
			Label:  activeStatusLabel(false),
			Value:  string(ports.RecordStatusFilterInactive),
			URL:    basePath + "?status=inactive",
			Active: selected == ports.RecordStatusFilterInactive,
		},
		{
			Label:  allStatusLabel(),
			Value:  string(ports.RecordStatusFilterAll),
			URL:    basePath + "?status=all",
			Active: selected == ports.RecordStatusFilterAll,
		},
	}
}

func newSearchStatusFilterTabs(query string, selected ports.RecordStatusFilter) []statusFilterTabView {
	return []statusFilterTabView{
		newSearchStatusFilterTab(activeStatusLabel(true), ports.RecordStatusFilterActive, query, selected),
		newSearchStatusFilterTab(activeStatusLabel(false), ports.RecordStatusFilterInactive, query, selected),
		newSearchStatusFilterTab(allStatusLabel(), ports.RecordStatusFilterAll, query, selected),
	}
}

func newSearchStatusFilterTab(
	label string,
	value ports.RecordStatusFilter,
	query string,
	selected ports.RecordStatusFilter,
) statusFilterTabView {
	values := url.Values{}
	values.Set("status", string(value))

	if strings.TrimSpace(query) != "" {
		values.Set("q", query)
	}

	return statusFilterTabView{
		Label:  label,
		Value:  string(value),
		URL:    "/search?" + values.Encode(),
		Active: selected == value,
	}
}
