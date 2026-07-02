package app

import (
	"context"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

// CreatePersonnelCommand contains the input data required to create personnel.
type CreatePersonnelCommand struct {
	FullName string
}

// CreatePersonnelService handles the personnel creation use case.
type CreatePersonnelService struct {
	personnelRepository ports.PersonnelRepository
	idGenerator         ports.IDGenerator
}

// NewCreatePersonnelService creates a CreatePersonnelService with its dependencies.
func NewCreatePersonnelService(
	personnelRepository ports.PersonnelRepository,
	idGenerator ports.IDGenerator,
) *CreatePersonnelService {
	return &CreatePersonnelService{
		personnelRepository: personnelRepository,
		idGenerator:         idGenerator,
	}
}

// Execute creates and persists a new personnel record.
func (s *CreatePersonnelService) Execute(ctx context.Context, cmd CreatePersonnelCommand) (domain.Personnel, error) {
	personnelID := domain.PersonnelID(s.idGenerator.NewID())

	personnel, err := domain.NewPersonnel(personnelID, cmd.FullName)
	if err != nil {
		return domain.Personnel{}, err
	}

	if err := s.personnelRepository.Save(ctx, personnel); err != nil {
		return domain.Personnel{}, err
	}

	return personnel, nil
}
