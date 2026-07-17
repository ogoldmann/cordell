package web

import (
	"time"

	"cordell/internal/domain"
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

type operatorDetailView struct {
	ID               string
	Username         string
	Role             string
	RoleLabel        string
	Active           bool
	CreatedAt        string
	CanDeactivate    bool
	CanReactivate    bool
	CanChangeRole    bool
	CanResetPassword bool
	RoleOptions      []operatorRoleOptionView
}

type operatorRoleOptionView struct {
	Value    string
	Label    string
	Selected bool
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

func newOperatorDetailView(
	operator ports.OperatorSummary,
	currentOperatorID domain.OperatorID,
) operatorDetailView {
	return operatorDetailView{
		ID:               string(operator.ID),
		Username:         operator.Username,
		Role:             operator.Role.String(),
		RoleLabel:        operator.Role.Label(),
		Active:           operator.Active,
		CreatedAt:        formatTimestamp(operator.CreatedAt),
		CanDeactivate:    operator.Active && operator.ID != currentOperatorID,
		CanReactivate:    !operator.Active,
		CanChangeRole:    operator.ID != currentOperatorID,
		CanResetPassword: operator.Active && operator.ID != currentOperatorID,
		RoleOptions:      newOperatorRoleOptionViews(operator.Role.String()),
	}
}

func newOperatorRoleOptionViews(selectedRole string) []operatorRoleOptionView {
	options := make([]operatorRoleOptionView, 0, len(domain.OperatorRoleOptions()))

	for _, role := range domain.OperatorRoleOptions() {
		options = append(options, operatorRoleOptionView{
			Value:    role.String(),
			Label:    role.Label(),
			Selected: role.String() == selectedRole,
		})
	}

	return options
}

func formatTimestamp(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.UTC().Format("2006-01-02 15:04 UTC")
}
