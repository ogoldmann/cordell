package app

import (
	"context"

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
