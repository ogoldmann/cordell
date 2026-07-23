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
	privateLayoutData
	Title                    string
	Error                    string
	FullName                 string
	Alias                    string
	RegistrationID           string
	SelectedRank             string
	SelectedSection          string
	SelectedOrganizationUnit string
	RankOptions              []selectOption
	SectionOptions           []selectOption
	OrganizationUnitOptions  []selectOption
}

type personnelShowPageData struct {
	privateLayoutData
	Title                          string
	Personnel                      personnelView
	CurrentCustody                 []currentCustodyView
	HasCurrentCustody              bool
	ShowInactiveCustodyWarning     bool
	HasInactiveAssetCurrentCustody bool
	History                        []custodyHistoryView
}

type personnelIndexPageData struct {
	privateLayoutData
	Title        string
	Personnel    []personnelView
	StatusFilter string
	StatusTabs   []statusFilterTabView
}

type personnelView struct {
	ID                    string
	FullName              string
	Alias                 string
	Rank                  string
	RankLabel             string
	RegistrationID        string
	Section               string
	SectionLabel          string
	OrganizationUnit      string
	OrganizationUnitLabel string
	DisplayName           string
	Active                bool
	StatusLabel           string
	CanDeactivate         bool
	CanReactivate         bool
}

type currentCustodyView struct {
	AssetID     string
	AssetName   string
	AssetActive bool
	StatusLabel string
	Quantity    int
}

type custodyHistoryView struct {
	ID              string
	Type            string
	TypeLabel       string
	OperatorID      string
	OperatorDisplay string
	Notes           string
	CreatedAt       string
	Lines           []custodyHistoryLineView
	HasCorrection   bool
	EditCount       int
	EditLabel       string
}

type custodyHistoryLineView struct {
	AssetID   string
	AssetName string
	Quantity  int
}

type selectOption struct {
	Value string
	Label string
}

func newPersonnelNewPageData(data personnelNewPageData) personnelNewPageData {
	data.RankOptions = rankOptions()
	data.SectionOptions = sectionOptions()
	data.OrganizationUnitOptions = organizationUnitOptions()

	if data.SelectedOrganizationUnit == "" {
		data.SelectedOrganizationUnit = string(domain.OrganizationUnitDefault)
	}

	return data
}

func (s *Server) handleNewPersonnelForm(w http.ResponseWriter, r *http.Request) {
	data := newPersonnelNewPageData(personnelNewPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Cadastrar militar",
	})

	if err := s.renderer.Render(w, http.StatusOK, "personnel_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreatePersonnel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderNewPersonnelFormWithError(
			w,
			http.StatusBadRequest,
			"Envio de formulário inválido.",
			r,
		)
		return
	}

	fullName := r.FormValue("full_name")
	alias := r.FormValue("alias")
	registrationID := r.FormValue("registration_id")
	rank := r.FormValue("rank")
	section := r.FormValue("section")
	organizationUnit := r.FormValue("organization_unit")

	personnel, err := s.services.CreatePersonnel.Execute(r.Context(), app.CreatePersonnelCommand{
		FullName:         fullName,
		Alias:            alias,
		Rank:             domain.PersonnelRank(rank),
		RegistrationID:   registrationID,
		Section:          domain.PersonnelSection(section),
		OrganizationUnit: domain.OrganizationUnit(organizationUnit),
	})
	if err != nil {
		s.renderNewPersonnelFormWithError(
			w,
			http.StatusBadRequest,
			humanizePersonnelError(err),
			r,
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

	currentCustody, err := s.services.ListCurrentCustody.Execute(r.Context(), app.ListCurrentCustodyCommand{
		PersonnelID: personnel.ID(),
	})
	if err != nil {
		s.logger.Error("failed to list current custody", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	history, err := s.services.ListCustodyHistory.Execute(r.Context(), app.ListCustodyHistoryCommand{
		PersonnelID: personnel.ID(),
		Limit:       50,
	})
	if err != nil {
		s.logger.Error("failed to list custody history", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := personnelShowPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             personnel.FullName(),
		Personnel:         newPersonnelView(personnel),
		CurrentCustody:    make([]currentCustodyView, 0, len(currentCustody)),
	}

	for _, item := range currentCustody {
		assetStatusLabel := activeStatusLabel(item.AssetActive)

		if !item.AssetActive {
			data.HasInactiveAssetCurrentCustody = true
		}

		data.CurrentCustody = append(data.CurrentCustody, currentCustodyView{
			AssetID:     string(item.AssetID),
			AssetName:   item.AssetName,
			AssetActive: item.AssetActive,
			StatusLabel: assetStatusLabel,
			Quantity:    item.Quantity,
		})
	}

	data.HasCurrentCustody = len(data.CurrentCustody) > 0
	data.ShowInactiveCustodyWarning = !personnel.Active() && data.HasCurrentCustody

	for _, entry := range history {
		historyView := custodyHistoryView{
			ID:              string(entry.ID),
			Type:            string(entry.Type),
			TypeLabel:       custodyTransactionTypeLabel(entry.Type),
			OperatorID:      string(entry.OperatorID),
			OperatorDisplay: militaryDisplayName(entry.OperatorRank, entry.OperatorAlias),
			Notes:           entry.Notes,
			CreatedAt:       entry.CreatedAt.Local().Format("02/01/2006 15:04"),
			Lines:           make([]custodyHistoryLineView, 0, len(entry.Lines)),
			HasCorrection:   entry.HasCorrection,
			EditCount:       entry.EditCount,
			EditLabel:       editCountLabel(entry.EditCount),
		}

		for _, line := range entry.Lines {
			historyView.Lines = append(historyView.Lines, custodyHistoryLineView{
				AssetID:   string(line.AssetID),
				AssetName: line.AssetName,
				Quantity:  line.Quantity,
			})
		}

		data.History = append(data.History, historyView)
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_show.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleDeactivatePersonnel(w http.ResponseWriter, r *http.Request) {
	personnelID := domain.PersonnelID(chi.URLParam(r, "id"))

	err := s.services.DeactivatePersonnel.Execute(r.Context(), app.DeactivatePersonnelCommand{
		ID: personnelID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyPersonnelID) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to deactivate personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.recordAuditEventOrLog(
		r,
		domain.AuditEventPersonnelDeactivated,
		domain.AuditEntityPersonnel,
		string(personnelID),
		nil,
	)

	http.Redirect(w, r, "/personnel/"+string(personnelID), http.StatusSeeOther)
}

func (s *Server) handleReactivatePersonnel(w http.ResponseWriter, r *http.Request) {
	personnelID := domain.PersonnelID(chi.URLParam(r, "id"))

	err := s.services.ReactivatePersonnel.Execute(r.Context(), app.ReactivatePersonnelCommand{
		ID: personnelID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyPersonnelID) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to reactivate personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.recordAuditEventOrLog(
		r,
		domain.AuditEventPersonnelReactivated,
		domain.AuditEntityPersonnel,
		string(personnelID),
		nil,
	)

	http.Redirect(w, r, "/personnel/"+string(personnelID), http.StatusSeeOther)
}

func (s *Server) renderNewPersonnelFormWithError(
	w http.ResponseWriter,
	status int,
	message string,
	r *http.Request,
) {
	data := newPersonnelNewPageData(personnelNewPageData{
		privateLayoutData:        newPrivateLayoutData(r),
		Title:                    "Cadastrar militar",
		Error:                    message,
		FullName:                 r.FormValue("full_name"),
		Alias:                    r.FormValue("alias"),
		RegistrationID:           r.FormValue("registration_id"),
		SelectedRank:             r.FormValue("rank"),
		SelectedSection:          r.FormValue("section"),
		SelectedOrganizationUnit: r.FormValue("organization_unit"),
	})

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
		return "O nome completo é obrigatório."
	case errors.Is(err, domain.ErrEmptyPersonnelAlias):
		return "O nome de guerra é obrigatório."
	case errors.Is(err, domain.ErrEmptyRegistrationID):
		return "A identidade é obrigatória."
	case errors.Is(err, domain.ErrInvalidRegistrationID):
		return "A identidade informada é inválida."
	case errors.Is(err, domain.ErrDuplicateRegistrationID):
		return "Esta identidade já está cadastrada."
	case errors.Is(err, domain.ErrInvalidPersonnelRank):
		return "O posto ou graduação é obrigatório."
	case errors.Is(err, domain.ErrInvalidPersonnelSection):
		return "A seção é obrigatória."
	case errors.Is(err, domain.ErrInvalidOrganizationUnit):
		return "A organização é obrigatória."
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "O militar é obrigatório."
	default:
		return "Não foi possível cadastrar o militar."
	}
}

func (s *Server) handleListPersonnel(w http.ResponseWriter, r *http.Request) {
	statusFilter := ports.NormalizeRecordStatusFilter(r.URL.Query().Get("status"))

	personnel, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit:        100,
		StatusFilter: string(statusFilter),
	})
	if err != nil {
		s.logger.Error("failed to list personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := personnelIndexPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             personnelPluralLabel(),
		Personnel:         make([]personnelView, 0, len(personnel)),
		StatusFilter:      string(statusFilter),
		StatusTabs:        newStatusFilterTabs("/personnel", statusFilter),
	}

	for _, item := range personnel {
		data.Personnel = append(data.Personnel, newPersonnelView(item))
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newPersonnelView(personnel domain.Personnel) personnelView {
	statusLabel := activeStatusLabel(personnel.Active())

	return personnelView{
		ID:                    string(personnel.ID()),
		FullName:              personnel.FullName(),
		Alias:                 personnel.Alias(),
		Rank:                  string(personnel.Rank()),
		RankLabel:             personnelRankLabel(personnel.Rank()),
		RegistrationID:        personnel.RegistrationID().String(),
		Section:               string(personnel.Section()),
		SectionLabel:          personnelSectionLabel(personnel.Section()),
		OrganizationUnit:      string(personnel.OrganizationUnit()),
		OrganizationUnitLabel: organizationUnitLabel(personnel.OrganizationUnit()),
		DisplayName:           militaryDisplayName(personnel.Rank(), personnel.Alias()),
		Active:                personnel.Active(),
		StatusLabel:           statusLabel,
		CanDeactivate:         personnel.Active(),
		CanReactivate:         !personnel.Active(),
	}
}

func rankOptions() []selectOption {
	options := domain.PersonnelRankOptions()
	result := make([]selectOption, 0, len(options))

	for _, option := range options {
		result = append(result, selectOption{
			Value: string(option.Value),
			Label: option.Value.DisplayLabel(),
		})
	}

	return result
}

func sectionOptions() []selectOption {
	options := domain.PersonnelSectionOptions()
	result := make([]selectOption, 0, len(options))

	for _, option := range options {
		result = append(result, selectOption{
			Value: string(option.Value),
			Label: option.Label,
		})
	}

	return result
}

func organizationUnitOptions() []selectOption {
	options := domain.OrganizationUnitOptions()
	result := make([]selectOption, 0, len(options))

	for _, option := range options {
		result = append(result, selectOption{
			Value: string(option.Value),
			Label: option.Label,
		})
	}

	return result
}

func personnelRankLabel(rank domain.PersonnelRank) string {
	return rank.DisplayLabel()
}

func personnelSectionLabel(section domain.PersonnelSection) string {
	for _, option := range domain.PersonnelSectionOptions() {
		if option.Value == section {
			return option.Label
		}
	}

	return string(section)
}

func organizationUnitLabel(unit domain.OrganizationUnit) string {
	for _, option := range domain.OrganizationUnitOptions() {
		if option.Value == unit {
			return option.Label
		}
	}

	return string(unit)
}
