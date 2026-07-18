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
	Title        string
	Assets       []assetView
	StatusFilter string
	StatusTabs   []statusFilterTabView
}

type assetNewPageData struct {
	privateLayoutData
	Title string
	Error string
	Name  string
}

type assetShowPageData struct {
	privateLayoutData
	Title                      string
	Asset                      assetView
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
		Title:             "Assets",
		Assets:            make([]assetView, 0, len(assets)),
		StatusFilter:      string(statusFilter),
		StatusTabs:        newStatusFilterTabs("/assets", statusFilter),
	}

	for _, item := range assets {
		data.Assets = append(data.Assets, newAssetView(item))
	}

	if err := s.renderer.Render(w, http.StatusOK, "assets_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleNewAssetForm(w http.ResponseWriter, r *http.Request) {
	data := assetNewPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Create asset",
	}

	if err := s.renderer.Render(w, http.StatusOK, "assets_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateAsset(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderNewAssetFormWithError(
			w,
			r,
			http.StatusBadRequest,
			"Invalid form submission.",
			"",
		)
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
			humanizeAssetError(err),
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

func (s *Server) handleShowAsset(w http.ResponseWriter, r *http.Request) {
	id := domain.AssetID(chi.URLParam(r, "id"))

	asset, err := s.services.GetAsset.Execute(r.Context(), app.GetAssetCommand{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			http.NotFound(w, r)
			return
		}

		if errors.Is(err, domain.ErrEmptyAssetID) {
			http.NotFound(w, r)
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

	data := assetShowPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             asset.Name(),
		Asset:             newAssetView(asset),
		Holders:           make([]assetHolderView, 0, len(holders)),
	}

	for _, holder := range holders {
		statusLabel := "Inactive"
		if holder.PersonnelActive {
			statusLabel = "Active"
		}

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

func (s *Server) handleConfirmDeactivateAsset(w http.ResponseWriter, r *http.Request) {
	assetID := domain.AssetID(chi.URLParam(r, "id"))

	asset, err := s.services.GetAsset.Execute(r.Context(), app.GetAssetCommand{
		ID: assetID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyAssetID) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to load asset deactivation confirmation", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !asset.Active() {
		http.Redirect(w, r, "/assets/"+string(assetID), http.StatusSeeOther)
		return
	}

	data := confirmationPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Deactivate asset",
		Kicker:            "Asset lifecycle",
		Heading:           "Deactivate " + asset.Name() + "?",
		Description:       "This asset will be removed from normal checkout workflows, but historical custody records, receipts, and audit events will remain available.",
		Warning:           "Deactivation does not settle current custody. Existing pending custody must still be returned or corrected later.",
		ConfirmLabel:      "Deactivate asset",
		CancelLabel:       "Cancel",
		ConfirmAction:     "/assets/" + string(assetID) + "/deactivate",
		CancelURL:         "/assets/" + string(assetID),
		ConfirmationStyle: "warning",
	}

	if err := s.renderer.Render(w, http.StatusOK, "confirmation.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleConfirmReactivateAsset(w http.ResponseWriter, r *http.Request) {
	assetID := domain.AssetID(chi.URLParam(r, "id"))

	asset, err := s.services.GetAsset.Execute(r.Context(), app.GetAssetCommand{
		ID: assetID,
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) || errors.Is(err, domain.ErrEmptyAssetID) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to load asset reactivation confirmation", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if asset.Active() {
		http.Redirect(w, r, "/assets/"+string(assetID), http.StatusSeeOther)
		return
	}

	data := confirmationPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Reactivate asset",
		Kicker:            "Asset lifecycle",
		Heading:           "Reactivate " + asset.Name() + "?",
		Description:       "This asset will become available for normal checkout workflows again.",
		Warning:           "",
		ConfirmLabel:      "Reactivate asset",
		CancelLabel:       "Cancel",
		ConfirmAction:     "/assets/" + string(assetID) + "/reactivate",
		CancelURL:         "/assets/" + string(assetID),
		ConfirmationStyle: "primary",
	}

	if err := s.renderer.Render(w, http.StatusOK, "confirmation.html", data); err != nil {
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
			http.NotFound(w, r)
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
			http.NotFound(w, r)
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
	message string,
	name string,
) {
	data := assetNewPageData{
		privateLayoutData: newPrivateLayoutData(r),
		Title:             "Create asset",
		Error:             message,
		Name:              name,
	}

	if err := s.renderer.Render(w, status, "assets_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func humanizeAssetError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyAssetName):
		return "Asset name is required."
	case errors.Is(err, domain.ErrDuplicateAssetName):
		return "Asset name is already registered."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Asset ID is required."
	default:
		return "Could not create asset."
	}
}
