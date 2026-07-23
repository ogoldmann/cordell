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
	AssetID              string
	Quantity             string
	CurrentAssetLabel    string
	CurrentAssetIsActive bool
	NeedsReplacement     bool
}

type custodyAssetOptionView struct {
	ID             string
	Label          string
	HasMaxQuantity bool
	MaxQuantity    int
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
		return "Adicione pelo menos uma linha de material."
	case errors.Is(err, domain.ErrEmptyAssetID):
		return "Cada linha deve ter um material."
	case errors.Is(err, domain.ErrInvalidQuantity):
		return "Cada linha deve ter uma quantidade válida."
	default:
		return "Revise as linhas de materiais e tente novamente."
	}
}

func defaultCustodyLineFormRows() []custodyLineFormRowView {
	return []custodyLineFormRowView{
		{
			Quantity: "1",
		},
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

		if row.Quantity == "" {
			row.Quantity = "1"
		}

		rows = append(rows, row)
	}

	if len(rows) == 0 {
		return defaultCustodyLineFormRows()
	}

	return rows
}

func ensureAtLeastOneCustodyLineFormRow(rows []custodyLineFormRowView) []custodyLineFormRowView {
	if len(rows) == 0 {
		return defaultCustodyLineFormRows()
	}

	for index := range rows {
		if rows[index].Quantity == "" {
			rows[index].Quantity = "1"
		}
	}

	return rows
}

func newCustodyAssetOptions(assets []domain.Asset) []custodyAssetOptionView {
	options := make([]custodyAssetOptionView, 0, len(assets))

	for _, asset := range assets {
		options = append(options, custodyAssetOptionView{
			ID:             string(asset.ID()),
			Label:          asset.Name(),
			HasMaxQuantity: false,
			MaxQuantity:    0,
		})
	}

	return options
}

func newReturnAssetOptions(items []app.CurrentCustodyItem) []custodyAssetOptionView {
	options := make([]custodyAssetOptionView, 0, len(items))

	for _, item := range items {
		label := item.AssetName + " · disponível " + strconv.Itoa(item.Quantity)
		if !item.AssetActive {
			label += " · " + activeStatusLabel(false)
		}

		options = append(options, custodyAssetOptionView{
			ID:             string(item.AssetID),
			Label:          label,
			HasMaxQuantity: true,
			MaxQuantity:    item.Quantity,
		})
	}

	return options
}

func custodyTransactionTotalQuantity(lines []domain.CustodyLine) int {
	total := 0

	for _, line := range lines {
		total += line.Quantity().Int()
	}

	return total
}
