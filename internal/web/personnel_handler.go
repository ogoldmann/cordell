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

type personnelNewPageData struct {
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
	Title          string
	Personnel      personnelView
	CurrentCustody []currentCustodyView
	History        []custodyHistoryView
}

type personnelIndexPageData struct {
	Title       string
	SearchQuery string
	Personnel   []personnelView
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
	Active                bool
}

type currentCustodyView struct {
	AssetID   string
	AssetName string
	Quantity  int
}

type custodyHistoryView struct {
	ID        string
	Type      string
	TypeLabel string
	Notes     string
	CreatedAt string
	Lines     []custodyHistoryLineView
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
		Title: "Create personnel",
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
			"Invalid form submission.",
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
		Title:          personnel.FullName(),
		Personnel:      newPersonnelView(personnel),
		CurrentCustody: make([]currentCustodyView, 0, len(currentCustody)),
	}

	for _, item := range currentCustody {
		data.CurrentCustody = append(data.CurrentCustody, currentCustodyView{
			AssetID:   string(item.AssetID),
			AssetName: item.AssetName,
			Quantity:  item.Quantity,
		})
	}

	for _, entry := range history {
		historyView := custodyHistoryView{
			ID:        string(entry.ID),
			Type:      string(entry.Type),
			TypeLabel: custodyTransactionTypeLabel(entry.Type),
			Notes:     entry.Notes,
			CreatedAt: entry.CreatedAt.Local().Format("2006-01-02 15:04"),
			Lines:     make([]custodyHistoryLineView, 0, len(entry.Lines)),
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

func (s *Server) renderNewPersonnelFormWithError(
	w http.ResponseWriter,
	status int,
	message string,
	r *http.Request,
) {
	data := newPersonnelNewPageData(personnelNewPageData{
		Title:                    "Create personnel",
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
		return "Full name is required."
	case errors.Is(err, domain.ErrEmptyPersonnelAlias):
		return "Alias is required."
	case errors.Is(err, domain.ErrEmptyRegistrationID):
		return "Registration ID is required."
	case errors.Is(err, domain.ErrInvalidRegistrationID):
		return "Registration ID is invalid."
	case errors.Is(err, domain.ErrDuplicateRegistrationID):
		return "Registration ID is already registered."
	case errors.Is(err, domain.ErrInvalidPersonnelRank):
		return "Rank is required."
	case errors.Is(err, domain.ErrInvalidPersonnelSection):
		return "Section is required."
	case errors.Is(err, domain.ErrInvalidOrganizationUnit):
		return "Organization unit is required."
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "Personnel ID is required."
	default:
		return "Could not create personnel."
	}
}

func (s *Server) handleListPersonnel(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	personnel, err := s.services.SearchPersonnel.Execute(r.Context(), app.SearchPersonnelCommand{
		Query: searchQuery,
		Limit: 50,
	})
	if err != nil {
		s.logger.Error("failed to list personnel", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := personnelIndexPageData{
		Title:       "Personnel",
		SearchQuery: searchQuery,
		Personnel:   make([]personnelView, 0, len(personnel)),
	}

	for _, item := range personnel {
		data.Personnel = append(data.Personnel, newPersonnelView(item))
	}

	if err := s.renderer.Render(w, http.StatusOK, "personnel_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func custodyTransactionTypeLabel(transactionType domain.CustodyTransactionType) string {
	switch transactionType {
	case domain.CustodyTransactionTypeCheckout:
		return "Checkout"
	case domain.CustodyTransactionTypeReturn:
		return "Return"
	default:
		return "Unknown"
	}
}

func newPersonnelView(personnel domain.Personnel) personnelView {
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
		Active:                personnel.Active(),
	}
}

func rankOptions() []selectOption {
	options := domain.PersonnelRankOptions()
	result := make([]selectOption, 0, len(options))

	for _, option := range options {
		result = append(result, selectOption{
			Value: string(option.Value),
			Label: option.Label,
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
	for _, option := range domain.PersonnelRankOptions() {
		if option.Value == rank {
			return option.Label
		}
	}

	return string(rank)
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
