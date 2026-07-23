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
	displayName := operator.Rank().Label() + " " + operator.Alias()

	return currentOperatorView{
		ID:                 string(operator.ID()),
		RegistrationID:     operator.RegistrationID().String(),
		Alias:              operator.Alias(),
		RankLabel:          operator.Rank().Label(),
		DisplayName:        displayName,
		Role:               operator.Role().String(),
		RoleLabel:          operatorRoleLabel(operator.Role()),
		CanManageOperators: operator.Role().CanManageOperators(),
	}
}
