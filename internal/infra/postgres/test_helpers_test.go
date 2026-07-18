package postgres_test

import (
	"testing"

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

func mustNewTestOperator(
	t *testing.T,
	id domain.OperatorID,
	registrationID string,
	alias string,
	rank domain.Rank,
	role domain.OperatorRole,
) domain.Operator {
	t.Helper()

	validRegistrationID, err := domain.NewRegistrationID(registrationID)
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	operator, err := domain.NewOperator(
		id,
		validRegistrationID,
		alias,
		rank,
		role,
		"$argon2id$hash",
	)
	if err != nil {
		t.Fatalf("expected valid operator, got %v", err)
	}

	return operator
}
