package app

import (
	"context"
	"strings"
	"unicode/utf8"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

const defaultListOperatorsLimit = 50
const maxListOperatorsLimit = 200

// ListOperatorsCommand contains list operators options.
type ListOperatorsCommand struct {
	Limit int
}

// ListOperatorsService lists operators for administration.
type ListOperatorsService struct {
	operatorRepository ports.OperatorRepository
}

// NewListOperatorsService creates a ListOperatorsService.
func NewListOperatorsService(operatorRepository ports.OperatorRepository) *ListOperatorsService {
	return &ListOperatorsService{
		operatorRepository: operatorRepository,
	}
}

// Execute lists operator summaries.
func (s *ListOperatorsService) Execute(ctx context.Context, cmd ListOperatorsCommand) ([]ports.OperatorSummary, error) {
	limit := normalizeListOperatorsLimit(cmd.Limit)

	return s.operatorRepository.List(ctx, limit)
}

func normalizeListOperatorsLimit(limit int) int {
	if limit <= 0 {
		return defaultListOperatorsLimit
	}

	if limit > maxListOperatorsLimit {
		return maxListOperatorsLimit
	}

	return limit
}

// DeactivateOperatorCommand contains input required to deactivate an operator.
type DeactivateOperatorCommand struct {
	CurrentOperatorID domain.OperatorID
	OperatorID        domain.OperatorID
}

// DeactivateOperatorService deactivates operators from administration workflows.
type DeactivateOperatorService struct {
	operatorRepository ports.OperatorRepository
	sessionRepository  ports.OperatorSessionRepository
}

// NewDeactivateOperatorService creates a DeactivateOperatorService.
func NewDeactivateOperatorService(
	operatorRepository ports.OperatorRepository,
	sessionRepository ports.OperatorSessionRepository,
) *DeactivateOperatorService {
	return &DeactivateOperatorService{
		operatorRepository: operatorRepository,
		sessionRepository:  sessionRepository,
	}
}

// Execute deactivates an operator and invalidates its sessions.
func (s *DeactivateOperatorService) Execute(ctx context.Context, cmd DeactivateOperatorCommand) error {
	if cmd.CurrentOperatorID == cmd.OperatorID {
		return domain.ErrCannotDeactivateCurrentOperator
	}

	operator, err := s.operatorRepository.FindByID(ctx, cmd.OperatorID)
	if err != nil {
		return err
	}

	if !operator.Active() {
		return nil
	}

	if operator.Role() == domain.OperatorRoleAdmin {
		activeAdminCount, err := s.operatorRepository.CountActiveAdmins(ctx)
		if err != nil {
			return err
		}

		if activeAdminCount <= 1 {
			return domain.ErrCannotDeactivateLastAdmin
		}
	}

	deactivated, err := s.operatorRepository.Deactivate(ctx, operator.ID())
	if err != nil {
		return err
	}

	if !deactivated {
		if operator.Role() == domain.OperatorRoleAdmin {
			return domain.ErrCannotDeactivateLastAdmin
		}

		return ports.ErrNotFound
	}

	return s.sessionRepository.DeleteByOperatorID(ctx, operator.ID())
}

// ChangeOperatorRoleCommand contains input required to change an operator role.
type ChangeOperatorRoleCommand struct {
	CurrentOperatorID domain.OperatorID
	OperatorID        domain.OperatorID
	Role              string
}

// ChangeOperatorRoleService changes operator roles from administration workflows.
type ChangeOperatorRoleService struct {
	operatorRepository ports.OperatorRepository
	sessionRepository  ports.OperatorSessionRepository
}

// NewChangeOperatorRoleService creates a ChangeOperatorRoleService.
func NewChangeOperatorRoleService(
	operatorRepository ports.OperatorRepository,
	sessionRepository ports.OperatorSessionRepository,
) *ChangeOperatorRoleService {
	return &ChangeOperatorRoleService{
		operatorRepository: operatorRepository,
		sessionRepository:  sessionRepository,
	}
}

// Execute changes an operator role and invalidates its sessions.
func (s *ChangeOperatorRoleService) Execute(ctx context.Context, cmd ChangeOperatorRoleCommand) error {
	if cmd.CurrentOperatorID == cmd.OperatorID {
		return domain.ErrCannotChangeCurrentOperatorRole
	}

	role, err := domain.NewOperatorRole(cmd.Role)
	if err != nil {
		return err
	}

	operator, err := s.operatorRepository.FindByID(ctx, cmd.OperatorID)
	if err != nil {
		return err
	}

	if operator.Role() == role {
		return nil
	}

	if operator.Active() &&
		operator.Role() == domain.OperatorRoleAdmin &&
		role != domain.OperatorRoleAdmin {
		activeAdminCount, err := s.operatorRepository.CountActiveAdmins(ctx)
		if err != nil {
			return err
		}

		if activeAdminCount <= 1 {
			return domain.ErrCannotDemoteLastAdmin
		}
	}

	changed, err := s.operatorRepository.ChangeRole(ctx, operator.ID(), role)
	if err != nil {
		return err
	}

	if !changed {
		if operator.Active() &&
			operator.Role() == domain.OperatorRoleAdmin &&
			role != domain.OperatorRoleAdmin {
			return domain.ErrCannotDemoteLastAdmin
		}

		return ports.ErrNotFound
	}

	return s.sessionRepository.DeleteByOperatorID(ctx, operator.ID())
}

// ResetOperatorPasswordCommand contains input required to reset an operator password.
type ResetOperatorPasswordCommand struct {
	CurrentOperatorID domain.OperatorID
	OperatorID        domain.OperatorID
	Password          string
}

// ResetOperatorPasswordService resets operator passwords from administration workflows.
type ResetOperatorPasswordService struct {
	operatorRepository ports.OperatorRepository
	sessionRepository  ports.OperatorSessionRepository
	passwordHasher     ports.PasswordHasher
}

// NewResetOperatorPasswordService creates a ResetOperatorPasswordService.
func NewResetOperatorPasswordService(
	operatorRepository ports.OperatorRepository,
	sessionRepository ports.OperatorSessionRepository,
	passwordHasher ports.PasswordHasher,
) *ResetOperatorPasswordService {
	return &ResetOperatorPasswordService{
		operatorRepository: operatorRepository,
		sessionRepository:  sessionRepository,
		passwordHasher:     passwordHasher,
	}
}

// Execute resets an operator password and invalidates its sessions.
func (s *ResetOperatorPasswordService) Execute(ctx context.Context, cmd ResetOperatorPasswordCommand) error {
	if cmd.CurrentOperatorID == cmd.OperatorID {
		return domain.ErrCannotResetCurrentOperatorPassword
	}

	if strings.TrimSpace(cmd.Password) == "" {
		return domain.ErrEmptyOperatorPassword
	}

	if utf8.RuneCountInString(cmd.Password) < minOperatorPasswordLength {
		return domain.ErrWeakOperatorPassword
	}

	operator, err := s.operatorRepository.FindByID(ctx, cmd.OperatorID)
	if err != nil {
		return err
	}

	if !operator.Active() {
		return ports.ErrNotFound
	}

	passwordHash, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return err
	}

	updated, err := s.operatorRepository.UpdatePasswordHash(ctx, operator.ID(), passwordHash)
	if err != nil {
		return err
	}

	if !updated {
		return ports.ErrNotFound
	}

	return s.sessionRepository.DeleteByOperatorID(ctx, operator.ID())
}

// GetOperatorAdminCommand contains input required to retrieve an operator for administration.
type GetOperatorAdminCommand struct {
	OperatorID domain.OperatorID
}

// GetOperatorAdminService retrieves operator administration details.
type GetOperatorAdminService struct {
	operatorRepository ports.OperatorRepository
}

// NewGetOperatorAdminService creates a GetOperatorAdminService.
func NewGetOperatorAdminService(operatorRepository ports.OperatorRepository) *GetOperatorAdminService {
	return &GetOperatorAdminService{
		operatorRepository: operatorRepository,
	}
}

// Execute retrieves an operator administration read model.
func (s *GetOperatorAdminService) Execute(ctx context.Context, cmd GetOperatorAdminCommand) (ports.OperatorSummary, error) {
	return s.operatorRepository.FindSummaryByID(ctx, cmd.OperatorID)
}
