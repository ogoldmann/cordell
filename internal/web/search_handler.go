package web

import (
	"net/http"
	"net/url"
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
	Results      globalSearchResultsView
}

type globalSearchResultsView struct {
	Query        string
	EscapedQuery string
	HasQuery     bool
	Personnel    []personnelView
	Assets       []assetView
}

const globalSearchResultLimit = 5

func (s *Server) handleGlobalSearch(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))
	statusFilter := ports.NormalizeRecordStatusFilter(r.URL.Query().Get("status"))

	results, err := s.newGlobalSearchResultsView(r, searchQuery, statusFilter, globalSearchResultLimit)
	if err != nil {
		s.logger.Error("failed to run global search", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if wantsPartialResponse(r) {
		if err := s.renderer.Render(w, http.StatusOK, "global_search_results", results); err != nil {
			s.logger.Error("failed to render global search partial", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		return
	}

	data := globalSearchPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             searchLabel(),
		SearchQuery:       results.Query,
		StatusFilter:      string(statusFilter),
		StatusTabs:        newSearchStatusFilterTabs(results.Query, statusFilter),
		Results:           results,
	}
	data.Breadcrumbs = searchBreadcrumbs()

	if err := s.renderer.Render(w, http.StatusOK, "search.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleNavbarSearchSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	results, err := s.newGlobalSearchResultsView(
		r,
		query,
		ports.RecordStatusFilterActive,
		globalSearchResultLimit,
	)
	if err != nil {
		s.logger.Error("failed to build navbar search suggestions", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	results = limitGlobalSearchResults(results, 4)

	if err := s.renderer.Render(w, http.StatusOK, "navbar_search_suggestions", results); err != nil {
		s.logger.Error("failed to render navbar search suggestions", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) newGlobalSearchResultsView(
	r *http.Request,
	query string,
	statusFilter ports.RecordStatusFilter,
	limit int,
) (globalSearchResultsView, error) {
	query = strings.TrimSpace(query)

	view := globalSearchResultsView{
		Query:        query,
		EscapedQuery: url.QueryEscape(query),
		HasQuery:     query != "",
	}

	if query == "" {
		return view, nil
	}

	result, err := s.services.GlobalSearch.Execute(r.Context(), app.GlobalSearchCommand{
		Query:        query,
		LimitPerType: limit,
		StatusFilter: string(statusFilter),
	})
	if err != nil {
		return globalSearchResultsView{}, err
	}

	view.Query = result.Query
	view.EscapedQuery = url.QueryEscape(result.Query)
	view.HasQuery = result.Query != ""
	view.Personnel = make([]personnelView, 0, len(result.Personnel))
	view.Assets = make([]assetView, 0, len(result.Assets))

	for _, item := range result.Personnel {
		view.Personnel = append(view.Personnel, newPersonnelView(item))
	}

	for _, item := range result.Assets {
		view.Assets = append(view.Assets, newAssetView(item))
	}

	return view, nil
}

func limitGlobalSearchResults(results globalSearchResultsView, limit int) globalSearchResultsView {
	if limit <= 0 {
		return results
	}

	if len(results.Personnel) > limit {
		results.Personnel = results.Personnel[:limit]
	}

	if len(results.Assets) > limit {
		results.Assets = results.Assets[:limit]
	}

	return results
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
