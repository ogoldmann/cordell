package app

import (
	"context"
	"strings"

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

// UpdatePersonnelCommand contains the input data required to update personnel.
type UpdatePersonnelCommand struct {
	ID               domain.PersonnelID
	FullName         string
	Alias            string
	Rank             domain.PersonnelRank
	RegistrationID   domain.RegistrationID
	Section          domain.PersonnelSection
	OrganizationUnit domain.OrganizationUnit
}

// UpdatePersonnelService handles the personnel update use case.
type UpdatePersonnelService struct {
	personnelRepository ports.PersonnelRepository
}

// NewUpdatePersonnelService creates an UpdatePersonnelService.
func NewUpdatePersonnelService(personnelRepository ports.PersonnelRepository) *UpdatePersonnelService {
	return &UpdatePersonnelService{
		personnelRepository: personnelRepository,
	}
}

// Execute updates editable personnel details.
func (s *UpdatePersonnelService) Execute(ctx context.Context, cmd UpdatePersonnelCommand) (domain.Personnel, error) {
	if cmd.ID == "" {
		return domain.Personnel{}, domain.ErrEmptyPersonnelID
	}

	personnel, err := s.personnelRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		return domain.Personnel{}, err
	}

	normalizedRegistrationID := domain.RegistrationID(domain.NormalizeRegistrationID(string(cmd.RegistrationID)))

	_, duplicateFound, err := s.personnelRepository.FindByRegistrationIDExcludingID(
		ctx,
		normalizedRegistrationID,
		cmd.ID,
	)
	if err != nil {
		return domain.Personnel{}, err
	}
	if duplicateFound {
		return domain.Personnel{}, domain.ErrDuplicateRegistrationID
	}

	if err := personnel.UpdateDetails(
		cmd.FullName,
		cmd.Alias,
		cmd.Rank,
		normalizedRegistrationID,
		cmd.Section,
		cmd.OrganizationUnit,
	); err != nil {
		return domain.Personnel{}, err
	}

	if err := s.personnelRepository.Update(ctx, personnel); err != nil {
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
	Limit        int
	StatusFilter string
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
	limit := normalizePersonnelLimit(cmd.Limit)
	statusFilter := ports.NormalizeRecordStatusFilter(cmd.StatusFilter)

	return s.personnelRepository.List(ctx, limit, statusFilter)
}

func normalizePersonnelLimit(limit int) int {
	if limit <= 0 {
		return defaultPersonnelListLimit
	}

	if limit > maxPersonnelListLimit {
		return maxPersonnelListLimit
	}

	return limit
}

// SearchPersonnelCommand contains the input data required to search personnel.
type SearchPersonnelCommand struct {
	Query        string
	Limit        int
	StatusFilter string
}

// SearchPersonnelService handles personnel search.
type SearchPersonnelService struct {
	personnelRepository ports.PersonnelRepository
}

// NewSearchPersonnelService creates a SearchPersonnelService with its dependencies.
func NewSearchPersonnelService(personnelRepository ports.PersonnelRepository) *SearchPersonnelService {
	return &SearchPersonnelService{
		personnelRepository: personnelRepository,
	}
}

// Execute searches personnel records by query.
func (s *SearchPersonnelService) Execute(ctx context.Context, cmd SearchPersonnelCommand) ([]domain.Personnel, error) {
	limit := normalizePersonnelLimit(cmd.Limit)
	query := strings.TrimSpace(cmd.Query)
	statusFilter := ports.NormalizeRecordStatusFilter(cmd.StatusFilter)

	if query == "" {
		return s.personnelRepository.List(ctx, limit, statusFilter)
	}

	return s.personnelRepository.Search(ctx, query, limit, statusFilter)
}

// DeactivatePersonnelCommand contains the input data required to deactivate personnel.
type DeactivatePersonnelCommand struct {
	ID domain.PersonnelID
}

// DeactivatePersonnelService handles personnel deactivation.
type DeactivatePersonnelService struct {
	personnelRepository ports.PersonnelRepository
}

// NewDeactivatePersonnelService creates a DeactivatePersonnelService.
func NewDeactivatePersonnelService(personnelRepository ports.PersonnelRepository) *DeactivatePersonnelService {
	return &DeactivatePersonnelService{
		personnelRepository: personnelRepository,
	}
}

// Execute marks personnel as inactive.
func (s *DeactivatePersonnelService) Execute(ctx context.Context, cmd DeactivatePersonnelCommand) error {
	if cmd.ID == "" {
		return domain.ErrEmptyPersonnelID
	}

	personnel, err := s.personnelRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if !personnel.Active() {
		return nil
	}

	deactivated, err := s.personnelRepository.Deactivate(ctx, personnel.ID())
	if err != nil {
		return err
	}

	if !deactivated {
		return ports.ErrNotFound
	}

	return nil
}

// ReactivatePersonnelCommand contains the input data required to reactivate personnel.
type ReactivatePersonnelCommand struct {
	ID domain.PersonnelID
}

// ReactivatePersonnelService handles personnel reactivation.
type ReactivatePersonnelService struct {
	personnelRepository ports.PersonnelRepository
}

// NewReactivatePersonnelService creates a ReactivatePersonnelService.
func NewReactivatePersonnelService(personnelRepository ports.PersonnelRepository) *ReactivatePersonnelService {
	return &ReactivatePersonnelService{
		personnelRepository: personnelRepository,
	}
}

// Execute marks personnel as active.
func (s *ReactivatePersonnelService) Execute(ctx context.Context, cmd ReactivatePersonnelCommand) error {
	if cmd.ID == "" {
		return domain.ErrEmptyPersonnelID
	}

	personnel, err := s.personnelRepository.FindByID(ctx, cmd.ID)
	if err != nil {
		return err
	}

	if personnel.Active() {
		return nil
	}

	reactivated, err := s.personnelRepository.Reactivate(ctx, personnel.ID())
	if err != nil {
		return err
	}

	if !reactivated {
		return ports.ErrNotFound
	}

	return nil
}
