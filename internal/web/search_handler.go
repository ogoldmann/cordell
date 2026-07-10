package web

import (
	"net/http"
	"strings"

	"cordell/internal/app"
	"cordell/internal/domain"
)

type globalSearchPageData struct {
	Title       string
	SearchQuery string
	HasQuery    bool
	HasResults  bool
	Personnel   []personnelView
	Assets      []assetView
}

func (s *Server) handleGlobalSearch(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	result, err := s.services.GlobalSearch.Execute(r.Context(), app.GlobalSearchCommand{
		Query:        searchQuery,
		LimitPerType: 8,
	})
	if err != nil {
		s.logger.Error("failed to run global search", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := globalSearchPageData{
		Title:       "Search",
		SearchQuery: result.Query,
		HasQuery:    result.Query != "",
		Personnel:   make([]personnelView, 0, len(result.Personnel)),
		Assets:      make([]assetView, 0, len(result.Assets)),
	}

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
	return assetView{
		ID:     string(asset.ID()),
		Name:   asset.Name(),
		Active: asset.Active(),
	}
}
