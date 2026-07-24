package web

import (
	"net/http"
	"strings"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"
)

type globalSearchPageData struct {
	privateLayoutData
	Title        string
	SearchQuery  string
	StatusFilter string
	StatusTabs   []statusFilterTabView
	HasQuery     bool
	HasResults   bool
	Personnel    []personnelView
	Assets       []assetView
}

func (s *Server) handleGlobalSearch(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))
	statusFilter := ports.NormalizeRecordStatusFilter(r.URL.Query().Get("status"))

	result, err := s.services.GlobalSearch.Execute(r.Context(), app.GlobalSearchCommand{
		Query:        searchQuery,
		LimitPerType: 8,
		StatusFilter: string(statusFilter),
	})
	if err != nil {
		s.logger.Error("failed to run global search", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := globalSearchPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             searchLabel(),
		SearchQuery:       result.Query,
		StatusFilter:      string(statusFilter),
		StatusTabs:        newSearchStatusFilterTabs(result.Query, statusFilter),
		HasQuery:          result.Query != "",
		Personnel:         make([]personnelView, 0, len(result.Personnel)),
		Assets:            make([]assetView, 0, len(result.Assets)),
	}
	data.Breadcrumbs = searchBreadcrumbs()

	for _, item := range result.Personnel {
		data.Personnel = append(data.Personnel, newPersonnelView(item))
	}

	for _, item := range result.Assets {
		data.Assets = append(data.Assets, newAssetView(item))
	}

	data.HasResults = len(data.Personnel) > 0 || len(data.Assets) > 0

	if err := s.renderer.Render(w, http.StatusOK, "search.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newAssetView(asset domain.Asset) assetView {
	statusLabel := activeStatusLabel(asset.Active())

	return assetView{
		ID:            string(asset.ID()),
		Name:          asset.Name(),
		Active:        asset.Active(),
		StatusLabel:   statusLabel,
		CanDeactivate: asset.Active(),
		CanReactivate: !asset.Active(),
	}
}
