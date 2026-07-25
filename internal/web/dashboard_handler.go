package web

import (
	"net/http"
)

type dashboardPageData struct {
	privateLayoutData
	Title           string
	WelcomePhrase   string
	Search          searchBarView
	OperationalDock dashboardOperationalDockView
}

func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	layout := newPrivateLayoutData(r)
	layout.HideNavbarSearch()
	layout.UseDefaultShell = false

	data := dashboardPageData{
		privateLayoutData: layout,
		Title:             dashboardLabel(),
		WelcomePhrase:     randomDashboardWelcomePhrase(),
		Search:            newHeroSearchBar("dashboard-global-search-input", "", "Ex.: sd silva, s1 john, rádio, identidade..."),
		OperationalDock:   newDashboardOperationalDockView(),
	}

	if err := s.renderer.Render(w, http.StatusOK, "dashboard.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}
