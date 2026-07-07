package postgres_test

import (
	"cordell/internal/app"
	"cordell/internal/domain"
)

type fixedIDGenerator struct {
	id string
}

func (g fixedIDGenerator) NewID() (string, error) {
	return g.id, nil
}

func validCreatePersonnelCommand(fullName string, alias string, registrationID string) app.CreatePersonnelCommand {
	return app.CreatePersonnelCommand{
		FullName:         fullName,
		Alias:            alias,
		Rank:             domain.PersonnelRankSergeant,
		RegistrationID:   registrationID,
		Section:          domain.PersonnelSectionOperations,
		OrganizationUnit: domain.OrganizationUnitDefault,
	}
}
