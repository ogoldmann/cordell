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
	Title       string
	SearchQuery string
	Assets      []assetView
}

type assetNewPageData struct {
	Title string
	Error string
	Name  string
}

type assetShowPageData struct {
	Title   string
	Asset   assetView
	Holders []assetHolderView
}

type assetView struct {
	ID     string
	Name   string
	Active bool
}

type assetHolderView struct {
	PersonnelID       string
	PersonnelFullName string
	Quantity          int
}

func (s *Server) handleListAssets(w http.ResponseWriter, r *http.Request) {
	searchQuery := strings.TrimSpace(r.URL.Query().Get("q"))

	assets, err := s.services.SearchAssets.Execute(r.Context(), app.SearchAssetsCommand{
		Query: searchQuery,
		Limit: 50,
	})
	if err != nil {
		s.logger.Error("failed to list assets", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := assetIndexPageData{
		Title:       "Assets",
		SearchQuery: searchQuery,
		Assets:      make([]assetView, 0, len(assets)),
	}

	for _, item := range assets {
		data.Assets = append(data.Assets, assetView{
			ID:     string(item.ID()),
			Name:   item.Name(),
			Active: item.Active(),
		})
	}

	if err := s.renderer.Render(w, http.StatusOK, "assets_index.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleNewAssetForm(w http.ResponseWriter, r *http.Request) {
	data := assetNewPageData{
		Title: "Create asset",
	}

	if err := s.renderer.Render(w, http.StatusOK, "assets_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateAsset(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderNewAssetFormWithError(
			w,
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
		Title: asset.Name(),
		Asset: assetView{
			ID:     string(asset.ID()),
			Name:   asset.Name(),
			Active: asset.Active(),
		},
		Holders: make([]assetHolderView, 0, len(holders)),
	}

	for _, holder := range holders {
		data.Holders = append(data.Holders, assetHolderView{
			PersonnelID:       string(holder.PersonnelID),
			PersonnelFullName: holder.PersonnelFullName,
			Quantity:          holder.Quantity,
		})
	}

	if err := s.renderer.Render(w, http.StatusOK, "assets_show.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) renderNewAssetFormWithError(
	w http.ResponseWriter,
	status int,
	message string,
	name string,
) {
	data := assetNewPageData{
		Title: "Create asset",
		Error: message,
		Name:  name,
	}

	if err := s.renderer.Render(w, status, "assets_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func humanizeAssetError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyAssetName):
		return "Asset name is required."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Asset ID is required."
	default:
		return "Could not create asset."
	}
}
