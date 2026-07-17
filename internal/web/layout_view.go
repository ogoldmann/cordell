package web

import "cordell/internal/domain"

type currentOperatorView struct {
	ID       string
	Username string
}

func newCurrentOperatorView(operator domain.Operator) currentOperatorView {
	return currentOperatorView{
		ID:       string(operator.ID()),
		Username: operator.Username(),
	}
}
