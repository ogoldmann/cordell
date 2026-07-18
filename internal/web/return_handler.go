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

type returnNewPageData struct {
	privateLayoutData
	Title               string
	Error               string
	Personnel           []personnelView
	CurrentCustody      []currentCustodyView
	SelectedPersonnelID string
	SelectedAssetID     string
	Quantity            string
	Notes               string
}

func (s *Server) handleNewReturnForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.buildReturnNewPageData(r, returnNewPageData{
		Title:               "Register return",
		SelectedPersonnelID: r.URL.Query().Get("personnel_id"),
		Quantity:            "1",
	})
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			http.NotFound(w, r)
			return
		}

		s.logger.Error("failed to build return form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, http.StatusOK, "return_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateReturn(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		s.renderReturnFormWithError(
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
		s.renderReturnFormWithError(
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

	_, err = s.services.RegisterReturn.Execute(r.Context(), app.RegisterReturnCommand{
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
		s.renderReturnFormWithError(
			w,
			r,
			http.StatusBadRequest,
			humanizeReturnError(err),
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

func (s *Server) renderReturnFormWithError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	message string,
) {
	data, err := s.buildReturnNewPageData(r, returnNewPageData{
		Title:               "Register return",
		Error:               message,
		SelectedPersonnelID: r.FormValue("personnel_id"),
		SelectedAssetID:     r.FormValue("asset_id"),
		Quantity:            r.FormValue("quantity"),
		Notes:               r.FormValue("notes"),
	})
	if err != nil {
		s.logger.Error("failed to rebuild return form data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, status, "return_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) buildReturnNewPageData(
	r *http.Request,
	data returnNewPageData,
) (returnNewPageData, error) {
	data.privateLayoutData = newPrivateLayoutData(r)

	personnel, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit: 100,
	})
	if err != nil {
		return returnNewPageData{}, err
	}

	data.Personnel = make([]personnelView, 0, len(personnel))
	for _, item := range personnel {
		data.Personnel = append(data.Personnel, newPersonnelView(item))
	}

	if data.SelectedPersonnelID != "" {
		currentCustody, err := s.services.ListCurrentCustody.Execute(r.Context(), app.ListCurrentCustodyCommand{
			PersonnelID: domain.PersonnelID(data.SelectedPersonnelID),
		})
		if err != nil {
			return returnNewPageData{}, err
		}

		data.CurrentCustody = make([]currentCustodyView, 0, len(currentCustody))
		for _, item := range currentCustody {
			data.CurrentCustody = append(data.CurrentCustody, currentCustodyView{
				AssetID:   string(item.AssetID),
				AssetName: item.AssetName,
				Quantity:  item.Quantity,
			})
		}
	}

	if data.Quantity == "" {
		data.Quantity = "1"
	}

	return data, nil
}

func humanizeReturnError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "Personnel is required."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Asset is required."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Quantity must be greater than zero."
	case errors.Is(err, domain.ErrInsufficientCustodyBalance):
		return "Return quantity cannot exceed the current custody quantity."
	case errors.Is(err, ports.ErrNotFound):
		return "Selected personnel or asset could not be found."
	default:
		return "Could not register return."
	}
}
