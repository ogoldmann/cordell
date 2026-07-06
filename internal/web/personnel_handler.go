package web

import (
	"errors"
	"fmt"
	"net/http"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"

	"github.com/go-chi/chi/v5"
)

type personnelNewPageData struct {
	Title    string
	Error    string
	FullName string
}

type personnelShowPageData struct {
	Title     string
	Personnel personnelView
}

type personnelIndexPageData struct {
	Title     string
	Personnel []personnelView
}

type personnelView struct {
	ID       string
	FullName string
	Active   bool
}

func (s *Server) handleNewPersonnelForm(w http.ResponseWriter, r *http.Request) {
	data := personnelNewPageData{
		Title: "Create personnel",
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreatePersonnel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderNewPersonnelFormWithError(
			w,
			http.StatusBadRequest,
			"Invalid form submission.",
			"",
		)
		return
	}

	fullName := r.FormValue("full_name")

	personnel, err := s.services.CreatePersonnel.Execute(r.Context(), app.CreatePersonnelCommand{
		FullName: fullName,
	})
	if err != nil {
		s.renderNewPersonnelFormWithError(
			w,
			http.StatusBadRequest,
			humanizePersonnelError(err),
			fullName,
		)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/personnel/%s", personnel.ID()),
		http.StatusSeeOther,
	)
}

func (s *Server) handleShowPersonnel(w http.ResponseWriter, r *http.Request) {
	id := domain.PersonnelID(chi.URLParam(r, "id"))

	personnel, err := s.services.GetPersonnel.Execute(r.Context(), app.GetPersonnelCommand{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			http.NotFound(w, r)
			return
		}

		if errors.Is(err, domain.ErrEmptyPersonnelID) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to show personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := personnelShowPageData{
		Title: personnel.FullName(),
		Personnel: personnelView{
			ID:       string(personnel.ID()),
			FullName: personnel.FullName(),
			Active:   personnel.Active(),
		},
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_show.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) renderNewPersonnelFormWithError(
	w http.ResponseWriter,
	status int,
	message string,
	fullName string,
) {
	data := personnelNewPageData{
		Title:    "Create personnel",
		Error:    message,
		FullName: fullName,
	}

	if err := s.renderer.Render(w, status, "personnel_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleRenderError(w http.ResponseWriter, err error) {
	s.logger.Error("failed to render template", "error", err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func humanizePersonnelError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyPersonnelName):
		return "Full name is required."
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "Personnel ID is required."
	default:
		return "Could not create personnel."
	}
}

func (s *Server) handleListPersonnel(w http.ResponseWriter, r *http.Request) {
	personnel, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit: 50,
	})
	if err != nil {
		s.logger.Error("failed to list personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := personnelIndexPageData{
		Title:     "Personnel",
		Personnel: make([]personnelView, 0, len(personnel)),
	}

	for _, item := range personnel {
		data.Personnel = append(data.Personnel, personnelView{
			ID:       string(item.ID()),
			FullName: item.FullName(),
			Active:   item.Active(),
		})
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}
