package web

import (
	"net/http"

	"cordell/internal/app"
)

type dashboardPageData struct {
	Title           string
	RecentPersonnel []personnelView
	RecentAssets    []assetView
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	personnel, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit: 5,
	})
	if err != nil {
		s.logger.Error("failed to list recent personnel for dashboard", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	assets, err := s.services.ListAssets.Execute(r.Context(), app.ListAssetsCommand{
		Limit: 5,
	})
	if err != nil {
		s.logger.Error("failed to list recent assets for dashboard", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := dashboardPageData{
		Title:           "Dashboard",
		RecentPersonnel: make([]personnelView, 0, len(personnel)),
		RecentAssets:    make([]assetView, 0, len(assets)),
	}

	for _, item := range personnel {
		data.RecentPersonnel = append(data.RecentPersonnel, newPersonnelView(item))
	}

	for _, item := range assets {
		data.RecentAssets = append(data.RecentAssets, assetView{
			ID:     string(item.ID()),
			Name:   item.Name(),
			Active: item.Active(),
		})
	}

	if err := s.renderer.Render(w, http.StatusOK, "dashboard.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}
