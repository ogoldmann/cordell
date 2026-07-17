package web

import "cordell/internal/domain"

type currentOperatorView struct {
	ID                 string
	Username           string
	Role               string
	RoleLabel          string
	CanManageOperators bool
}

func newCurrentOperatorView(operator domain.Operator) currentOperatorView {
	return currentOperatorView{
		ID:                 string(operator.ID()),
		Username:           operator.Username(),
		Role:               operator.Role().String(),
		RoleLabel:          operator.Role().Label(),
		CanManageOperators: operator.Role().CanManageOperators(),
	}
}
