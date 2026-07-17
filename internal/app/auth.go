package app

import (
	"context"
	"errors"
	"strings"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

const defaultOperatorSessionDuration = 12 * time.Hour

// AuthenticateOperatorCommand contains login credentials.
type AuthenticateOperatorCommand struct {
	RegistrationID string
	Password       string
}

// AuthenticateOperatorService verifies operator credentials.
type AuthenticateOperatorService struct {
	operatorRepository ports.OperatorRepository
	passwordHasher     ports.PasswordHasher
}

// NewAuthenticateOperatorService creates an AuthenticateOperatorService.
func NewAuthenticateOperatorService(
	operatorRepository ports.OperatorRepository,
	passwordHasher ports.PasswordHasher,
) *AuthenticateOperatorService {
	return &AuthenticateOperatorService{
		operatorRepository: operatorRepository,
		passwordHasher:     passwordHasher,
	}
}

// Execute authenticates an operator.
func (s *AuthenticateOperatorService) Execute(ctx context.Context, cmd AuthenticateOperatorCommand) (domain.Operator, error) {
	registrationID, err := domain.NewRegistrationID(cmd.RegistrationID)
	if err != nil {
		return domain.Operator{}, domain.ErrInvalidCredentials
	}

	operator, err := s.operatorRepository.FindByRegistrationID(ctx, registrationID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return domain.Operator{}, domain.ErrInvalidCredentials
		}

		return domain.Operator{}, err
	}

	if !operator.Active() {
		return domain.Operator{}, domain.ErrInvalidCredentials
	}

	matches, err := s.passwordHasher.Verify(cmd.Password, operator.PasswordHash())
	if err != nil {
		return domain.Operator{}, err
	}

	if !matches {
		return domain.Operator{}, domain.ErrInvalidCredentials
	}

	return operator, nil
}

// CreateOperatorSessionCommand contains input required to create a session.
type CreateOperatorSessionCommand struct {
	OperatorID domain.OperatorID
	Now        time.Time
}

// CreateOperatorSessionResult contains a session and its raw token.
type CreateOperatorSessionResult struct {
	Session domain.OperatorSession
	Token   string
}

// CreateOperatorSessionService creates operator sessions.
type CreateOperatorSessionService struct {
	sessionRepository  ports.OperatorSessionRepository
	idGenerator        ports.IDGenerator
	tokenGenerator     ports.SessionTokenGenerator
	csrfTokenGenerator ports.SessionTokenGenerator
	tokenHasher        ports.SessionTokenHasher
}

// NewCreateOperatorSessionService creates a CreateOperatorSessionService.
func NewCreateOperatorSessionService(
	sessionRepository ports.OperatorSessionRepository,
	idGenerator ports.IDGenerator,
	tokenGenerator ports.SessionTokenGenerator,
	csrfTokenGenerator ports.SessionTokenGenerator,
	tokenHasher ports.SessionTokenHasher,
) *CreateOperatorSessionService {
	return &CreateOperatorSessionService{
		sessionRepository:  sessionRepository,
		idGenerator:        idGenerator,
		tokenGenerator:     tokenGenerator,
		csrfTokenGenerator: csrfTokenGenerator,
		tokenHasher:        tokenHasher,
	}
}

// Execute creates an operator session and returns the raw token for the cookie.
func (s *CreateOperatorSessionService) Execute(ctx context.Context, cmd CreateOperatorSessionCommand) (CreateOperatorSessionResult, error) {
	now := cmd.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}

	id, err := s.idGenerator.NewID()
	if err != nil {
		return CreateOperatorSessionResult{}, err
	}

	token, err := s.tokenGenerator.NewToken()
	if err != nil {
		return CreateOperatorSessionResult{}, err
	}

	csrfToken, err := s.csrfTokenGenerator.NewToken()
	if err != nil {
		return CreateOperatorSessionResult{}, err
	}

	session, err := domain.NewOperatorSession(
		domain.OperatorSessionID(id),
		cmd.OperatorID,
		s.tokenHasher.Hash(token),
		csrfToken,
		now.Add(defaultOperatorSessionDuration),
		now,
	)
	if err != nil {
		return CreateOperatorSessionResult{}, err
	}

	if err := s.sessionRepository.Save(ctx, session); err != nil {
		return CreateOperatorSessionResult{}, err
	}

	return CreateOperatorSessionResult{
		Session: session,
		Token:   token,
	}, nil
}

// GetOperatorBySessionTokenCommand contains input required to load an operator by session token.
type GetOperatorBySessionTokenCommand struct {
	Token string
	Now   time.Time
}

// GetOperatorBySessionTokenService loads the current operator from a session token.
type GetOperatorBySessionTokenService struct {
	sessionRepository  ports.OperatorSessionRepository
	operatorRepository ports.OperatorRepository
	tokenHasher        ports.SessionTokenHasher
}

// NewGetOperatorBySessionTokenService creates a GetOperatorBySessionTokenService.
func NewGetOperatorBySessionTokenService(
	sessionRepository ports.OperatorSessionRepository,
	operatorRepository ports.OperatorRepository,
	tokenHasher ports.SessionTokenHasher,
) *GetOperatorBySessionTokenService {
	return &GetOperatorBySessionTokenService{
		sessionRepository:  sessionRepository,
		operatorRepository: operatorRepository,
		tokenHasher:        tokenHasher,
	}
}

// Execute loads an active operator from a raw session token.
func (s *GetOperatorBySessionTokenService) Execute(ctx context.Context, cmd GetOperatorBySessionTokenCommand) (CurrentOperatorSession, error) {

	token := strings.TrimSpace(cmd.Token)
	if token == "" {
		return CurrentOperatorSession{}, ports.ErrNotFound
	}

	now := cmd.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}

	session, err := s.sessionRepository.FindByTokenHash(ctx, s.tokenHasher.Hash(token))
	if err != nil {
		return CurrentOperatorSession{}, err
	}

	if session.Expired(now) {
		_ = s.sessionRepository.DeleteByTokenHash(ctx, session.TokenHash())
		return CurrentOperatorSession{}, domain.ErrExpiredOperatorSession
	}

	operator, err := s.operatorRepository.FindByID(ctx, session.OperatorID())
	if err != nil {
		return CurrentOperatorSession{}, err
	}

	if !operator.Active() {
		return CurrentOperatorSession{}, ports.ErrNotFound
	}

	return CurrentOperatorSession{
		Operator: operator,
		Session:  session,
	}, nil
}

// DeleteOperatorSessionCommand contains input required to delete a session.
type DeleteOperatorSessionCommand struct {
	Token string
}

// DeleteOperatorSessionService deletes operator sessions.
type DeleteOperatorSessionService struct {
	sessionRepository ports.OperatorSessionRepository
	tokenHasher       ports.SessionTokenHasher
}

// NewDeleteOperatorSessionService creates a DeleteOperatorSessionService.
func NewDeleteOperatorSessionService(
	sessionRepository ports.OperatorSessionRepository,
	tokenHasher ports.SessionTokenHasher,
) *DeleteOperatorSessionService {
	return &DeleteOperatorSessionService{
		sessionRepository: sessionRepository,
		tokenHasher:       tokenHasher,
	}
}

// Execute deletes an operator session.
func (s *DeleteOperatorSessionService) Execute(ctx context.Context, cmd DeleteOperatorSessionCommand) error {
	token := strings.TrimSpace(cmd.Token)
	if token == "" {
		return nil
	}

	return s.sessionRepository.DeleteByTokenHash(ctx, s.tokenHasher.Hash(token))
}

// CurrentOperatorSession contains the authenticated operator and session.
type CurrentOperatorSession struct {
	Operator domain.Operator
	Session  domain.OperatorSession
}

// DeleteExpiredOperatorSessionsCommand contains input required to delete expired sessions.
type DeleteExpiredOperatorSessionsCommand struct {
	Now time.Time
}

// DeleteExpiredOperatorSessionsService deletes expired operator sessions.
type DeleteExpiredOperatorSessionsService struct {
	sessionRepository ports.OperatorSessionRepository
}

// NewDeleteExpiredOperatorSessionsService creates a DeleteExpiredOperatorSessionsService.
func NewDeleteExpiredOperatorSessionsService(
	sessionRepository ports.OperatorSessionRepository,
) *DeleteExpiredOperatorSessionsService {
	return &DeleteExpiredOperatorSessionsService{
		sessionRepository: sessionRepository,
	}
}

// Execute deletes expired operator sessions.
func (s *DeleteExpiredOperatorSessionsService) Execute(ctx context.Context, cmd DeleteExpiredOperatorSessionsCommand) error {
	now := cmd.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return s.sessionRepository.DeleteExpired(ctx, now)
}
