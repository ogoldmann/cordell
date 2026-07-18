package web

import "cordell/internal/ports"

type statusFilterTabView struct {
	Label  string
	Value  string
	URL    string
	Active bool
}

func newStatusFilterTabs(basePath string, selected ports.RecordStatusFilter) []statusFilterTabView {
	return []statusFilterTabView{
		{
			Label:  "Active",
			Value:  string(ports.RecordStatusFilterActive),
			URL:    basePath + "?status=active",
			Active: selected == ports.RecordStatusFilterActive,
		},
		{
			Label:  "Inactive",
			Value:  string(ports.RecordStatusFilterInactive),
			URL:    basePath + "?status=inactive",
			Active: selected == ports.RecordStatusFilterInactive,
		},
		{
			Label:  "All",
			Value:  string(ports.RecordStatusFilterAll),
			URL:    basePath + "?status=all",
			Active: selected == ports.RecordStatusFilterAll,
		},
	}
}
