package web

import (
	"errors"
	"net/http"
	"strconv"

	"cordell/internal/app"
	"cordell/internal/domain"
	"cordell/internal/ports"
)

type returnNewPageData struct {
	privateLayoutData
	Title                string
	Error                string
	PersonnelOptions     []personnelOptionView
	SelectedPersonnelID  string
	SelectedPersonnel    personnelReturnView
	HasSelectedPersonnel bool
	CurrentItems         []returnCurrentCustodyItemView
	HasCurrentItems      bool
}

type personnelReturnView struct {
	ID          string
	DisplayName string
	FullName    string
}

type returnCurrentCustodyItemView struct {
	AssetID   string
	AssetName string
	Quantity  int
}

func (s *Server) handleNewReturnForm(w http.ResponseWriter, r *http.Request) {
	selectedPersonnelID := domain.PersonnelID(r.URL.Query().Get("personnel_id"))

	data, err := s.buildReturnFormPageData(r, selectedPersonnelID, "")
	if err != nil {
		s.logger.Error("failed to build return form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, http.StatusOK, "return_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) handleCreateReturn(w http.ResponseWriter, r *http.Request) {
	currentOperator, ok := currentOperatorFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	personnelID := domain.PersonnelID(r.FormValue("personnel_id"))
	assetID := domain.AssetID(r.FormValue("asset_id"))
	quantity, err := parsePositiveInt(r.FormValue("quantity"))
	if err != nil {
		s.renderReturnFormError(
			w,
			r,
			http.StatusBadRequest,
			personnelID,
			"Quantity must be a positive number.",
		)
		return
	}

	notes := r.FormValue("notes")

	transaction, err := s.services.RegisterReturn.Execute(r.Context(), app.RegisterReturnCommand{
		PersonnelID: personnelID,
		OperatorID:  currentOperator.ID(),
		Lines: []app.CustodyLineCommand{
			{
				AssetID:  assetID,
				Quantity: quantity,
			},
		},
		Notes: notes,
	})
	if err != nil {
		s.renderReturnFormError(
			w,
			r,
			http.StatusBadRequest,
			personnelID,
			humanizeReturnWebError(err),
		)
		return
	}

	s.recordAuditEventOrLog(
		r,
		domain.AuditEventCustodyReturnCreated,
		domain.AuditEntityCustodyTransaction,
		string(transaction.ID()),
		map[string]string{
			"personnel_id": string(transaction.PersonnelID()),
		},
	)

	http.Redirect(w, r, "/personnel/"+string(transaction.PersonnelID()), http.StatusSeeOther)
}

func (s *Server) renderReturnFormError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	personnelID domain.PersonnelID,
	message string,
) {
	data, err := s.buildReturnFormPageData(r, personnelID, message)
	if err != nil {
		s.logger.Error("failed to build return form error page", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := s.renderer.Render(w, status, "return_new.html", data); err != nil {
		s.handleRenderError(w, err)
	}
}

func (s *Server) buildReturnFormPageData(
	r *http.Request,
	selectedPersonnelID domain.PersonnelID,
	errorMessage string,
) (returnNewPageData, error) {
	personnelList, err := s.services.ListPersonnel.Execute(r.Context(), app.ListPersonnelCommand{
		Limit: 100,
	})
	if err != nil {
		return returnNewPageData{}, err
	}

	data := returnNewPageData{
		privateLayoutData:   newPrivateLayoutData(r),
		Title:               "Register return",
		Error:               errorMessage,
		PersonnelOptions:    make([]personnelOptionView, 0, len(personnelList)),
		SelectedPersonnelID: string(selectedPersonnelID),
	}

	for _, personnel := range personnelList {
		data.PersonnelOptions = append(data.PersonnelOptions, newPersonnelOptionView(personnel))
	}

	if selectedPersonnelID == "" {
		return data, nil
	}

	personnel, err := s.services.GetPersonnel.Execute(r.Context(), app.GetPersonnelCommand{
		ID: selectedPersonnelID,
	})
	if err != nil {
		return returnNewPageData{}, err
	}

	data.SelectedPersonnel = personnelReturnView{
		ID:          string(personnel.ID()),
		DisplayName: militaryDisplayName(personnel.Rank(), personnel.Alias()),
		FullName:    personnel.FullName(),
	}
	data.HasSelectedPersonnel = true

	currentItems, err := s.services.ListCurrentCustody.Execute(r.Context(), app.ListCurrentCustodyCommand{
		PersonnelID: selectedPersonnelID,
	})
	if err != nil {
		return returnNewPageData{}, err
	}

	data.CurrentItems = make([]returnCurrentCustodyItemView, 0, len(currentItems))

	for _, item := range currentItems {
		data.CurrentItems = append(data.CurrentItems, returnCurrentCustodyItemView{
			AssetID:   string(item.AssetID),
			AssetName: item.AssetName,
			Quantity:  item.Quantity,
		})
	}

	data.HasCurrentItems = len(data.CurrentItems) > 0

	return data, nil
}

func humanizeReturnWebError(err error) string {
	switch {
	case errors.Is(err, domain.ErrEmptyPersonnelID):
		return "Personnel is required."
	case errors.Is(err, domain.ErrEmptyOperatorID):
		return "Authenticated operator is required."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Asset is required."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Quantity must be a positive number."
	case errors.Is(err, domain.ErrInsufficientCustodyBalance):
		return "Return quantity is greater than the available custody balance."
	case errors.Is(err, ports.ErrNotFound):
		return "Personnel or asset not found."
	default:
		return "Could not register return."
	}
}

func parsePositiveInt(value string) (int, error) {
	quantity, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if quantity <= 0 {
		return 0, domain.ErrInvalidQuantity
	}

	return quantity, nil
}
