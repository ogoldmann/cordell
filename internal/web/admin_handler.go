package web

import "net/http"

type adminIndexPageData struct {
	privateLayoutData
	Title string
}

func (s *Server) handleAdminIndex(w http.ResponseWriter, r *http.Request) {
	if err := s.renderer.Render(w, http.StatusOK, "admin_index.html", adminIndexPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Admin",
	}); err != nil {
		s.handleRenderError(w, err)
	}
}
