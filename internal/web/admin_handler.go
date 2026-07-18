package web

import (
	"errors"
	"net/http"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"

	"github.com/go-chi/chi/v5"
)

type adminIndexPageData struct {
	privateLayoutData
	Title string
}

type adminOperatorsIndexPageData struct {
	privateLayoutData
	Title     string
	Error     string
	Operators []operatorSummaryView
}

type adminOperatorShowPageData struct {
	privateLayoutData
	Title    string
	Error    string
	Operator operatorDetailView
}

type adminOperatorNewPageData struct {
	privateLayoutData
	Title          string
	Error          string
	RegistrationID string
	Alias          string
	SelectedRank   string
	RankOptions    []operatorRankOptionView
	SelectedRole   string
	RoleOptions    []operatorRoleOptionView
}

type operatorRankOptionView struct {
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
	s.renderAdminOperatorsIndex(w, r, http.StatusOK, "")
}

func (s *Server) renderAdminOperatorsIndex(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message string,
) {
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
		Error:             message,
		Operators:         make([]operatorSummaryView, 0, len(operators)),
	}

	for _, operator := range operators {
		data.Operators = append(data.Operators, newOperatorSummaryView(operator))
	}

	if err := s.renderer.Render(w, status, "admin_operators_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleShowAdminOperator(w http.ResponseWriter, r *http.Request) {
	operatorID := domain.OperatorID(chi.URLParam(r, "id"))

	s.renderAdminOperatorShow(w, r, http.StatusOK, operatorID, "")
}

func (s *Server) renderAdminOperatorShow(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	operatorID domain.OperatorID,
	message string,
) {
	operator, err := s.services.GetOperatorAdmin.Execute(r.Context(), app.GetOperatorAdminCommand{
		OperatorID: operatorID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to get operator admin detail", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := adminOperatorShowPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             operatorDisplayName(operator.Rank, operator.Alias),
		Error:             message,
		Operator:          newOperatorDetailView(operator, currentOperator.ID()),
	}

	if err := s.renderer.Render(w, status, "admin_operator_show.html", data); err != nil {
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

	registrationID := r.FormValue("registration_id")
	alias := r.FormValue("alias")
	rank := r.FormValue("rank")
	role := r.FormValue("role")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		s.renderAdminOperatorFormWithError(
			w,
			r,
			http.StatusBadRequest,
			"Passwords do not match.",
			registrationID,
			alias,
			rank,
			role,
		)
		return
	}

	createdOperator, err := s.services.CreateOperator.Execute(r.Context(), app.CreateOperatorCommand{
		RegistrationID: registrationID,
		Alias:          alias,
		Rank:           rank,
		Role:           role,
		Password:       password,
	})
	if err != nil {
		s.renderAdminOperatorFormWithError(
			w,
			r,
			http.StatusBadRequest,
			humanizeCreateOperatorWebError(err),
			registrationID,
			alias,
			rank,
			role,
		)
		return
	}

	if err := s.recordAuditEvent(
		r,
		domain.AuditEventOperatorCreated,
		domain.AuditEntityOperator,
		string(createdOperator.ID()),
		map[string]string{
			"role": createdOperator.Role().String(),
			"rank": createdOperator.Rank().String(),
		},
	); err != nil {
		s.logger.Error("failed to record audit event", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/operators", http.StatusSeeOther)
}

func (s *Server) renderAdminOperatorFormWithError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message string,
	registrationID string,
	alias string,
	rank string,
	role string,
) {
	data := newAdminOperatorNewPageData(r, adminOperatorNewPageData{
		Title:          "Create operator",
		Error:          message,
		RegistrationID: registrationID,
		Alias:          alias,
		SelectedRank:   rank,
		SelectedRole:   role,
	})

	if err := s.renderer.Render(w, status, "admin_operator_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleDeactivateAdminOperator(w http.ResponseWriter, r *http.Request) {
	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	operatorID := domain.OperatorID(chi.URLParam(r, "id"))

	err := s.services.DeactivateOperator.Execute(r.Context(), app.DeactivateOperatorCommand{
		CurrentOperatorID: currentOperator.ID(),
		OperatorID:        operatorID,
	})
	if err != nil {
		s.renderAdminOperatorShow(
			w,
			r,
			http.StatusBadRequest,
			operatorID,
			humanizeDeactivateOperatorWebError(err),
		)
		return
	}

	if err := s.recordAuditEvent(
		r,
		domain.AuditEventOperatorDeactivated,
		domain.AuditEntityOperator,
		string(operatorID),
		nil,
	); err != nil {
		s.logger.Error("failed to record audit event", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/operators/"+string(operatorID), http.StatusSeeOther)
}

func (s *Server) handleChangeAdminOperatorRole(w http.ResponseWriter, r *http.Request) {
	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	operatorID := domain.OperatorID(chi.URLParam(r, "id"))
	role := r.FormValue("role")

	err := s.services.ChangeOperatorRole.Execute(r.Context(), app.ChangeOperatorRoleCommand{
		CurrentOperatorID: currentOperator.ID(),
		OperatorID:        operatorID,
		Role:              role,
	})
	if err != nil {
		s.renderAdminOperatorShow(
			w,
			r,
			http.StatusBadRequest,
			operatorID,
			humanizeChangeOperatorRoleWebError(err),
		)
		return
	}

	if err := s.recordAuditEvent(
		r,
		domain.AuditEventOperatorRoleChanged,
		domain.AuditEntityOperator,
		string(operatorID),
		map[string]string{
			"new_role": role,
		},
	); err != nil {
		s.logger.Error("failed to record audit event", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/operators/"+string(operatorID), http.StatusSeeOther)
}

func (s *Server) handleResetAdminOperatorPassword(w http.ResponseWriter, r *http.Request) {
	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	operatorID := domain.OperatorID(chi.URLParam(r, "id"))
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		s.renderAdminOperatorShow(
			w,
			r,
			http.StatusBadRequest,
			operatorID,
			"Passwords do not match.",
		)
		return
	}

	err := s.services.ResetOperatorPassword.Execute(r.Context(), app.ResetOperatorPasswordCommand{
		CurrentOperatorID: currentOperator.ID(),
		OperatorID:        operatorID,
		Password:          password,
	})
	if err != nil {
		s.renderAdminOperatorShow(
			w,
			r,
			http.StatusBadRequest,
			operatorID,
			humanizeResetOperatorPasswordWebError(err),
		)
		return
	}

	if err := s.recordAuditEvent(
		r,
		domain.AuditEventOperatorPasswordReset,
		domain.AuditEntityOperator,
		string(operatorID),
		nil,
	); err != nil {
		s.logger.Error("failed to record audit event", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/operators/"+string(operatorID), http.StatusSeeOther)
}

func (s *Server) handleReactivateAdminOperator(w http.ResponseWriter, r *http.Request) {
	operatorID := domain.OperatorID(chi.URLParam(r, "id"))

	err := s.services.ReactivateOperator.Execute(r.Context(), app.ReactivateOperatorCommand{
		OperatorID: operatorID,
	})
	if err != nil {
		s.renderAdminOperatorShow(
			w,
			r,
			http.StatusBadRequest,
			operatorID,
			humanizeReactivateOperatorWebError(err),
		)
		return
	}

	if err := s.recordAuditEvent(
		r,
		domain.AuditEventOperatorReactivated,
		domain.AuditEntityOperator,
		string(operatorID),
		nil,
	); err != nil {
		s.logger.Error("failed to record audit event", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/operators/"+string(operatorID), http.StatusSeeOther)
}

func newAdminOperatorNewPageData(
	r *http.Request,
	data adminOperatorNewPageData,
) adminOperatorNewPageData {
	data.privateLayoutData = newPrivateLayoutData(r)

	if data.SelectedRole == "" {
		data.SelectedRole = domain.OperatorRoleOperator.String()
	}

	data.RankOptions = newOperatorRankOptionViews(data.SelectedRank)
	data.RoleOptions = newOperatorRoleOptionViews(data.SelectedRole)

	return data
}

func newOperatorRankOptionViews(selectedRank string) []operatorRankOptionView {
	options := make([]operatorRankOptionView, 0, len(domain.RankOptions()))

	for _, rank := range domain.RankOptions() {
		options = append(options, operatorRankOptionView{
			Value:    rank.Value.String(),
			Label:    rank.Label,
			Selected: rank.Value.String() == selectedRank,
		})
	}

	return options
}

func humanizeCreateOperatorWebError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyRegistrationID):
		return "Registration ID is required."
	case errors.Is(err, domain.ErrInvalidRegistrationID):
		return "Registration ID is invalid."
	case errors.Is(err, domain.ErrDuplicateRegistrationID):
		return "Registration ID is already registered."
	case errors.Is(err, domain.ErrEmptyOperatorAlias):
		return "Alias is required."
	case errors.Is(err, domain.ErrInvalidOperatorRank):
		return "Rank is required."
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

func humanizeDeactivateOperatorWebError(err error) string {
	switch {
	case errors.Is(err, domain.ErrCannotDeactivateCurrentOperator):
		return "You cannot deactivate your own operator account."
	case errors.Is(err, domain.ErrCannotDeactivateLastAdmin):
		return "You cannot deactivate the last active admin."
	case errors.Is(err, ports.ErrNotFound):
		return "Operator not found."
	default:
		return "Could not deactivate operator."
	}
}

func humanizeChangeOperatorRoleWebError(err error) string {
	switch {
	case errors.Is(err, domain.ErrCannotChangeCurrentOperatorRole):
		return "You cannot change your own operator role."
	case errors.Is(err, domain.ErrCannotDemoteLastAdmin):
		return "You cannot demote the last active admin."
	case errors.Is(err, domain.ErrEmptyOperatorRole):
		return "Role is required."
	case errors.Is(err, domain.ErrInvalidOperatorRole):
		return "Role must be either admin or operator."
	case errors.Is(err, ports.ErrNotFound):
		return "Operator not found."
	default:
		return "Could not change operator role."
	}
}

func humanizeResetOperatorPasswordWebError(err error) string {
	switch {
	case errors.Is(err, domain.ErrCannotResetCurrentOperatorPassword):
		return "You cannot reset your own password from this admin action."
	case errors.Is(err, domain.ErrEmptyOperatorPassword):
		return "Password is required."
	case errors.Is(err, domain.ErrWeakOperatorPassword):
		return "Password must have at least 15 characters."
	case errors.Is(err, ports.ErrNotFound):
		return "Operator not found or inactive."
	default:
		return "Could not reset operator password."
	}
}

func humanizeReactivateOperatorWebError(err error) string {
	switch {
	case errors.Is(err, ports.ErrNotFound):
		return "Operator not found."
	default:
		return "Could not reactivate operator."
	}
}
