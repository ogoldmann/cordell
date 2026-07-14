package web

import (
	"context"

	"cordell/internal/domain"
)

type authContextKey string

const currentOperatorContextKey authContextKey = "current_operator"

func withCurrentOperator(ctx context.Context, operator domain.Operator) context.Context {
	return context.WithValue(ctx, currentOperatorContextKey, operator)
}

func currentOperatorFromContext(ctx context.Context) (domain.Operator, bool) {
	operator, ok := ctx.Value(currentOperatorContextKey).(domain.Operator)
	return operator, ok
}
