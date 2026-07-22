package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"cordell/internal/app"
	"cordell/internal/domain"
)

var errNoCustodyLineSubmitted = errors.New("no custody line submitted")

type custodyLineFormRowView struct {
	AssetID  string
	Quantity string
}

type custodyAssetOptionView struct {
	ID    string
	Label string
}

func parseCustodyLineCommandsFromRequest(r *http.Request) ([]app.CustodyLineCommand, error) {
	assetIDs := r.PostForm["asset_id"]
	quantities := r.PostForm["quantity"]

	if len(assetIDs) != len(quantities) {
		return nil, domain.ErrInvalidQuantity
	}

	lines := make([]app.CustodyLineCommand, 0, len(assetIDs))

	for index := range assetIDs {
		assetID := strings.TrimSpace(assetIDs[index])
		quantityText := strings.TrimSpace(quantities[index])

		if assetID == "" && quantityText == "" {
			continue
		}

		if assetID == "" {
			return nil, domain.ErrEmptyAssetID
		}

		quantity, err := strconv.Atoi(quantityText)
		if err != nil || quantity <= 0 {
			return nil, domain.ErrInvalidQuantity
		}

		lines = append(lines, app.CustodyLineCommand{
			AssetID:  domain.AssetID(assetID),
			Quantity: quantity,
		})
	}

	if len(lines) == 0 {
		return nil, errNoCustodyLineSubmitted
	}

	return lines, nil
}

func humanizeCustodyLineFormError(err error) string {
	switch {
	case errors.Is(err, errNoCustodyLineSubmitted):
		return "Add at least one asset line."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Each checkout line must have an asset."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Each checkout line must have a valid positive quantity."
	default:
		return "Review the asset lines and try again."
	}
}

func defaultCustodyLineFormRows() []custodyLineFormRowView {
	return []custodyLineFormRowView{
		{},
	}
}

func custodyLineFormRowsFromRequest(r *http.Request) []custodyLineFormRowView {
	assetIDs := r.PostForm["asset_id"]
	quantities := r.PostForm["quantity"]

	rowCount := len(assetIDs)
	if len(quantities) > rowCount {
		rowCount = len(quantities)
	}

	rows := make([]custodyLineFormRowView, 0, rowCount)

	for index := 0; index < rowCount; index++ {
		row := custodyLineFormRowView{}

		if index < len(assetIDs) {
			row.AssetID = strings.TrimSpace(assetIDs[index])
		}

		if index < len(quantities) {
			row.Quantity = strings.TrimSpace(quantities[index])
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return defaultCustodyLineFormRows()
	}

	return rows
}

func newCustodyAssetOptions(assets []domain.Asset) []custodyAssetOptionView {
	options := make([]custodyAssetOptionView, 0, len(assets))

	for _, asset := range assets {
		options = append(options, custodyAssetOptionView{
			ID:    string(asset.ID()),
			Label: asset.Name(),
		})
	}

	return options
}
