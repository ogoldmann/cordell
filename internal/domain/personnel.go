package domain

import "strings"

// PersonnelID identifies a personnel record.
type PersonnelID string

// Personnel represents a person who can receive assets under custody.
type Personnel struct {
	id               PersonnelID
	fullName         string
	alias            string
	rank             PersonnelRank
	registrationID   RegistrationID
	section          PersonnelSection
	organizationUnit OrganizationUnit
	active           bool
}

// NewPersonnel creates an active Personnel after validating required fields.
func NewPersonnel(
	id PersonnelID,
	fullName string,
	alias string,
	rank PersonnelRank,
	registrationID RegistrationID,
	section PersonnelSection,
	organizationUnit OrganizationUnit,
) (Personnel, error) {
	return buildPersonnel(
		id,
		fullName,
		alias,
		rank,
		registrationID,
		section,
		organizationUnit,
		true,
	)
}

// ReconstitutePersonnel rebuilds a Personnel from persisted state.
func ReconstitutePersonnel(
	id PersonnelID,
	fullName string,
	alias string,
	rank PersonnelRank,
	registrationID RegistrationID,
	section PersonnelSection,
	organizationUnit OrganizationUnit,
	active bool,
) (Personnel, error) {
	return buildPersonnel(
		id,
		fullName,
		alias,
		rank,
		registrationID,
		section,
		organizationUnit,
		active,
	)
}

func buildPersonnel(
	id PersonnelID,
	fullName string,
	alias string,
	rank PersonnelRank,
	registrationID RegistrationID,
	section PersonnelSection,
	organizationUnit OrganizationUnit,
	active bool,
) (Personnel, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Personnel{}, ErrEmptyPersonnelID
	}

	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return Personnel{}, ErrEmptyPersonnelName
	}

	alias = strings.TrimSpace(alias)
	if alias == "" {
		return Personnel{}, ErrEmptyPersonnelAlias
	}

	if !IsValidPersonnelRank(rank) {
		return Personnel{}, ErrInvalidPersonnelRank
	}

	if registrationID == "" {
		return Personnel{}, ErrEmptyRegistrationID
	}

	if !IsValidPersonnelSection(section) {
		return Personnel{}, ErrInvalidPersonnelSection
	}

	if !IsValidOrganizationUnit(organizationUnit) {
		return Personnel{}, ErrInvalidOrganizationUnit
	}

	return Personnel{
		id:               id,
		fullName:         fullName,
		alias:            alias,
		rank:             rank,
		registrationID:   registrationID,
		section:          section,
		organizationUnit: organizationUnit,
		active:           active,
	}, nil
}

// ID returns the personnel identifier.
func (p Personnel) ID() PersonnelID {
	return p.id
}

// FullName returns the personnel full name.
func (p Personnel) FullName() string {
	return p.fullName
}

// Alias returns the personnel operational alias.
func (p Personnel) Alias() string {
	return p.alias
}

// Rank returns the personnel rank.
func (p Personnel) Rank() PersonnelRank {
	return p.rank
}

// RegistrationID returns the personnel registration identifier.
func (p Personnel) RegistrationID() RegistrationID {
	return p.registrationID
}

// Section returns the personnel section.
func (p Personnel) Section() PersonnelSection {
	return p.section
}

// OrganizationUnit returns the personnel organization unit.
func (p Personnel) OrganizationUnit() OrganizationUnit {
	return p.organizationUnit
}

// Active reports whether the personnel record is active.
func (p Personnel) Active() bool {
	return p.active
}
