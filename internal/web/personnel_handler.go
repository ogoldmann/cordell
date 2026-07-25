package web

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"

	"github.com/go-chi/chi/v5"
)

type personnelFormPageData struct {
	privateLayoutData
	Header                   pageHeaderView
	FormActions              formActionsView
	Feedback                 *feedbackMessageView
	IdentitySection          sectionHeaderView
	AssignmentSection        sectionHeaderView
	ActionURL                string
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
	IdentitySection                sectionHeaderView
	AssignmentSection              sectionHeaderView
	IdentityFields                 []detailFieldView
	AssignmentFields               []detailFieldView
	CurrentCustody                 []currentCustodyView
	HasCurrentCustody              bool
	ShowInactiveCustodyWarning     bool
	HasInactiveAssetCurrentCustody bool
	History                        []custodyHistoryView
	CustodyTimeline                custodyTimelineView
}

type personnelIndexPageData struct {
	privateLayoutData
	Header       pageHeaderView
	EmptyState   emptyStateView
	Title        string
	Personnel    []personnelView
	SearchQuery  string
	StatusFilter string
	StatusTabs   []statusFilterTabView
}

type personnelView struct {
	ID                     string
	FullName               string
	Alias                  string
	Rank                   string
	RankLabel              string
	RegistrationID         string
	Section                string
	SectionLabel           string
	SectionShortLabel      string
	OrganizationUnit       string
	OrganizationUnitLabel  string
	DisplayName            string
	Label                  string
	Active                 bool
	StatusLabel            string
	CanDeactivate          bool
	CanReactivate          bool
	CurrentCustodyQuantity int64
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
	PersonnelLabel  string
	PersonnelURL    string
	Notes           string
	CreatedAt       string
	DateLabel       string
	TimeLabel       string
	Lines           []custodyHistoryLineView
	HasCorrection   bool
	EditCount       int
	EditLabel       string
	EditCountLabel  string
}

type custodyHistoryLineView struct {
	AssetID   string
	AssetName string
	Quantity  int
}

type selectOption struct {
	Value    string
	Label    string
	Selected bool
}

func newPersonnelFormPageData(data personnelFormPageData) personnelFormPageData {
	data.RankOptions = rankOptions()
	data.SectionOptions = sectionOptions()
	data.OrganizationUnitOptions = organizationUnitOptions()

	if data.SelectedOrganizationUnit == "" {
		data.SelectedOrganizationUnit = string(domain.OrganizationUnitDefault)
	}

	return data
}

func (s *Server) handleNewPersonnelForm(w http.ResponseWriter, r *http.Request) {
	var feedback *feedbackMessageView
	if r.URL.Query().Get("created") == "1" {
		feedback = newSuccessFeedback("Militar cadastrado. Cadastre o próximo militar.")
	}

	data := newPersonnelCreatePageData(r, personnelFormPageData{}, feedback)

	if err := s.renderer.Render(w, http.StatusOK, "personnel_form.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreatePersonnel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		data := newPersonnelCreatePageDataFromMessage(r, "Envio de formulário inválido.")
		if renderErr := s.renderer.Render(w, http.StatusBadRequest, "personnel_form.html", data); renderErr != nil {
			s.handleRenderError(w, renderErr)
		}
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
			err,
			r,
		)
		return
	}

	if wantsSaveAndCreateAnother(r) {
		http.Redirect(w, r, "/personnel/new?created=1", http.StatusSeeOther)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/personnel/%s", personnel.ID()),
		http.StatusSeeOther,
	)
}

func (s *Server) handleEditPersonnelForm(w http.ResponseWriter, r *http.Request) {
	id := domain.PersonnelID(chi.URLParam(r, "id"))

	personnel, err := s.services.GetPersonnel.Execute(r.Context(), app.GetPersonnelCommand{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyPersonnelID) {
			s.renderNotFound(w, r)
			return
		}

		s.logger.Error("failed to get personnel for edit", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := newPersonnelEditPageData(r, newPersonnelView(personnel), nil)

	if err := s.renderer.Render(w, http.StatusOK, "personnel_form.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleUpdatePersonnel(w http.ResponseWriter, r *http.Request) {
	id := domain.PersonnelID(chi.URLParam(r, "id"))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	personnel, err := s.services.UpdatePersonnel.Execute(r.Context(), app.UpdatePersonnelCommand{
		ID:               id,
		FullName:         r.FormValue("full_name"),
		Alias:            r.FormValue("alias"),
		Rank:             domain.PersonnelRank(r.FormValue("rank")),
		RegistrationID:   domain.RegistrationID(r.FormValue("registration_id")),
		Section:          domain.PersonnelSection(r.FormValue("section")),
		OrganizationUnit: domain.OrganizationUnit(r.FormValue("organization_unit")),
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyPersonnelID) {
			s.renderNotFound(w, r)
			return
		}

		data := newPersonnelEditPageDataFromRequest(r, id, err)

		if renderErr := s.renderer.Render(w, http.StatusUnprocessableEntity, "personnel_form.html", data); renderErr != nil {
			s.handleRenderError(w, renderErr)
		}

		return
	}

	s.recordAuditEventOrLog(
		r,
		domain.AuditEventPersonnelUpdated,
		domain.AuditEntityPersonnel,
		string(personnel.ID()),
		map[string]string{
			"personnel_id": string(personnel.ID()),
		},
	)

	http.Redirect(w, r, "/personnel/"+string(personnel.ID()), http.StatusSeeOther)
}

func (s *Server) handleShowPersonnel(w http.ResponseWriter, r *http.Request) {
	id := domain.PersonnelID(chi.URLParam(r, "id"))

	personnel, err := s.services.GetPersonnel.Execute(r.Context(), app.GetPersonnelCommand{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			s.renderNotFound(w, r)
			return
		}

		if errors.Is(err, domain.ErrEmptyPersonnelID) {
			s.renderNotFound(w, r)
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

	view := newPersonnelView(personnel)

	data := personnelShowPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             personnel.FullName(),
		Personnel:         view,
		IdentitySection: newSectionHeader(
			"Identificação",
			"Identificação do militar",
			"Dados principais do registro do militar.",
		),
		AssignmentSection: newSectionHeader(
			"Classificação",
			"Dados organizacionais",
			"Posto/graduação, seção e organização do militar.",
		),
		IdentityFields: []detailFieldView{
			newDetailField("Nome completo", view.FullName),
			newDetailField("Nome de guerra", view.DisplayName),
			newDetailField("Identidade", view.RegistrationID),
			newDetailField("Situação", view.StatusLabel),
		},
		AssignmentFields: []detailFieldView{
			newDetailField("Posto/Graduação", view.RankLabel),
			newDetailField("Seção", view.SectionLabel),
			newDetailField("Organização", view.OrganizationUnitLabel),
		},
		CurrentCustody: make([]currentCustodyView, 0, len(currentCustody)),
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		personnelBreadcrumb(),
		currentBreadcrumb(data.Personnel.DisplayName),
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
			PersonnelLabel:  view.DisplayName,
			PersonnelURL:    "/personnel/" + view.ID,
			Notes:           entry.Notes,
			CreatedAt:       entry.CreatedAt.Local().Format("02/01/2006 15:04"),
			DateLabel:       formatTimelineDate(entry.CreatedAt),
			TimeLabel:       formatTimelineTime(entry.CreatedAt),
			Lines:           make([]custodyHistoryLineView, 0, len(entry.Lines)),
			HasCorrection:   entry.HasCorrection,
			EditCount:       entry.EditCount,
			EditLabel:       editCountLabel(entry.EditCount),
			EditCountLabel:  custodyEditCountLabel(entry.EditCount),
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

	data.CustodyTimeline = custodyTimelineView{
		Items: newCustodyTimelineItemsFromPersonnelHistoryItems(data.History),
		EmptyState: newEmptyState(
			"Nenhuma movimentação de custódia",
			"Este militar ainda não possui cautelas ou descautelas registradas.",
			newPageAction("Registrar cautela", "/custody/checkouts/new"),
		),
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_show.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newCustodyTimelineItemsFromPersonnelHistoryItems(
	items []custodyHistoryView,
) []custodyTimelineItemView {
	timelineItems := make([]custodyTimelineItemView, 0, len(items))

	for _, item := range items {
		timelineItems = append(timelineItems, newCustodyTimelineItemFromPersonnelHistoryItem(item))
	}

	return timelineItems
}

func newCustodyTimelineItemFromPersonnelHistoryItem(item custodyHistoryView) custodyTimelineItemView {
	lines := make([]custodyTimelineLineView, 0, len(item.Lines))
	for _, line := range item.Lines {
		lines = append(lines, custodyTimelineLineView{
			AssetID:   line.AssetID,
			AssetName: line.AssetName,
			AssetURL:  "/assets/" + line.AssetID,
			Quantity:  formatTimelineQuantity(line.Quantity),
		})
	}

	receiptURL := "/custody/transactions/" + item.ID

	return custodyTimelineItemView{
		ID:                   item.ID,
		URL:                  receiptURL,
		TypeLabel:            item.TypeLabel,
		TypeTone:             custodyTimelineTypeTone(item.TypeLabel),
		DateLabel:            item.DateLabel,
		TimeLabel:            item.TimeLabel,
		PersonnelLabel:       item.PersonnelLabel,
		PersonnelURL:         item.PersonnelURL,
		OperatorLabel:        item.OperatorDisplay,
		Edited:               item.HasCorrection,
		EditCountLabel:       item.EditCountLabel,
		Notes:                item.Notes,
		Lines:                lines,
		PrimaryActionLabel:   "Abrir recibo",
		PrimaryActionURL:     receiptURL,
		SecondaryActionLabel: "Editar",
		SecondaryActionURL:   receiptURL + "/edit",
	}
}

func (s *Server) handleDeactivatePersonnel(w http.ResponseWriter, r *http.Request) {
	personnelID := domain.PersonnelID(chi.URLParam(r, "id"))

	err := s.services.DeactivatePersonnel.Execute(r.Context(), app.DeactivatePersonnelCommand{
		ID: personnelID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyPersonnelID) {
			s.renderNotFound(w, r)
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
			s.renderNotFound(w, r)
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
	err error,
	r *http.Request,
) {
	data := newPersonnelCreatePageData(r, personnelFormPageData{
		FullName:                 r.FormValue("full_name"),
		Alias:                    r.FormValue("alias"),
		RegistrationID:           r.FormValue("registration_id"),
		SelectedRank:             r.FormValue("rank"),
		SelectedSection:          r.FormValue("section"),
		SelectedOrganizationUnit: r.FormValue("organization_unit"),
	}, personnelFeedbackFromError(err))

	if err := s.renderer.Render(w, status, "personnel_form.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newPersonnelCreatePageDataFromMessage(r *http.Request, message string) personnelFormPageData {
	data := newPersonnelCreatePageData(r, personnelFormPageData{
		FullName:                 r.FormValue("full_name"),
		Alias:                    r.FormValue("alias"),
		RegistrationID:           r.FormValue("registration_id"),
		SelectedRank:             r.FormValue("rank"),
		SelectedSection:          r.FormValue("section"),
		SelectedOrganizationUnit: r.FormValue("organization_unit"),
	}, newErrorFeedback(message))
	data.Error = message

	return data
}

func newPersonnelCreatePageData(
	r *http.Request,
	form personnelFormPageData,
	feedback *feedbackMessageView,
) personnelFormPageData {
	errorMessage := ""
	if feedback != nil && feedback.Kind == "error" {
		errorMessage = feedback.Message
	}

	form.privateLayoutData = newPrivateLayoutData(r)
	form.Title = "Cadastrar militar"
	form.Header = newPageHeader(
		personnelPluralLabel(),
		"Cadastrar militar",
		"Cadastre um militar que poderá receber materiais sob custódia.",
		nil,
	)
	form.FormActions = newFormActionsWithSaveAndCreateAnother(
		"Cadastrar militar",
		"Cancelar",
		"/personnel",
		"Salvar e cadastrar outro",
	)
	form.Feedback = feedback
	form.IdentitySection = newSectionHeader(
		"Identificação",
		"Identificação do militar",
		"Informe os dados principais do militar.",
	)
	form.AssignmentSection = newSectionHeader(
		"Classificação",
		"Dados organizacionais",
		"Informe posto/graduação, seção e organização.",
	)
	form.ActionURL = "/personnel"
	form.Error = errorMessage
	form.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		personnelBreadcrumb(),
		currentBreadcrumb("Cadastrar militar"),
	}

	return newPersonnelFormPageData(form)
}

func newPersonnelEditPageData(
	r *http.Request,
	personnel personnelView,
	err error,
) personnelFormPageData {
	errorMessage := ""
	if err != nil {
		errorMessage = duplicatePersonnelErrorMessage(err)
	}

	data := personnelFormPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Editar militar",
		Header: newPageHeader(
			personnelPluralLabel(),
			"Editar militar",
			"Atualize os dados cadastrais do militar.",
			nil,
		),
		FormActions: newFormActions(
			"Salvar alterações",
			"Cancelar",
			"/personnel/"+personnel.ID,
		),
		Feedback: personnelFeedbackFromError(err),
		IdentitySection: newSectionHeader(
			"Identificação",
			"Identificação do militar",
			"Atualize os dados principais do militar.",
		),
		AssignmentSection: newSectionHeader(
			"Classificação",
			"Dados organizacionais",
			"Atualize posto/graduação, seção e organização.",
		),
		ActionURL:                "/personnel/" + personnel.ID,
		Error:                    errorMessage,
		FullName:                 personnel.FullName,
		Alias:                    personnel.Alias,
		RegistrationID:           personnel.RegistrationID,
		SelectedRank:             personnel.Rank,
		SelectedSection:          personnel.Section,
		SelectedOrganizationUnit: personnel.OrganizationUnit,
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		personnelBreadcrumb(),
		{
			Label: personnel.DisplayName,
			URL:   "/personnel/" + personnel.ID,
		},
		currentBreadcrumb("Editar"),
	}

	return newPersonnelFormPageData(data)
}

func newPersonnelEditPageDataFromRequest(
	r *http.Request,
	id domain.PersonnelID,
	err error,
) personnelFormPageData {
	errorMessage := ""
	if err != nil {
		errorMessage = duplicatePersonnelErrorMessage(err)
	}

	data := personnelFormPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Editar militar",
		Header: newPageHeader(
			personnelPluralLabel(),
			"Editar militar",
			"Atualize os dados cadastrais do militar.",
			nil,
		),
		FormActions: newFormActions(
			"Salvar alterações",
			"Cancelar",
			"/personnel/"+string(id),
		),
		Feedback: personnelFeedbackFromError(err),
		IdentitySection: newSectionHeader(
			"Identificação",
			"Identificação do militar",
			"Atualize os dados principais do militar.",
		),
		AssignmentSection: newSectionHeader(
			"Classificação",
			"Dados organizacionais",
			"Atualize posto/graduação, seção e organização.",
		),
		ActionURL:                "/personnel/" + string(id),
		Error:                    errorMessage,
		FullName:                 r.FormValue("full_name"),
		Alias:                    r.FormValue("alias"),
		RegistrationID:           r.FormValue("registration_id"),
		SelectedRank:             r.FormValue("rank"),
		SelectedSection:          r.FormValue("section"),
		SelectedOrganizationUnit: r.FormValue("organization_unit"),
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		personnelBreadcrumb(),
		{
			Label: string(id),
			URL:   "/personnel/" + string(id),
		},
		currentBreadcrumb("Editar"),
	}

	return newPersonnelFormPageData(data)
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
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	var personnel []domain.Personnel
	var err error

	if searchQuery == "" {
		personnel, err = s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
			Limit:        100,
			StatusFilter: string(statusFilter),
		})
	} else {
		personnel, err = s.services.SearchPersonnel.Execute(r.Context(), app.SearchPersonnelCommand{
			Query:        searchQuery,
			Limit:        100,
			StatusFilter: string(statusFilter),
		})
	}
	if err != nil {
		s.logger.Error("failed to list personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	currentCustodySummaries, err := s.services.ListCurrentCustodySummary.Execute(r.Context())
	if err != nil {
		s.logger.Error("failed to list current custody summaries by personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	currentCustodyByPersonnelID := make(map[string]int64, len(currentCustodySummaries))
	for _, summary := range currentCustodySummaries {
		currentCustodyByPersonnelID[string(summary.PersonnelID)] = summary.CurrentCustodyQuantity
	}

	data := personnelIndexPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Header: newPageHeader(
			"Militares",
			personnelPluralLabel(),
			"Militares que podem cautelar e descautelar materiais.",
			newPageAction("Cadastrar Militar", "/personnel/new"),
		),
		EmptyState: newEmptyState(
			"Nenhum militar cadastrado",
			"Cadastre o primeiro militar para começar a usar o Cordell.",
			newPageAction("Cadastrar militar", "/personnel/new"),
		),
		Title:        personnelPluralLabel(),
		Personnel:    make([]personnelView, 0, len(personnel)),
		SearchQuery:  searchQuery,
		StatusFilter: string(statusFilter),
		StatusTabs:   newStatusFilterTabs("/personnel", statusFilter, searchQuery),
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		currentBreadcrumb(personnelPluralLabel()),
	}

	for _, item := range personnel {
		view := newPersonnelView(item)
		view.CurrentCustodyQuantity = currentCustodyByPersonnelID[view.ID]
		data.Personnel = append(data.Personnel, view)
	}

	data.UseDefaultShell = false

	if wantsPartialResponse(r) {
		if err := s.renderer.Render(w, http.StatusOK, "personnel_table", data); err != nil {
			s.logger.Error("failed to render personnel list partial", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		return
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newPersonnelView(personnel domain.Personnel) personnelView {
	statusLabel := activeStatusLabel(personnel.Active())
	displayName := militaryDisplayName(personnel.Rank(), personnel.Alias())
	sectionShortLabel := personnel.Section().Abbreviation()

	return personnelView{
		ID:                    string(personnel.ID()),
		FullName:              personnel.FullName(),
		Alias:                 personnel.Alias(),
		Rank:                  string(personnel.Rank()),
		RankLabel:             personnelRankLabel(personnel.Rank()),
		RegistrationID:        personnel.RegistrationID().String(),
		Section:               string(personnel.Section()),
		SectionLabel:          personnelSectionLabel(personnel.Section()),
		SectionShortLabel:     sectionShortLabel,
		OrganizationUnit:      string(personnel.OrganizationUnit()),
		OrganizationUnitLabel: organizationUnitLabel(personnel.OrganizationUnit()),
		DisplayName:           displayName,
		Label:                 newPersonnelOptionLabel(displayName, personnel.FullName(), sectionShortLabel),
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
			Label: option.Value.DisplayLabel(),
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
	return section.DisplayLabel()
}

func organizationUnitLabel(unit domain.OrganizationUnit) string {
	for _, option := range domain.OrganizationUnitOptions() {
		if option.Value == unit {
			return option.Label
		}
	}

	return string(unit)
}
