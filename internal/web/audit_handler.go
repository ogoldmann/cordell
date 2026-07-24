package web

import (
	"net/http"

	"cordell/internal/app"
)

type auditEventsIndexPageData struct {
	privateLayoutData
	Title  string
	Events []auditEventView
}

func (s *Server) handleAdminAuditEventsIndex(w http.ResponseWriter, r *http.Request) {
	entries, err := s.services.ListAuditEvents.Execute(r.Context(), app.ListAuditEventsCommand{
		Limit: 100,
	})
	if err != nil {
		s.logger.Error("failed to list audit events", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := auditEventsIndexPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Audit events",
		Events:            make([]auditEventView, 0, len(entries)),
	}
	data.Breadcrumbs = adminAuditEventsBreadcrumbs()

	for _, entry := range entries {
		data.Events = append(data.Events, newAuditEventView(entry))
	}

	if err := s.renderer.Render(w, http.StatusOK, "admin_audit_events_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}
