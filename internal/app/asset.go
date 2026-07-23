package app

import (
	"context"
	"strings"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

// CreateAssetCommand contains the input data required to create an asset.
type CreateAssetCommand struct {
	Name string
}

// CreateAssetService handles the asset creation use case.
type CreateAssetService struct {
	assetRepository ports.AssetRepository
	idGenerator     ports.IDGenerator
}

// NewCreateAssetService creates a CreateAssetService with its dependencies.
func NewCreateAssetService(
	assetRepository ports.AssetRepository,
	idGenerator ports.IDGenerator,
) *CreateAssetService {
	return &CreateAssetService{
		assetRepository: assetRepository,
		idGenerator:     idGenerator,
	}
}

// Execute creates and persists a new asset record.
func (s *CreateAssetService) Execute(ctx context.Context, cmd CreateAssetCommand) (domain.Asset, error) {
	id, err := s.idGenerator.NewID()
	if err != nil {
		return domain.Asset{}, err
	}

	assetID := domain.AssetID(id)

	asset, err := domain.NewAsset(assetID, cmd.Name)
	if err != nil {
		return domain.Asset{}, err
	}

	if err := s.assetRepository.Save(ctx, asset); err != nil {
		return domain.Asset{}, err
	}

	return asset, nil
}

// UpdateAssetCommand contains the input data required to update an asset.
type UpdateAssetCommand struct {
	ID   domain.AssetID
	Name string
}

// UpdateAssetService handles the asset update use case.
type UpdateAssetService struct {
	assetRepository ports.AssetRepository
}

// NewUpdateAssetService creates an UpdateAssetService.
func NewUpdateAssetService(assetRepository ports.AssetRepository) *UpdateAssetService {
	return &UpdateAssetService{
		assetRepository: assetRepository,
	}
}

// Execute updates editable asset details.
func (s *UpdateAssetService) Execute(ctx context.Context, cmd UpdateAssetCommand) (domain.Asset, error) {
	if cmd.ID == "" {
		return domain.Asset{}, domain.ErrEmptyAssetID
	}

	asset, err := s.assetRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		return domain.Asset{}, err
	}

	normalizedName := domain.NormalizeAssetName(cmd.Name)

	_, duplicateFound, err := s.assetRepository.FindByNameExcludingID(ctx, normalizedName, cmd.ID)
	if err != nil {
		return domain.Asset{}, err
	}
	if duplicateFound {
		return domain.Asset{}, domain.ErrDuplicateAssetName
	}

	if err := asset.UpdateDetails(normalizedName); err != nil {
		return domain.Asset{}, err
	}

	if err := s.assetRepository.Update(ctx, asset); err != nil {
		return domain.Asset{}, err
	}

	return asset, nil
}

const (
	defaultAssetListLimit = 50
	maxAssetListLimit     = 100
)

// GetAssetCommand contains the input data required to retrieve an asset.
type GetAssetCommand struct {
	ID domain.AssetID
}

// GetAssetService handles the asset retrieval use case.
type GetAssetService struct {
	assetRepository ports.AssetRepository
}

// NewGetAssetService creates a GetAssetService with its dependencies.
func NewGetAssetService(assetRepository ports.AssetRepository) *GetAssetService {
	return &GetAssetService{
		assetRepository: assetRepository,
	}
}

// Execute retrieves an asset record by identifier.
func (s *GetAssetService) Execute(ctx context.Context, cmd GetAssetCommand) (domain.Asset, error) {
	if cmd.ID == "" {
		return domain.Asset{}, domain.ErrEmptyAssetID
	}

	return s.assetRepository.FindByID(ctx, cmd.ID)
}

// ListAssetsCommand contains the input data required to list assets.
type ListAssetsCommand struct {
	Limit        int
	StatusFilter string
}

// ListAssetsService handles the asset listing use case.
type ListAssetsService struct {
	assetRepository ports.AssetRepository
}

// NewListAssetsService creates a ListAssetsService with its dependencies.
func NewListAssetsService(assetRepository ports.AssetRepository) *ListAssetsService {
	return &ListAssetsService{
		assetRepository: assetRepository,
	}
}

// Execute retrieves a limited list of asset records.
func (s *ListAssetsService) Execute(ctx context.Context, cmd ListAssetsCommand) ([]domain.Asset, error) {
	limit := normalizeAssetLimit(cmd.Limit)
	statusFilter := ports.NormalizeRecordStatusFilter(cmd.StatusFilter)

	return s.assetRepository.List(ctx, limit, statusFilter)
}

func normalizeAssetLimit(limit int) int {
	if limit <= 0 {
		return defaultAssetListLimit
	}

	if limit > maxAssetListLimit {
		return maxAssetListLimit
	}

	return limit
}

// SearchAssetsCommand contains the input data required to search assets.
type SearchAssetsCommand struct {
	Query        string
	Limit        int
	StatusFilter string
}

// SearchAssetsService handles asset search.
type SearchAssetsService struct {
	assetRepository ports.AssetRepository
}

// NewSearchAssetsService creates a SearchAssetsService with its dependencies.
func NewSearchAssetsService(assetRepository ports.AssetRepository) *SearchAssetsService {
	return &SearchAssetsService{
		assetRepository: assetRepository,
	}
}

// Execute searches asset records by query.
func (s *SearchAssetsService) Execute(ctx context.Context, cmd SearchAssetsCommand) ([]domain.Asset, error) {
	limit := normalizeAssetLimit(cmd.Limit)
	query := strings.TrimSpace(cmd.Query)
	statusFilter := ports.NormalizeRecordStatusFilter(cmd.StatusFilter)

	if query == "" {
		return s.assetRepository.List(ctx, limit, statusFilter)
	}

	return s.assetRepository.Search(ctx, query, limit, statusFilter)
}

// DeactivateAssetCommand contains the input data required to deactivate an asset.
type DeactivateAssetCommand struct {
	ID domain.AssetID
}

// DeactivateAssetService handles asset deactivation.
type DeactivateAssetService struct {
	assetRepository ports.AssetRepository
}

// NewDeactivateAssetService creates a DeactivateAssetService.
func NewDeactivateAssetService(assetRepository ports.AssetRepository) *DeactivateAssetService {
	return &DeactivateAssetService{
		assetRepository: assetRepository,
	}
}

// Execute marks an asset as inactive.
func (s *DeactivateAssetService) Execute(ctx context.Context, cmd DeactivateAssetCommand) error {
	if cmd.ID == "" {
		return domain.ErrEmptyAssetID
	}

	asset, err := s.assetRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if !asset.Active() {
		return nil
	}

	deactivated, err := s.assetRepository.Deactivate(ctx, asset.ID())
	if err != nil {
		return err
	}

	if !deactivated {
		return ports.ErrNotFound
	}

	return nil
}

// ReactivateAssetCommand contains the input data required to reactivate an asset.
type ReactivateAssetCommand struct {
	ID domain.AssetID
}

// ReactivateAssetService handles asset reactivation.
type ReactivateAssetService struct {
	assetRepository ports.AssetRepository
}

// NewReactivateAssetService creates a ReactivateAssetService.
func NewReactivateAssetService(assetRepository ports.AssetRepository) *ReactivateAssetService {
	return &ReactivateAssetService{
		assetRepository: assetRepository,
	}
}

// Execute marks an asset as active.
func (s *ReactivateAssetService) Execute(ctx context.Context, cmd ReactivateAssetCommand) error {
	if cmd.ID == "" {
		return domain.ErrEmptyAssetID
	}

	asset, err := s.assetRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if asset.Active() {
		return nil
	}

	reactivated, err := s.assetRepository.Reactivate(ctx, asset.ID())
	if err != nil {
		return err
	}

	if !reactivated {
		return ports.ErrNotFound
	}

	return nil
}
