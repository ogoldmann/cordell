package app

import (
	"context"
	"strings"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

const defaultGlobalSearchLimitPerGroup = 8
const maxGlobalSearchLimitPerGroup = 20

// GlobalSearchCommand contains the input data required to run a global search.
type GlobalSearchCommand struct {
	Query        string
	LimitPerType int
}

// GlobalSearchResult contains grouped global search results.
type GlobalSearchResult struct {
	Query     string
	Personnel []domain.Personnel
	Assets    []domain.Asset
}

// GlobalSearchService handles global search across searchable records.
type GlobalSearchService struct {
	personnelRepository ports.PersonnelRepository
	assetRepository     ports.AssetRepository
}

// NewGlobalSearchService creates a GlobalSearchService with its dependencies.
func NewGlobalSearchService(
	personnelRepository ports.PersonnelRepository,
	assetRepository ports.AssetRepository,
) *GlobalSearchService {
	return &GlobalSearchService{
		personnelRepository: personnelRepository,
		assetRepository:     assetRepository,
	}
}

// Execute searches across personnel and assets.
func (s *GlobalSearchService) Execute(ctx context.Context, cmd GlobalSearchCommand) (GlobalSearchResult, error) {
	query := strings.TrimSpace(cmd.Query)
	limit := normalizeGlobalSearchLimit(cmd.LimitPerType)

	result := GlobalSearchResult{
		Query: query,
	}

	if query == "" {
		return result, nil
	}

	personnel, err := s.personnelRepository.Search(ctx, query, limit)
	if err != nil {
		return GlobalSearchResult{}, err
	}

	assets, err := s.assetRepository.Search(ctx, query, limit)
	if err != nil {
		return GlobalSearchResult{}, err
	}

	result.Personnel = personnel
	result.Assets = assets

	return result, nil
}

func normalizeGlobalSearchLimit(limit int) int {
	if limit <= 0 {
		return defaultGlobalSearchLimitPerGroup
	}

	if limit > maxGlobalSearchLimitPerGroup {
		return maxGlobalSearchLimitPerGroup
	}

	return limit
}
