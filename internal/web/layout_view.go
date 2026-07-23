package web

import "cordell/internal/domain"

type currentOperatorView struct {
	ID                 string
	RegistrationID     string
	Alias              string
	RankLabel          string
	DisplayName        string
	Role               string
	RoleLabel          string
	CanManageOperators bool
}

func newCurrentOperatorView(operator domain.Operator) currentOperatorView {
	displayName := militaryDisplayName(operator.Rank(), operator.Alias())

	return currentOperatorView{
		ID:                 string(operator.ID()),
		RegistrationID:     operator.RegistrationID().String(),
		Alias:              operator.Alias(),
		RankLabel:          operator.Rank().DisplayLabel(),
		DisplayName:        displayName,
		Role:               operator.Role().String(),
		RoleLabel:          operatorRoleLabel(operator.Role()),
		CanManageOperators: operator.Role().CanManageOperators(),
	}
}
