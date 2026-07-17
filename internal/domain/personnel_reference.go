package domain

// Rank represents a predefined military/personnel rank.
type Rank string

const (
	RankPrivate    Rank = "private"
	RankCorporal   Rank = "corporal"
	RankSergeant   Rank = "sergeant"
	RankLieutenant Rank = "lieutenant"
	RankCaptain    Rank = "captain"
)

type RankOption struct {
	Value Rank
	Label string
}

var rankOptions = []RankOption{
	{Value: RankPrivate, Label: "Private"},
	{Value: RankCorporal, Label: "Corporal"},
	{Value: RankSergeant, Label: "Sergeant"},
	{Value: RankLieutenant, Label: "Lieutenant"},
	{Value: RankCaptain, Label: "Captain"},
}

func RankOptions() []RankOption {
	options := make([]RankOption, len(rankOptions))
	copy(options, rankOptions)

	return options
}

func IsValidRank(rank Rank) bool {
	for _, option := range rankOptions {
		if option.Value == rank {
			return true
		}
	}

	return false
}

func (r Rank) String() string {
	return string(r)
}

func (r Rank) Label() string {
	for _, option := range rankOptions {
		if option.Value == r {
			return option.Label
		}
	}

	return string(r)
}

// PersonnelRank is kept as a compatibility alias for personnel code.
type PersonnelRank = Rank

const (
	PersonnelRankPrivate    PersonnelRank = RankPrivate
	PersonnelRankCorporal   PersonnelRank = RankCorporal
	PersonnelRankSergeant   PersonnelRank = RankSergeant
	PersonnelRankLieutenant PersonnelRank = RankLieutenant
	PersonnelRankCaptain    PersonnelRank = RankCaptain
)

type PersonnelRankOption = RankOption

func PersonnelRankOptions() []PersonnelRankOption {
	return RankOptions()
}

func IsValidPersonnelRank(rank PersonnelRank) bool {
	return IsValidRank(rank)
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
