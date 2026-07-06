package app

import (
	"context"

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
