package web

import (
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

type operatorSummaryView struct {
	ID             string
	RegistrationID string
	Alias          string
	Rank           string
	RankLabel      string
	DisplayName    string
	Role           string
	RoleLabel      string
	Active         bool
	CreatedAt      string
}

type operatorDetailView struct {
	ID               string
	RegistrationID   string
	Alias            string
	Rank             string
	RankLabel        string
	DisplayName      string
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
		ID:             string(operator.ID),
		RegistrationID: operator.RegistrationID.String(),
		Alias:          operator.Alias,
		Rank:           operator.Rank.String(),
		RankLabel:      operator.Rank.Label(),
		DisplayName:    operatorDisplayName(operator.Rank, operator.Alias),
		Role:           operator.Role.String(),
		RoleLabel:      operatorRoleLabel(operator.Role),
		Active:         operator.Active,
		CreatedAt:      formatTimestamp(operator.CreatedAt),
	}
}

func newOperatorDetailView(
	operator ports.OperatorSummary,
	currentOperatorID domain.OperatorID,
) operatorDetailView {
	return operatorDetailView{
		ID:               string(operator.ID),
		RegistrationID:   operator.RegistrationID.String(),
		Alias:            operator.Alias,
		Rank:             operator.Rank.String(),
		RankLabel:        operator.Rank.Label(),
		DisplayName:      operatorDisplayName(operator.Rank, operator.Alias),
		Role:             operator.Role.String(),
		RoleLabel:        operatorRoleLabel(operator.Role),
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
			Label:    operatorRoleLabel(role),
			Selected: role.String() == selectedRole,
		})
	}

	return options
}

func formatTimestamp(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Local().Format("02/01/2006 15:04")
}

func operatorDisplayName(rank domain.Rank, alias string) string {
	if alias == "" {
		return rank.Label()
	}

	return rank.Label() + " " + alias
}
