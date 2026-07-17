package web

import (
	"context"

	"cordell/internal/domain"
)

type authContextKey string

const (
	currentOperatorContextKey authContextKey = "current_operator"
	currentSessionContextKey  authContextKey = "current_session"
)

func withCurrentOperator(ctx context.Context, operator domain.Operator) context.Context {
	return context.WithValue(ctx, currentOperatorContextKey, operator)
}

func currentOperatorFromContext(ctx context.Context) (domain.Operator, bool) {
	operator, ok := ctx.Value(currentOperatorContextKey).(domain.Operator)
	return operator, ok
}

func withCurrentSession(ctx context.Context, session domain.OperatorSession) context.Context {
	return context.WithValue(ctx, currentSessionContextKey, session)
}

func currentSessionFromContext(ctx context.Context) (domain.OperatorSession, bool) {
	session, ok := ctx.Value(currentSessionContextKey).(domain.OperatorSession)
	return session, ok
}
