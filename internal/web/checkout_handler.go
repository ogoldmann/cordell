package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"
)

type checkoutNewPageData struct {
	privateLayoutData
	Title               string
	Error               string
	Personnel           []personnelView
	Assets              []assetView
	SelectedPersonnelID string
	SelectedAssetID     string
	Quantity            string
	Notes               string
}

func (s *Server) handleNewCheckoutForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.buildCheckoutNewPageData(r, checkoutNewPageData{
		Title:               "Register checkout",
		SelectedPersonnelID: r.URL.Query().Get("personnel_id"),
		SelectedAssetID:     r.URL.Query().Get("asset_id"),
		Quantity:            "1",
	})
	if err != nil {
		s.logger.Error("failed to build checkout form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, http.StatusOK, "checkout_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateCheckout(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderCheckoutFormWithError(
			w,
			r,
			http.StatusBadRequest,
			"Invalid form submission.",
		)
		return
	}

	personnelID := r.FormValue("personnel_id")
	assetID := r.FormValue("asset_id")
	quantityText := r.FormValue("quantity")
	notes := r.FormValue("notes")

	quantity, err := strconv.Atoi(quantityText)
	if err != nil {
		s.renderCheckoutFormWithError(
			w,
			r,
			http.StatusBadRequest,
			"Quantity must be a valid number.",
		)
		return
	}

	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err = s.services.RegisterCheckout.Execute(r.Context(), app.RegisterCheckoutCommand{
		PersonnelID: domain.PersonnelID(personnelID),
		OperatorID:  currentOperator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  domain.AssetID(assetID),
				Quantity: quantity,
			},
		},
		Notes: notes,
	})
	if err != nil {
		s.renderCheckoutFormWithError(
			w,
			r,
			http.StatusBadRequest,
			humanizeCheckoutError(err),
		)
		return
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf("/personnel/%s", personnelID),
		http.StatusSeeOther,
	)
}

func (s *Server) renderCheckoutFormWithError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message string,
) {
	data, err := s.buildCheckoutNewPageData(r, checkoutNewPageData{
		Title:               "Register checkout",
		Error:               message,
		SelectedPersonnelID: r.FormValue("personnel_id"),
		SelectedAssetID:     r.FormValue("asset_id"),
		Quantity:            r.FormValue("quantity"),
		Notes:               r.FormValue("notes"),
	})
	if err != nil {
		s.logger.Error("failed to rebuild checkout form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, status, "checkout_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) buildCheckoutNewPageData(
	r *http.Request,
	data checkoutNewPageData,
) (checkoutNewPageData, error) {
	data.privateLayoutData = newPrivateLayoutData(r)

	personnel, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit: 100,
	})
	if err != nil {
		return checkoutNewPageData{}, err
	}

	assets, err := s.services.ListAssets.Execute(r.Context(), app.ListAssetsCommand{
		Limit: 100,
	})
	if err != nil {
		return checkoutNewPageData{}, err
	}

	data.Personnel = make([]personnelView, 0, len(personnel))
	for _, item := range personnel {
		data.Personnel = append(data.Personnel, newPersonnelView(item))
	}

	data.Assets = make([]assetView, 0, len(assets))
	for _, item := range assets {
		data.Assets = append(data.Assets, assetView{
			ID:     string(item.ID()),
			Name:   item.Name(),
			Active: item.Active(),
		})
	}

	if data.Quantity == "" {
		data.Quantity = "1"
	}

	return data, nil
}

func humanizeCheckoutError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "Personnel is required."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Asset is required."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Quantity must be greater than zero."
	case errors.Is(err, ports.ErrNotFound):
		return "Selected personnel or asset could not be found."
	default:
		return "Could not register checkout."
	}
}
