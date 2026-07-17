package postgres

import (
	"context"
	"errors"
	"time"

	"cordell/internal/domain"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// OperatorSessionRepository persists operator sessions in PostgreSQL.
type OperatorSessionRepository struct {
	queries *db.Queries
}

// NewOperatorSessionRepository creates an OperatorSessionRepository.
func NewOperatorSessionRepository(queries *db.Queries) *OperatorSessionRepository {
	return &OperatorSessionRepository{
		queries: queries,
	}
}

// Save persists an operator session.
func (r *OperatorSessionRepository) Save(ctx context.Context, session domain.OperatorSession) error {
	return r.queries.CreateOperatorSession(ctx, db.CreateOperatorSessionParams{
		ID:         string(session.ID()),
		OperatorID: string(session.OperatorID()),
		TokenHash:  session.TokenHash(),
		CsrfToken:  session.CSRFToken(),
		ExpiresAt:  timeToTimestamptz(session.ExpiresAt()),
		CreatedAt:  timeToTimestamptz(session.CreatedAt()),
	})
}

// FindByTokenHash retrieves an operator session by token hash.
func (r *OperatorSessionRepository) FindByTokenHash(ctx context.Context, tokenHash string) (domain.OperatorSession, error) {
	row, err := r.queries.GetOperatorSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.OperatorSession{}, ports.ErrNotFound
		}

		return domain.OperatorSession{}, err
	}

	return domain.NewOperatorSession(
		domain.OperatorSessionID(row.ID),
		domain.OperatorID(row.OperatorID),
		row.TokenHash,
		row.CsrfToken,
		row.ExpiresAt.Time,
		row.CreatedAt.Time,
	)
}

// DeleteByTokenHash deletes an operator session by token hash.
func (r *OperatorSessionRepository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	return r.queries.DeleteOperatorSessionByTokenHash(ctx, tokenHash)
}

// DeleteExpired deletes expired operator sessions.
func (r *OperatorSessionRepository) DeleteExpired(ctx context.Context, now time.Time) error {
	return r.queries.DeleteExpiredOperatorSessions(ctx, timeToTimestamptz(now))
}

func timeToTimestamptz(value time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  value.UTC(),
		Valid: true,
	}
}

var _ ports.OperatorSessionRepository = (*OperatorSessionRepository)(nil)
