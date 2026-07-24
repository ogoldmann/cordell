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

type assetIndexPageData struct {
	privateLayoutData
	Header       pageHeaderView
	EmptyState   emptyStateView
	Title        string
	Assets       []assetView
	SearchQuery  string
	StatusFilter string
	StatusTabs   []statusFilterTabView
}

type assetFormPageData struct {
	privateLayoutData
	Header         pageHeaderView
	FormActions    formActionsView
	Feedback       *feedbackMessageView
	DetailsSection sectionHeaderView
	ActionURL      string
	Title          string
	Error          string
	Form           assetFormView
}

type assetShowPageData struct {
	privateLayoutData
	Title                      string
	Asset                      assetView
	DetailsSection             sectionHeaderView
	DetailsFields              []detailFieldView
	CustodyHistorySection      sectionHeaderView
	CustodyTimeline            custodyTimelineView
	Holders                    []assetHolderView
	HasHolders                 bool
	HasInactiveHolders         bool
	ShowInactiveCustodyWarning bool
}

type assetView struct {
	ID            string
	Name          string
	Active        bool
	StatusLabel   string
	CanDeactivate bool
	CanReactivate bool
}

type assetHolderView struct {
	PersonnelID          string
	PersonnelFullName    string
	PersonnelDisplay     string
	PersonnelActive      bool
	PersonnelStatusLabel string
	Quantity             int
}

type assetFormView struct {
	Name string
}

func (s *Server) handleListAssets(w http.ResponseWriter, r *http.Request) {
	statusFilter := ports.NormalizeRecordStatusFilter(r.URL.Query().Get("status"))
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	var assets []domain.Asset
	var err error

	if searchQuery == "" {
		assets, err = s.services.ListAssets.Execute(r.Context(), app.ListAssetsCommand{
			Limit:        100,
			StatusFilter: string(statusFilter),
		})
	} else {
		assets, err = s.services.SearchAssets.Execute(r.Context(), app.SearchAssetsCommand{
			Query:        searchQuery,
			Limit:        100,
			StatusFilter: string(statusFilter),
		})
	}
	if err != nil {
		s.logger.Error("failed to list assets", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := assetIndexPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Header: newPageHeader(
			"Materiais",
			assetPluralLabel(),
			"Materiais que podem ser cautelados aos militares.",
			newPageAction("Cadastrar material", "/assets/new"),
		),
		EmptyState: newEmptyState(
			"Nenhum material cadastrado",
			"Cadastre o primeiro material para começar a usar o Cordell.",
			newPageAction("Cadastrar material", "/assets/new"),
		),
		Title:        assetPluralLabel(),
		Assets:       make([]assetView, 0, len(assets)),
		SearchQuery:  searchQuery,
		StatusFilter: string(statusFilter),
		StatusTabs:   newStatusFilterTabs("/assets", statusFilter, searchQuery),
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		currentBreadcrumb(assetPluralLabel()),
	}

	for _, item := range assets {
		data.Assets = append(data.Assets, newAssetView(item))
	}

	if wantsPartialResponse(r) {
		if err := s.renderer.Render(w, http.StatusOK, "asset_list", data); err != nil {
			s.logger.Error("failed to render asset list partial", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		return
	}

	if err := s.renderer.Render(w, http.StatusOK, "assets_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleNewAssetForm(w http.ResponseWriter, r *http.Request) {
	var feedback *feedbackMessageView
	if r.URL.Query().Get("created") == "1" {
		feedback = newSuccessFeedback("Material cadastrado. Cadastre o próximo material.")
	}

	data := newAssetCreatePageData(r, assetFormView{}, feedback)

	if err := s.renderer.Render(w, http.StatusOK, "asset_form.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateAsset(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		data := newAssetCreatePageDataFromMessage(r, "Envio de formulário inválido.")
		if renderErr := s.renderer.Render(w, http.StatusBadRequest, "asset_form.html", data); renderErr != nil {
			s.handleRenderError(w, renderErr)
		}
		return
	}

	name := r.FormValue("name")

	asset, err := s.services.CreateAsset.Execute(r.Context(), app.CreateAssetCommand{
		Name: name,
	})
	if err != nil {
		s.renderNewAssetFormWithError(
			w,
			r,
			http.StatusBadRequest,
			err,
			name,
		)
		return
	}

	if wantsSaveAndCreateAnother(r) {
		http.Redirect(w, r, "/assets/new?created=1", http.StatusSeeOther)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/assets/%s", asset.ID()),
		http.StatusSeeOther,
	)
}

func (s *Server) handleEditAssetForm(w http.ResponseWriter, r *http.Request) {
	id := domain.AssetID(chi.URLParam(r, "id"))

	asset, err := s.services.GetAsset.Execute(r.Context(), app.GetAssetCommand{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyAssetID) {
			s.renderNotFound(w, r)
			return
		}

		s.logger.Error("failed to get asset for edit", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := newAssetEditPageData(r, newAssetView(asset), nil)

	if err := s.renderer.Render(w, http.StatusOK, "asset_form.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleUpdateAsset(w http.ResponseWriter, r *http.Request) {
	id := domain.AssetID(chi.URLParam(r, "id"))

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	asset, err := s.services.UpdateAsset.Execute(r.Context(), app.UpdateAssetCommand{
		ID:   id,
		Name: r.FormValue("name"),
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyAssetID) {
			s.renderNotFound(w, r)
			return
		}

		data := newAssetEditPageDataFromRequest(r, id, err)

		if renderErr := s.renderer.Render(w, http.StatusUnprocessableEntity, "asset_form.html", data); renderErr != nil {
			s.handleRenderError(w, renderErr)
		}

		return
	}

	s.recordAuditEventOrLog(
		r,
		domain.AuditEventAssetUpdated,
		domain.AuditEntityAsset,
		string(asset.ID()),
		map[string]string{
			"asset_id": string(asset.ID()),
		},
	)

	http.Redirect(w, r, "/assets/"+string(asset.ID()), http.StatusSeeOther)
}

func (s *Server) handleShowAsset(w http.ResponseWriter, r *http.Request) {
	id := domain.AssetID(chi.URLParam(r, "id"))

	asset, err := s.services.GetAsset.Execute(r.Context(), app.GetAssetCommand{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			s.renderNotFound(w, r)
			return
		}

		if errors.Is(err, domain.ErrEmptyAssetID) {
			s.renderNotFound(w, r)
			return
		}

		s.logger.Error("failed to show asset", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	holders, err := s.services.ListCurrentAssetHolders.Execute(r.Context(), app.ListCurrentAssetHoldersCommand{
		AssetID: asset.ID(),
	})
	if err != nil {
		s.logger.Error("failed to list current asset holders", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	historyItems, err := s.services.ListAssetCustodyHistory.Execute(r.Context(), app.ListAssetCustodyHistoryCommand{
		AssetID: asset.ID(),
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyAssetID) {
			s.renderNotFound(w, r)
			return
		}

		s.logger.Error("failed to list asset custody history", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	view := newAssetView(asset)
	timelineItems := make([]custodyTimelineItemView, 0, len(historyItems))
	for _, item := range historyItems {
		timelineItems = append(timelineItems, newCustodyTimelineItemFromAssetHistoryItem(item, asset.ID()))
	}

	data := assetShowPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             asset.Name(),
		Asset:             view,
		DetailsSection: newSectionHeader(
			"Identificação",
			"Dados do material",
			"Informações principais do material.",
		),
		DetailsFields: []detailFieldView{
			newDetailField("Nome", view.Name),
			newDetailField("Situação", view.StatusLabel),
		},
		CustodyHistorySection: newSectionHeader(
			"Histórico",
			"Histórico de custódia",
			"Movimentações de cautela e descautela envolvendo este material.",
		),
		CustodyTimeline: custodyTimelineView{
			Items: timelineItems,
			EmptyState: newEmptyState(
				"Nenhuma movimentação de custódia",
				"Este material ainda não possui cautelas ou descautelas registradas.",
				newPageAction("Registrar cautela", "/custody/checkouts/new"),
			),
		},
		Holders: make([]assetHolderView, 0, len(holders)),
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		assetsBreadcrumb(),
		currentBreadcrumb(data.Asset.Name),
	}

	for _, holder := range holders {
		statusLabel := activeStatusLabel(holder.PersonnelActive)

		if !holder.PersonnelActive {
			data.HasInactiveHolders = true
		}

		data.Holders = append(data.Holders, assetHolderView{
			PersonnelID:          string(holder.PersonnelID),
			PersonnelFullName:    holder.PersonnelFullName,
			PersonnelDisplay:     militaryDisplayName(holder.PersonnelRank, holder.PersonnelAlias),
			PersonnelActive:      holder.PersonnelActive,
			PersonnelStatusLabel: statusLabel,
			Quantity:             holder.Quantity,
		})
	}

	data.HasHolders = len(data.Holders) > 0
	data.ShowInactiveCustodyWarning = !asset.Active() && data.HasHolders

	if err := s.renderer.Render(w, http.StatusOK, "assets_show.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newCustodyTimelineItemFromAssetHistoryItem(
	item ports.AssetCustodyHistoryItem,
	currentAssetID domain.AssetID,
) custodyTimelineItemView {
	lines := make([]custodyTimelineLineView, 0, len(item.Lines))
	for _, line := range item.Lines {
		lines = append(lines, custodyTimelineLineView{
			AssetID:     string(line.AssetID),
			AssetName:   line.AssetName,
			AssetURL:    "/assets/" + string(line.AssetID),
			Quantity:    formatTimelineQuantity(line.Quantity),
			Highlighted: line.AssetID == currentAssetID,
		})
	}

	typeLabel := custodyTransactionTypeLabel(item.Type)
	transactionURL := "/custody/transactions/" + string(item.ID)

	return custodyTimelineItemView{
		ID:                   string(item.ID),
		URL:                  transactionURL,
		TypeLabel:            typeLabel,
		TypeTone:             custodyTimelineTypeTone(typeLabel),
		DateLabel:            formatTimelineDate(item.CreatedAt),
		TimeLabel:            formatTimelineTime(item.CreatedAt),
		PersonnelLabel:       militaryDisplayName(domain.Rank(item.PersonnelRank), item.PersonnelAlias),
		PersonnelURL:         "/personnel/" + string(item.PersonnelID),
		OperatorLabel:        militaryDisplayName(item.OperatorRank, item.OperatorAlias),
		OperatorURL:          "/admin/operators/" + string(item.OperatorID),
		Edited:               item.EditCount > 0,
		EditCountLabel:       custodyEditCountLabel(item.EditCount),
		Notes:                item.Notes,
		Lines:                lines,
		PrimaryActionLabel:   "Abrir recibo",
		PrimaryActionURL:     transactionURL,
		SecondaryActionLabel: "Editar",
		SecondaryActionURL:   transactionURL + "/edit",
	}
}

func (s *Server) handleDeactivateAsset(w http.ResponseWriter, r *http.Request) {
	assetID := domain.AssetID(chi.URLParam(r, "id"))

	err := s.services.DeactivateAsset.Execute(r.Context(), app.DeactivateAssetCommand{
		ID: assetID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyAssetID) {
			s.renderNotFound(w, r)
			return
		}

		s.logger.Error("failed to deactivate asset", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.recordAuditEventOrLog(
		r,
		domain.AuditEventAssetDeactivated,
		domain.AuditEntityAsset,
		string(assetID),
		nil,
	)

	http.Redirect(w, r, "/assets/"+string(assetID), http.StatusSeeOther)
}

func (s *Server) handleReactivateAsset(w http.ResponseWriter, r *http.Request) {
	assetID := domain.AssetID(chi.URLParam(r, "id"))

	err := s.services.ReactivateAsset.Execute(r.Context(), app.ReactivateAssetCommand{
		ID: assetID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyAssetID) {
			s.renderNotFound(w, r)
			return
		}

		s.logger.Error("failed to reactivate asset", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.recordAuditEventOrLog(
		r,
		domain.AuditEventAssetReactivated,
		domain.AuditEntityAsset,
		string(assetID),
		nil,
	)

	http.Redirect(w, r, "/assets/"+string(assetID), http.StatusSeeOther)
}

func (s *Server) renderNewAssetFormWithError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	err error,
	name string,
) {
	data := newAssetCreatePageData(r, assetFormView{
		Name: name,
	}, assetFeedbackFromError(err))

	if err := s.renderer.Render(w, status, "asset_form.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newAssetCreatePageDataFromMessage(r *http.Request, message string) assetFormPageData {
	data := newAssetCreatePageData(r, assetFormView{
		Name: r.FormValue("name"),
	}, newErrorFeedback(message))
	data.Error = message

	return data
}

func newAssetCreatePageData(
	r *http.Request,
	form assetFormView,
	feedback *feedbackMessageView,
) assetFormPageData {
	errorMessage := ""
	if feedback != nil && feedback.Kind == "error" {
		errorMessage = feedback.Message
	}

	data := assetFormPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Cadastrar material",
		Header: newPageHeader(
			assetPluralLabel(),
			"Cadastrar material",
			"Cadastre um material que poderá ser controlado por custódia.",
			nil,
		),
		FormActions: newFormActionsWithSaveAndCreateAnother(
			"Cadastrar material",
			"Cancelar",
			"/assets",
			"Salvar e cadastrar outro",
		),
		Feedback: feedback,
		DetailsSection: newSectionHeader(
			"Identificação",
			"Dados do material",
			"Informe o nome do material que será controlado por custódia.",
		),
		ActionURL: "/assets",
		Error:     errorMessage,
		Form:      form,
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		assetsBreadcrumb(),
		currentBreadcrumb("Cadastrar material"),
	}

	return data
}

func newAssetEditPageData(
	r *http.Request,
	asset assetView,
	err error,
) assetFormPageData {
	errorMessage := ""
	if err != nil {
		errorMessage = duplicateAssetErrorMessage(err)
	}

	data := assetFormPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Editar material",
		Header: newPageHeader(
			assetPluralLabel(),
			"Editar material",
			"Atualize os dados cadastrais do material.",
			nil,
		),
		FormActions: newFormActions(
			"Salvar alterações",
			"Cancelar",
			"/assets/"+asset.ID,
		),
		Feedback: assetFeedbackFromError(err),
		DetailsSection: newSectionHeader(
			"Identificação",
			"Dados do material",
			"Atualize o nome do material.",
		),
		ActionURL: "/assets/" + asset.ID,
		Error:     errorMessage,
		Form: assetFormView{
			Name: asset.Name,
		},
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		assetsBreadcrumb(),
		{
			Label: asset.Name,
			URL:   "/assets/" + asset.ID,
		},
		currentBreadcrumb("Editar"),
	}

	return data
}

func newAssetEditPageDataFromRequest(
	r *http.Request,
	id domain.AssetID,
	err error,
) assetFormPageData {
	errorMessage := ""
	if err != nil {
		errorMessage = duplicateAssetErrorMessage(err)
	}

	data := assetFormPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Editar material",
		Header: newPageHeader(
			assetPluralLabel(),
			"Editar material",
			"Atualize os dados cadastrais do material.",
			nil,
		),
		FormActions: newFormActions(
			"Salvar alterações",
			"Cancelar",
			"/assets/"+string(id),
		),
		Feedback: assetFeedbackFromError(err),
		DetailsSection: newSectionHeader(
			"Identificação",
			"Dados do material",
			"Atualize o nome do material.",
		),
		ActionURL: "/assets/" + string(id),
		Error:     errorMessage,
		Form: assetFormView{
			Name: r.FormValue("name"),
		},
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		assetsBreadcrumb(),
		currentBreadcrumb("Editar material"),
	}

	return data
}

func humanizeAssetError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyAssetName):
		return "O nome do material é obrigatório."
	case errors.Is(err, domain.ErrDuplicateAssetName):
		return "Este material já está cadastrado."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "O material é obrigatório."
	default:
		return "Não foi possível cadastrar o material."
	}
}
