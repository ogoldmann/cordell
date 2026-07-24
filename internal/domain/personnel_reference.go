package domain

// Rank represents a predefined military/personnel rank.
type Rank string

const (
	RankBrigadierGeneral          Rank = "brigadier_general"
	RankColonel                   Rank = "colonel"
	RankLieutenantColonel         Rank = "lieutenant_colonel"
	RankMajor                     Rank = "major"
	RankCaptain                   Rank = "captain"
	RankFirstLieutenant           Rank = "first_lieutenant"
	RankSecondLieutenant          Rank = "second_lieutenant"
	RankAspirant                  Rank = "aspirant"
	RankWarrantOfficer            Rank = "warrant_officer"
	RankFirstSergeant             Rank = "first_sergeant"
	RankSecondSergeant            Rank = "second_sergeant"
	RankThirdSergeant             Rank = "sergeant"
	RankCorporal                  Rank = "corporal"
	RankSoldier                   Rank = "private"
	RankVariableEffectiveCorporal Rank = "variable_effective_corporal"
	RankVariableEffectiveSoldier  Rank = "variable_effective_soldier"
	RankStudent                   Rank = "student"

	// Backward-compatible aliases kept for older tests/data and previous code.
	RankPrivate    Rank = RankSoldier
	RankSergeant   Rank = RankThirdSergeant
	RankLieutenant Rank = RankFirstLieutenant
)

// RankOption describes one selectable rank in the shared rank catalog.
type RankOption struct {
	Value        Rank
	Label        string
	Abbreviation string
	Order        int
	SearchTerms  []string
}

var rankOptions = []RankOption{
	{Value: RankBrigadierGeneral, Label: "General de Brigada", Abbreviation: "Gen Bda", Order: 10, SearchTerms: []string{"general", "gen"}},
	{Value: RankColonel, Label: "Coronel", Abbreviation: "Cel", Order: 20, SearchTerms: []string{"cel"}},
	{Value: RankLieutenantColonel, Label: "Tenente-Coronel", Abbreviation: "Ten Cel", Order: 30, SearchTerms: []string{"ten cel", "tencel", "tenente coronel"}},
	{Value: RankMajor, Label: "Major", Abbreviation: "Maj", Order: 40, SearchTerms: []string{"maj"}},
	{Value: RankCaptain, Label: "Capitão", Abbreviation: "Cap", Order: 50, SearchTerms: []string{"cap"}},
	{Value: RankFirstLieutenant, Label: "Primeiro-Tenente", Abbreviation: "1º Ten", Order: 60, SearchTerms: []string{"1 ten", "1º ten", "primeiro tenente"}},
	{Value: RankSecondLieutenant, Label: "Segundo-Tenente", Abbreviation: "2º Ten", Order: 70, SearchTerms: []string{"2 ten", "2º ten", "segundo tenente"}},
	{Value: RankAspirant, Label: "Aspirante", Abbreviation: "Asp", Order: 80, SearchTerms: []string{"asp"}},
	{Value: RankWarrantOfficer, Label: "Subtenente", Abbreviation: "Sub Ten", Order: 90, SearchTerms: []string{"st", "sub ten", "subten", "subtenente"}},
	{Value: RankFirstSergeant, Label: "Primeiro-Sargento", Abbreviation: "1º Sgt", Order: 100, SearchTerms: []string{"1 sgt", "1º sgt", "primeiro sargento"}},
	{Value: RankSecondSergeant, Label: "Segundo-Sargento", Abbreviation: "2º Sgt", Order: 110, SearchTerms: []string{"2 sgt", "2º sgt", "segundo sargento"}},
	{Value: RankThirdSergeant, Label: "Terceiro-Sargento", Abbreviation: "3º Sgt", Order: 120, SearchTerms: []string{"3 sgt", "3º sgt", "terceiro sargento", "sgt"}},
	{Value: RankCorporal, Label: "Cabo", Abbreviation: "Cb", Order: 130, SearchTerms: []string{"cb"}},
	{Value: RankSoldier, Label: "Soldado", Abbreviation: "Sd", Order: 140, SearchTerms: []string{"sd"}},
	{Value: RankVariableEffectiveCorporal, Label: "Cabo do Efetivo Variável", Abbreviation: "Cb EV", Order: 150, SearchTerms: []string{"cb ev", "cabo ev"}},
	{Value: RankVariableEffectiveSoldier, Label: "Soldado do Efetivo Variável", Abbreviation: "Sd EV", Order: 160, SearchTerms: []string{"sd ev", "soldado ev"}},
	{Value: RankStudent, Label: "Aluno", Abbreviation: "Al", Order: 170, SearchTerms: []string{"al"}},
}

// RankOptions returns the available shared rank catalog options.
func RankOptions() []RankOption {
	options := make([]RankOption, len(rankOptions))
	copy(options, rankOptions)

	return options
}

// IsValidRank reports whether rank is part of the shared rank catalog.
func IsValidRank(rank Rank) bool {
	for _, option := range rankOptions {
		if option.Value == rank {
			return true
		}
	}

	return false
}

// String returns the internal rank value.
func (r Rank) String() string {
	return string(r)
}

// Label returns the Portuguese full label for the rank.
func (r Rank) Label() string {
	for _, option := range rankOptions {
		if option.Value == r {
			return option.Label
		}
	}

	return string(r)
}

// Abbreviation returns the operational abbreviation for the rank.
func (r Rank) Abbreviation() string {
	for _, option := range rankOptions {
		if option.Value == r {
			return option.Abbreviation
		}
	}

	return string(r)
}

// DisplayLabel returns the full rank label with its operational abbreviation.
func (r Rank) DisplayLabel() string {
	label := r.Label()
	abbreviation := r.Abbreviation()

	if abbreviation == "" || abbreviation == label {
		return label
	}

	return label + " (" + abbreviation + ")"
}

// PersonnelRank is kept as a compatibility alias for personnel code.
type PersonnelRank = Rank

const (
	PersonnelRankBrigadierGeneral          PersonnelRank = RankBrigadierGeneral
	PersonnelRankColonel                   PersonnelRank = RankColonel
	PersonnelRankLieutenantColonel         PersonnelRank = RankLieutenantColonel
	PersonnelRankMajor                     PersonnelRank = RankMajor
	PersonnelRankCaptain                   PersonnelRank = RankCaptain
	PersonnelRankFirstLieutenant           PersonnelRank = RankFirstLieutenant
	PersonnelRankSecondLieutenant          PersonnelRank = RankSecondLieutenant
	PersonnelRankAspirant                  PersonnelRank = RankAspirant
	PersonnelRankWarrantOfficer            PersonnelRank = RankWarrantOfficer
	PersonnelRankFirstSergeant             PersonnelRank = RankFirstSergeant
	PersonnelRankSecondSergeant            PersonnelRank = RankSecondSergeant
	PersonnelRankThirdSergeant             PersonnelRank = RankThirdSergeant
	PersonnelRankCorporal                  PersonnelRank = RankCorporal
	PersonnelRankSoldier                   PersonnelRank = RankSoldier
	PersonnelRankVariableEffectiveCorporal PersonnelRank = RankVariableEffectiveCorporal
	PersonnelRankVariableEffectiveSoldier  PersonnelRank = RankVariableEffectiveSoldier
	PersonnelRankStudent                   PersonnelRank = RankStudent

	// Backward-compatible aliases.
	PersonnelRankPrivate    PersonnelRank = RankPrivate
	PersonnelRankSergeant   PersonnelRank = RankSergeant
	PersonnelRankLieutenant PersonnelRank = RankLieutenant
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
	PersonnelSectionCommand        PersonnelSection = "command"
	PersonnelSectionPersonnel      PersonnelSection = "personnel"
	PersonnelSectionIntelligence   PersonnelSection = "intelligence"
	PersonnelSectionOperations     PersonnelSection = "operations"
	PersonnelSectionLogistics      PersonnelSection = "logistics"
	PersonnelSectionCommunications PersonnelSection = "communications"
	PersonnelSectionSupply         PersonnelSection = "supply"
	PersonnelSectionArmory         PersonnelSection = "armory"
	PersonnelSectionMaintenance    PersonnelSection = "maintenance"
	PersonnelSectionHealth         PersonnelSection = "health"

	// Backward-compatible aliases kept for older tests/data and previous code.
	PersonnelSectionAdministration PersonnelSection = PersonnelSectionPersonnel
)

// PersonnelSectionOption represents a selectable personnel section.
type PersonnelSectionOption struct {
	Value        PersonnelSection
	Label        string
	Abbreviation string
	Order        int
	SearchTerms  []string
}

var personnelSectionOptions = []PersonnelSectionOption{
	{Value: PersonnelSectionCommand, Label: "Comando", Abbreviation: "Cmdo", Order: 10, SearchTerms: []string{"cmdo", "comando"}},
	{Value: PersonnelSectionPersonnel, Label: "1ª Seção / Pessoal", Abbreviation: "S1", Order: 20, SearchTerms: []string{"s1", "pessoal", "1 secao", "1 seção"}},
	{Value: PersonnelSectionIntelligence, Label: "2ª Seção / Inteligência", Abbreviation: "S2", Order: 30, SearchTerms: []string{"s2", "inteligencia", "inteligência", "2 secao", "2 seção"}},
	{Value: PersonnelSectionOperations, Label: "3ª Seção / Operações", Abbreviation: "S3", Order: 40, SearchTerms: []string{"s3", "operacoes", "operações", "3 secao", "3 seção"}},
	{Value: PersonnelSectionLogistics, Label: "4ª Seção / Logística", Abbreviation: "S4", Order: 50, SearchTerms: []string{"s4", "logistica", "logística", "4 secao", "4 seção"}},
	{Value: PersonnelSectionCommunications, Label: "Comunicações", Abbreviation: "Com", Order: 60, SearchTerms: []string{"com", "comunicacoes", "comunicações"}},
	{Value: PersonnelSectionSupply, Label: "Almoxarifado", Abbreviation: "Almx", Order: 70, SearchTerms: []string{"almx", "almox", "almoxarifado"}},
	{Value: PersonnelSectionArmory, Label: "Reserva de Armamento", Abbreviation: "Res Armt", Order: 80, SearchTerms: []string{"res armt", "armamento", "reserva armamento"}},
	{Value: PersonnelSectionMaintenance, Label: "Manutenção", Abbreviation: "Mnt", Order: 90, SearchTerms: []string{"mnt", "manutencao", "manutenção"}},
	{Value: PersonnelSectionHealth, Label: "Saúde", Abbreviation: "Sau", Order: 100, SearchTerms: []string{"sau", "saude", "saúde"}},
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

// String returns the internal personnel section value.
func (s PersonnelSection) String() string {
	return string(s)
}

// Label returns the Portuguese full label for the personnel section.
func (s PersonnelSection) Label() string {
	for _, option := range personnelSectionOptions {
		if option.Value == s {
			return option.Label
		}
	}

	return string(s)
}

// Abbreviation returns the operational abbreviation for the personnel section.
func (s PersonnelSection) Abbreviation() string {
	for _, option := range personnelSectionOptions {
		if option.Value == s {
			return option.Abbreviation
		}
	}

	return string(s)
}

// DisplayLabel returns the full personnel section label with its abbreviation.
func (s PersonnelSection) DisplayLabel() string {
	label := s.Label()
	abbreviation := s.Abbreviation()

	if abbreviation == "" || abbreviation == label {
		return label
	}

	return label + " (" + abbreviation + ")"
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

func (u OrganizationUnit) Label() string {
	for _, option := range organizationUnitOptions {
		if option.Value == u {
			return option.Label
		}
	}

	return string(u)
}
