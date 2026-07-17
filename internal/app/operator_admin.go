package app

import (
	"context"

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
