package web

import (
	"errors"
	"net/http"

	"cordell/internal/app"
	"cordell/internal/domain"
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

type adminOperatorNewPageData struct {
	privateLayoutData
	Title           string
	Error           string
	Username        string
	SelectedRole    string
	RoleOptions     []operatorRoleOptionView
	Password        string
	ConfirmPassword string
}

type operatorRoleOptionView struct {
	Value    string
	Label    string
	Selected bool
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

func (s *Server) handleNewAdminOperatorForm(w http.ResponseWriter, r *http.Request) {
	data := newAdminOperatorNewPageData(r, adminOperatorNewPageData{
		Title:        "Create operator",
		SelectedRole: domain.OperatorRoleOperator.String(),
	})

	if err := s.renderer.Render(w, http.StatusOK, "admin_operator_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateAdminOperator(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	role := r.FormValue("role")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		s.renderAdminOperatorFormWithError(
			w,
			r,
			http.StatusBadRequest,
			"Passwords do not match.",
			username,
			role,
		)
		return
	}

	_, err := s.services.CreateOperator.Execute(r.Context(), app.CreateOperatorCommand{
		Username: username,
		Role:     role,
		Password: password,
	})
	if err != nil {
		s.renderAdminOperatorFormWithError(
			w,
			r,
			http.StatusBadRequest,
			humanizeCreateOperatorWebError(err),
			username,
			role,
		)
		return
	}

	http.Redirect(w, r, "/admin/operators", http.StatusSeeOther)
}

func (s *Server) renderAdminOperatorFormWithError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message string,
	username string,
	role string,
) {
	data := newAdminOperatorNewPageData(r, adminOperatorNewPageData{
		Title:        "Create operator",
		Error:        message,
		Username:     username,
		SelectedRole: role,
	})

	if err := s.renderer.Render(w, status, "admin_operator_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newAdminOperatorNewPageData(
	r *http.Request,
	data adminOperatorNewPageData,
) adminOperatorNewPageData {
	data.privateLayoutData = newPrivateLayoutData(r)

	if data.SelectedRole == "" {
		data.SelectedRole = domain.OperatorRoleOperator.String()
	}

	data.RoleOptions = make([]operatorRoleOptionView, 0, len(domain.OperatorRoleOptions()))

	for _, role := range domain.OperatorRoleOptions() {
		data.RoleOptions = append(data.RoleOptions, operatorRoleOptionView{
			Value:    role.String(),
			Label:    role.Label(),
			Selected: role.String() == data.SelectedRole,
		})
	}

	return data
}

func humanizeCreateOperatorWebError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyOperatorUsername):
		return "Username is required."
	case errors.Is(err, domain.ErrInvalidOperatorUsername):
		return "Username must be 3-64 characters and contain only lowercase letters, numbers, dots, underscores, or hyphens."
	case errors.Is(err, domain.ErrDuplicateOperatorUsername):
		return "Username is already registered."
	case errors.Is(err, domain.ErrEmptyOperatorRole):
		return "Role is required."
	case errors.Is(err, domain.ErrInvalidOperatorRole):
		return "Role must be either admin or operator."
	case errors.Is(err, domain.ErrEmptyOperatorPassword):
		return "Password is required."
	case errors.Is(err, domain.ErrWeakOperatorPassword):
		return "Password must have at least 15 characters."
	default:
		return "Could not create operator."
	}
}
