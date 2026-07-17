package web

import (
	"net/http"

	"cordell/internal/app"
)

type adminIndexPageData struct {
	privateLayoutData
	Title string
}

type adminOperatorsIndexPageData struct {
	privateLayoutData
	Title     string
	Operators []operatorSummaryView
}

func (s *Server) handleAdminIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.renderer.Render(w, http.StatusOK, "admin_index.html", adminIndexPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Admin",
	}); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleAdminOperatorsIndex(w http.ResponseWriter, r *http.Request) {
	operators, err := s.services.ListOperators.Execute(r.Context(), app.ListOperatorsCommand{
		Limit: 100,
	})
	if err != nil {
		s.logger.Error("failed to list operators", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := adminOperatorsIndexPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Operators",
		Operators:         make([]operatorSummaryView, 0, len(operators)),
	}

	for _, operator := range operators {
		data.Operators = append(data.Operators, newOperatorSummaryView(operator))
	}

	if err := s.renderer.Render(w, http.StatusOK, "admin_operators_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}
