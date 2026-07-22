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
	PersonnelOptions     []returnPersonnelOptionView
	SelectedPersonnelID  string
	SelectedPersonnel    personnelReturnView
	HasSelectedPersonnel bool
	AssetOptions         []custodyAssetOptionView
	LineRows             []custodyLineFormRowView
	Notes                string
	TransactionID        string
	HasCurrentCustody    bool
}

type personnelReturnView struct {
	ID          string
	DisplayName string
	FullName    string
}

type returnPersonnelOptionView struct {
	ID       string
	Label    string
	Selected bool
}

func (s *Server) handleNewReturnForm(w http.ResponseWriter, r *http.Request) {
	data, err := s.buildReturnFormPageData(r, returnNewPageData{
		Title:               "Register return",
		SelectedPersonnelID: r.URL.Query().Get("personnel_id"),
		LineRows:            defaultCustodyLineFormRows(),
	})
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
	transactionID := domain.CustodyTransactionID(r.FormValue("transaction_id"))
	if transactionID == "" {
		s.renderReturnFormError(
			w,
			r,
			http.StatusBadRequest,
			personnelID,
			"Form transaction ID is missing. Please reload the page and try again.",
		)
		return
	}

	lines, err := parseCustodyLineCommandsFromRequest(r)
	if err != nil {
		s.renderReturnFormError(
			w,
			r,
			http.StatusBadRequest,
			personnelID,
			humanizeCustodyLineFormError(err),
		)
		return
	}

	notes := r.FormValue("notes")

	result, err := s.services.RegisterReturn.Execute(r.Context(), app.RegisterReturnCommand{
		TransactionID: transactionID,
		PersonnelID:   personnelID,
		OperatorID:    currentOperator.ID(),
		Lines:         lines,
		Notes:         notes,
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

	if result.Created {
		s.recordAuditEventOrLog(
			r,
			domain.AuditEventCustodyReturnCreated,
			domain.AuditEntityCustodyTransaction,
			string(result.Transaction.ID()),
			map[string]string{
				"transaction_id": string(result.Transaction.ID()),
				"personnel_id":   string(result.Transaction.PersonnelID()),
				"line_count":     strconv.Itoa(len(result.Transaction.Lines())),
				"total_quantity": strconv.Itoa(custodyTransactionTotalQuantity(result.Transaction.Lines())),
			},
		)
	}

	http.Redirect(w, r, "/custody/transactions/"+string(result.Transaction.ID()), http.StatusSeeOther)
}

func (s *Server) renderReturnFormError(
	w http.ResponseWriter,
	r *http.Request,
	status int,
	personnelID domain.PersonnelID,
	message string,
) {
	data, err := s.buildReturnFormPageData(r, returnNewPageData{
		Title:               "Register return",
		Error:               message,
		SelectedPersonnelID: string(personnelID),
		LineRows:            custodyLineFormRowsFromRequest(r),
		Notes:               r.FormValue("notes"),
		TransactionID:       r.FormValue("transaction_id"),
	})
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
	data returnNewPageData,
) (returnNewPageData, error) {
	personnelList, err := s.services.ListPersonnelWithCurrentCustody.Execute(r.Context())
	if err != nil {
		return returnNewPageData{}, err
	}

	data.privateLayoutData = newPrivateLayoutData(r)
	data.PersonnelOptions = newReturnPersonnelOptions(personnelList, data.SelectedPersonnelID)
	data.LineRows = ensureAtLeastOneCustodyLineFormRow(data.LineRows)

	if data.TransactionID == "" {
		transactionID, err := newFormTransactionID()
		if err != nil {
			return returnNewPageData{}, err
		}

		data.TransactionID = string(transactionID)
	}

	if data.SelectedPersonnelID == "" {
		return data, nil
	}

	selectedPersonnelID := domain.PersonnelID(data.SelectedPersonnelID)

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

	data.AssetOptions = newReturnAssetOptions(currentItems)
	data.HasCurrentCustody = len(data.AssetOptions) > 0

	return data, nil
}

func newReturnPersonnelOptions(
	personnel []app.PersonnelWithCurrentCustody,
	selectedID string,
) []returnPersonnelOptionView {
	options := make([]returnPersonnelOptionView, 0, len(personnel))

	for _, item := range personnel {
		label := militaryDisplayName(item.Rank, item.Alias) + " - " + item.FullName
		label += " - " + strconv.Itoa(item.TotalQuantity) + " item(s)"

		if !item.Active {
			label += " - Inactive"
		}

		options = append(options, returnPersonnelOptionView{
			ID:       string(item.ID),
			Label:    label,
			Selected: string(item.ID) == selectedID,
		})
	}

	return options
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
