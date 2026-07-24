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

func newStatusFilterTabs(basePath string, selected ports.RecordStatusFilter, query string) []statusFilterTabView {
	return []statusFilterTabView{
		newStatusFilterTab(basePath, activeStatusLabel(true), ports.RecordStatusFilterActive, selected, query),
		newStatusFilterTab(basePath, activeStatusLabel(false), ports.RecordStatusFilterInactive, selected, query),
		newStatusFilterTab(basePath, allStatusLabel(), ports.RecordStatusFilterAll, selected, query),
	}
}

func newStatusFilterTab(
	basePath string,
	label string,
	value ports.RecordStatusFilter,
	selected ports.RecordStatusFilter,
	query string,
) statusFilterTabView {
	values := url.Values{}
	values.Set("status", string(value))

	if strings.TrimSpace(query) != "" {
		values.Set("q", query)
	}

	return statusFilterTabView{
		Label:  label,
		Value:  string(value),
		URL:    basePath + "?" + values.Encode(),
		Active: selected == value,
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
