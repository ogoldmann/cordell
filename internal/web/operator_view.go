package web

import (
	"time"

	"cordell/internal/ports"
)

type operatorSummaryView struct {
	ID        string
	Username  string
	Role      string
	RoleLabel string
	Active    bool
	CreatedAt string
}

func newOperatorSummaryView(operator ports.OperatorSummary) operatorSummaryView {
	return operatorSummaryView{
		ID:        string(operator.ID),
		Username:  operator.Username,
		Role:      operator.Role.String(),
		RoleLabel: operator.Role.Label(),
		Active:    operator.Active,
		CreatedAt: formatTimestamp(operator.CreatedAt),
	}
}

func formatTimestamp(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.UTC().Format("2006-01-02 15:04 UTC")
}
