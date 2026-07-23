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

type assetIndexPageData struct {
	privateLayoutData
	Header       pageHeaderView
	EmptyState   emptyStateView
	Title        string
	Assets       []assetView
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

	assets, err := s.services.ListAssets.Execute(r.Context(), app.ListAssetsCommand{
		Limit:        100,
		StatusFilter: string(statusFilter),
	})
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
		StatusFilter: string(statusFilter),
		StatusTabs:   newStatusFilterTabs("/assets", statusFilter),
	}
	data.Breadcrumbs = []breadcrumbItemView{
		homeBreadcrumb(),
		currentBreadcrumb(assetPluralLabel()),
	}

	for _, item := range assets {
		data.Assets = append(data.Assets, newAssetView(item))
	}

	if err := s.renderer.Render(w, http.StatusOK, "assets_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleNewAssetForm(w http.ResponseWriter, r *http.Request) {
	data := newAssetCreatePageData(r, assetFormView{}, nil)

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

	view := newAssetView(asset)

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
	}, err)

	if err := s.renderer.Render(w, status, "asset_form.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func newAssetCreatePageDataFromMessage(r *http.Request, message string) assetFormPageData {
	data := newAssetCreatePageData(r, assetFormView{
		Name: r.FormValue("name"),
	}, nil)
	data.Error = message
	data.Feedback = newErrorFeedback(message)

	return data
}

func newAssetCreatePageData(
	r *http.Request,
	form assetFormView,
	err error,
) assetFormPageData {
	errorMessage := ""
	if err != nil {
		errorMessage = duplicateAssetErrorMessage(err)
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
		FormActions: newFormActions(
			"Cadastrar material",
			"Cancelar",
			"/assets",
		),
		Feedback: assetFeedbackFromError(err),
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
