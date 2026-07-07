package app

import (
	"context"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

// CreatePersonnelCommand contains the input data required to create personnel.
type CreatePersonnelCommand struct {
	FullName         string
	Alias            string
	Rank             domain.PersonnelRank
	RegistrationID   string
	Section          domain.PersonnelSection
	OrganizationUnit domain.OrganizationUnit
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
	id, err := s.idGenerator.NewID()
	if err != nil {
		return domain.Personnel{}, err
	}

	registrationID, err := domain.NewRegistrationID(cmd.RegistrationID)
	if err != nil {
		return domain.Personnel{}, err
	}

	personnel, err := domain.NewPersonnel(
		domain.PersonnelID(id),
		cmd.FullName,
		cmd.Alias,
		cmd.Rank,
		registrationID,
		cmd.Section,
		cmd.OrganizationUnit,
	)
	if err != nil {
		return domain.Personnel{}, err
	}

	if err := s.personnelRepository.Save(ctx, personnel); err != nil {
		return domain.Personnel{}, err
	}

	return personnel, nil
}

// GetPersonnelCommand contains the input data required to retrieve personnel.
type GetPersonnelCommand struct {
	ID domain.PersonnelID
}

// GetPersonnelService handles the personnel retrieval use case.
type GetPersonnelService struct {
	personnelRepository ports.PersonnelRepository
}

// NewGetPersonnelService creates a GetPersonnelService with its dependencies.
func NewGetPersonnelService(personnelRepository ports.PersonnelRepository) *GetPersonnelService {
	return &GetPersonnelService{
		personnelRepository: personnelRepository,
	}
}

// Execute retrieves a personnel record by identifier.
func (s *GetPersonnelService) Execute(ctx context.Context, cmd GetPersonnelCommand) (domain.Personnel, error) {
	if cmd.ID == "" {
		return domain.Personnel{}, domain.ErrEmptyPersonnelID
	}

	return s.personnelRepository.FindByID(ctx, cmd.ID)
}

const (
	defaultPersonnelListLimit = 50
	maxPersonnelListLimit     = 100
)

// ListPersonnelCommand contains the input data required to list personnel.
type ListPersonnelCommand struct {
	Limit int
}

// ListPersonnelService handles the personnel listing use case.
type ListPersonnelService struct {
	personnelRepository ports.PersonnelRepository
}

// NewListPersonnelService creates a ListPersonnelService with its dependencies.
func NewListPersonnelService(personnelRepository ports.PersonnelRepository) *ListPersonnelService {
	return &ListPersonnelService{
		personnelRepository: personnelRepository,
	}
}

// Execute retrieves a limited list of personnel records.
func (s *ListPersonnelService) Execute(ctx context.Context, cmd ListPersonnelCommand) ([]domain.Personnel, error) {
	limit := cmd.Limit

	if limit <= 0 {
		limit = defaultPersonnelListLimit
	}

	if limit > maxPersonnelListLimit {
		limit = maxPersonnelListLimit
	}

	return s.personnelRepository.List(ctx, limit)
}
