package domain

// PersonnelRank represents a predefined personnel rank.
type PersonnelRank string

const (
	// PersonnelRankPrivate represents the private rank.
	PersonnelRankPrivate PersonnelRank = "private"

	// PersonnelRankCorporal represents the corporal rank.
	PersonnelRankCorporal PersonnelRank = "corporal"

	// PersonnelRankSergeant represents the sergeant rank.
	PersonnelRankSergeant PersonnelRank = "sergeant"

	// PersonnelRankLieutenant represents the lieutenant rank.
	PersonnelRankLieutenant PersonnelRank = "lieutenant"

	// PersonnelRankCaptain represents the captain rank.
	PersonnelRankCaptain PersonnelRank = "captain"
)

// PersonnelRankOption represents a selectable personnel rank.
type PersonnelRankOption struct {
	Value PersonnelRank
	Label string
}

var personnelRankOptions = []PersonnelRankOption{
	{Value: PersonnelRankPrivate, Label: "Private"},
	{Value: PersonnelRankCorporal, Label: "Corporal"},
	{Value: PersonnelRankSergeant, Label: "Sergeant"},
	{Value: PersonnelRankLieutenant, Label: "Lieutenant"},
	{Value: PersonnelRankCaptain, Label: "Captain"},
}

// PersonnelRankOptions returns the available personnel rank options.
func PersonnelRankOptions() []PersonnelRankOption {
	options := make([]PersonnelRankOption, len(personnelRankOptions))
	copy(options, personnelRankOptions)

	return options
}

// IsValidPersonnelRank reports whether rank is an allowed personnel rank.
func IsValidPersonnelRank(rank PersonnelRank) bool {
	for _, option := range personnelRankOptions {
		if option.Value == rank {
			return true
		}
	}

	return false
}

// PersonnelSection represents a predefined personnel section.
type PersonnelSection string

const (
	// PersonnelSectionAdministration represents the administration section.
	PersonnelSectionAdministration PersonnelSection = "administration"

	// PersonnelSectionOperations represents the operations section.
	PersonnelSectionOperations PersonnelSection = "operations"

	// PersonnelSectionLogistics represents the logistics section.
	PersonnelSectionLogistics PersonnelSection = "logistics"

	// PersonnelSectionMaintenance represents the maintenance section.
	PersonnelSectionMaintenance PersonnelSection = "maintenance"
)

// PersonnelSectionOption represents a selectable personnel section.
type PersonnelSectionOption struct {
	Value PersonnelSection
	Label string
}

var personnelSectionOptions = []PersonnelSectionOption{
	{Value: PersonnelSectionAdministration, Label: "Administration"},
	{Value: PersonnelSectionOperations, Label: "Operations"},
	{Value: PersonnelSectionLogistics, Label: "Logistics"},
	{Value: PersonnelSectionMaintenance, Label: "Maintenance"},
}

// PersonnelSectionOptions returns the available personnel section options.
func PersonnelSectionOptions() []PersonnelSectionOption {
	options := make([]PersonnelSectionOption, len(personnelSectionOptions))
	copy(options, personnelSectionOptions)

	return options
}

// IsValidPersonnelSection reports whether section is an allowed personnel section.
func IsValidPersonnelSection(section PersonnelSection) bool {
	for _, option := range personnelSectionOptions {
		if option.Value == section {
			return true
		}
	}

	return false
}

// OrganizationUnit represents a predefined organization unit.
type OrganizationUnit string

const (
	// OrganizationUnitDefault represents the default organization unit.
	OrganizationUnitDefault OrganizationUnit = "default_unit"
)

// OrganizationUnitOption represents a selectable organization unit.
type OrganizationUnitOption struct {
	Value OrganizationUnit
	Label string
}

var organizationUnitOptions = []OrganizationUnitOption{
	{Value: OrganizationUnitDefault, Label: "Default Unit"},
}

// OrganizationUnitOptions returns the available organization unit options.
func OrganizationUnitOptions() []OrganizationUnitOption {
	options := make([]OrganizationUnitOption, len(organizationUnitOptions))
	copy(options, organizationUnitOptions)

	return options
}

// IsValidOrganizationUnit reports whether unit is an allowed organization unit.
func IsValidOrganizationUnit(unit OrganizationUnit) bool {
	for _, option := range organizationUnitOptions {
		if option.Value == unit {
			return true
		}
	}

	return false
}
