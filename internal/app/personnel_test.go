package app

import (
	"context"
	"testing"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestCreatePersonnelServiceExecute(t *testing.T) {
	repository := &fakePersonnelRepository{}
	idGenerator := fixedIDGenerator{id: "personnel-1"}

	service := NewCreatePersonnelService(repository, idGenerator)

	personnel, err := service.Execute(context.Background(), validCreatePersonnelCommand("John Doe", "Doe", "52998224725"))
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if personnel.ID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", personnel.ID())
	}

	if personnel.FullName() != "John Doe" {
		t.Fatalf("expected trimmed full name John Doe, got %s", personnel.FullName())
	}

	if len(repository.saved) != 1 {
		t.Fatalf("expected 1 saved personnel, got %d", len(repository.saved))
	}

	if personnel.Alias() != "Doe" {
		t.Fatalf("expected alias Doe, got %s", personnel.Alias())
	}

	if personnel.RegistrationID() != "52998224725" {
		t.Fatalf("expected registration id 52998224725, got %s", personnel.RegistrationID())
	}
}

func TestCreatePersonnelServiceRejectsInvalidPersonnel(t *testing.T) {
	repository := &fakePersonnelRepository{}
	idGenerator := fixedIDGenerator{id: "personnel-1"}

	service := NewCreatePersonnelService(repository, idGenerator)

	_, err := service.Execute(context.Background(), validCreatePersonnelCommand("   ", "Doe", "52998224725"))
	if err != domain.ErrEmptyPersonnelName {
		t.Fatalf("expected ErrEmptyPersonnelName, got %v", err)
	}

	if len(repository.saved) != 0 {
		t.Fatalf("expected no saved personnel, got %d", len(repository.saved))
	}
}

func TestGetPersonnelServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	service := NewGetPersonnelService(repository)

	found, err := service.Execute(context.Background(), GetPersonnelCommand{
		ID: "personnel-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if found.ID() != "personnel-1" {
		t.Fatalf("expected personnel id personnel-1, got %s", found.ID())
	}

	if found.FullName() != "John Doe" {
		t.Fatalf("expected personnel name John Doe, got %s", found.FullName())
	}
}

func TestGetPersonnelServiceRejectsEmptyID(t *testing.T) {
	repository := &fakePersonnelRepository{}
	service := NewGetPersonnelService(repository)

	_, err := service.Execute(context.Background(), GetPersonnelCommand{
		ID: "",
	})
	if err != domain.ErrEmptyPersonnelID {
		t.Fatalf("expected ErrEmptyPersonnelID, got %v", err)
	}
}

func TestListPersonnelServiceExecute(t *testing.T) {
	firstPersonnel := mustBuildPersonnel(t, "personnel-1")
	secondPersonnel := mustBuildPersonnel(t, "personnel-2")

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			firstPersonnel.ID():  firstPersonnel,
			secondPersonnel.ID(): secondPersonnel,
		},
	}

	service := NewListPersonnelService(repository)

	personnel, err := service.Execute(context.Background(), ListPersonnelCommand{
		Limit:        10,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(personnel) != 2 {
		t.Fatalf("expected 2 personnel records, got %d", len(personnel))
	}
}

func TestListPersonnelServiceAppliesDefaultLimit(t *testing.T) {
	repository := &fakePersonnelRepository{}
	service := NewListPersonnelService(repository)

	_, err := service.Execute(context.Background(), ListPersonnelCommand{
		Limit:        0,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestListPersonnelServiceCapsLimit(t *testing.T) {
	repository := &fakePersonnelRepository{}
	service := NewListPersonnelService(repository)

	_, err := service.Execute(context.Background(), ListPersonnelCommand{
		Limit:        1000,
		StatusFilter: string(ports.RecordStatusFilterActive),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestListPersonnelServiceDefaultsToActiveStatusFilter(t *testing.T) {
	repository := &fakePersonnelRepository{}
	service := NewListPersonnelService(repository)

	_, err := service.Execute(context.Background(), ListPersonnelCommand{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.lastStatusFilter != ports.RecordStatusFilterActive {
		t.Fatalf("expected active status filter, got %s", repository.lastStatusFilter)
	}
}

func TestDeactivatePersonnelServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	service := NewDeactivatePersonnelService(repository)

	err := service.Execute(context.Background(), DeactivatePersonnelCommand{
		ID: personnel.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := repository.FindByID(context.Background(), personnel.ID())
	if err != nil {
		t.Fatalf("expected personnel to exist, got %v", err)
	}

	if updated.Active() {
		t.Fatal("expected personnel to be inactive")
	}
}

func TestDeactivatePersonnelServiceIsNoOpForInactivePersonnel(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	inactivePersonnel, err := domain.ReconstitutePersonnel(
		personnel.ID(),
		personnel.FullName(),
		personnel.Alias(),
		personnel.Rank(),
		personnel.RegistrationID(),
		personnel.Section(),
		personnel.OrganizationUnit(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive personnel, got %v", err)
	}

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			inactivePersonnel.ID(): inactivePersonnel,
		},
	}

	service := NewDeactivatePersonnelService(repository)

	err = service.Execute(context.Background(), DeactivatePersonnelCommand{
		ID: inactivePersonnel.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReactivatePersonnelServiceExecute(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")
	inactivePersonnel, err := domain.ReconstitutePersonnel(
		personnel.ID(),
		personnel.FullName(),
		personnel.Alias(),
		personnel.Rank(),
		personnel.RegistrationID(),
		personnel.Section(),
		personnel.OrganizationUnit(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive personnel, got %v", err)
	}

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			inactivePersonnel.ID(): inactivePersonnel,
		},
	}

	service := NewReactivatePersonnelService(repository)

	err = service.Execute(context.Background(), ReactivatePersonnelCommand{
		ID: inactivePersonnel.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := repository.FindByID(context.Background(), inactivePersonnel.ID())
	if err != nil {
		t.Fatalf("expected personnel to exist, got %v", err)
	}

	if !updated.Active() {
		t.Fatal("expected personnel to be active")
	}
}

func TestReactivatePersonnelServiceIsNoOpForActivePersonnel(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	service := NewReactivatePersonnelService(repository)

	err := service.Execute(context.Background(), ReactivatePersonnelCommand{
		ID: personnel.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCreatePersonnelServiceRejectsInvalidRegistrationID(t *testing.T) {
	repository := &fakePersonnelRepository{}
	service := NewCreatePersonnelService(repository, fixedIDGenerator{id: "personnel-1"})

	_, err := service.Execute(context.Background(), CreatePersonnelCommand{
		FullName:         "John Doe",
		Alias:            "Doe",
		Rank:             domain.PersonnelRankSergeant,
		RegistrationID:   "11111111111",
		Section:          domain.PersonnelSectionOperations,
		OrganizationUnit: domain.OrganizationUnitDefault,
	})
	if err != domain.ErrInvalidRegistrationID {
		t.Fatalf("expected ErrInvalidRegistrationID, got %v", err)
	}
}

func TestSearchPersonnelServiceExecute(t *testing.T) {
	firstPersonnel := mustBuildPersonnel(t, "personnel-1")

	secondRegistrationID, err := domain.NewRegistrationID("11144477735")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	secondPersonnel, err := domain.NewPersonnel(
		"personnel-2",
		"Jane Doe",
		"Jane",
		domain.PersonnelRankCorporal,
		secondRegistrationID,
		domain.PersonnelSectionLogistics,
		domain.OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			firstPersonnel.ID():  firstPersonnel,
			secondPersonnel.ID(): secondPersonnel,
		},
	}

	service := NewSearchPersonnelService(repository)

	personnel, err := service.Execute(context.Background(), SearchPersonnelCommand{
		Query: "jane",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(personnel) != 1 {
		t.Fatalf("expected 1 personnel record, got %d", len(personnel))
	}

	if personnel[0].Alias() != "Jane" {
		t.Fatalf("expected alias Jane, got %s", personnel[0].Alias())
	}
}

func TestSearchPersonnelServiceFallsBackToListWhenQueryIsEmpty(t *testing.T) {
	personnel := mustBuildPersonnel(t, "personnel-1")

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			personnel.ID(): personnel,
		},
	}

	service := NewSearchPersonnelService(repository)

	result, err := service.Execute(context.Background(), SearchPersonnelCommand{
		Query: "   ",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 personnel record, got %d", len(result))
	}
}

func TestSearchPersonnelServiceFindsPersonnelByCombinedTerms(t *testing.T) {
	firstRegistrationID, err := domain.NewRegistrationID("52998224725")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	firstPersonnel, err := domain.NewPersonnel(
		"personnel-1",
		"John Alpha",
		"Doe",
		domain.PersonnelRankSergeant,
		firstRegistrationID,
		domain.PersonnelSectionOperations,
		domain.OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	secondRegistrationID, err := domain.NewRegistrationID("11144477735")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	secondPersonnel, err := domain.NewPersonnel(
		"personnel-2",
		"Jane Doe",
		"Smith",
		domain.PersonnelRankCorporal,
		secondRegistrationID,
		domain.PersonnelSectionLogistics,
		domain.OrganizationUnitDefault,
	)
	if err != nil {
		t.Fatalf("expected valid personnel, got %v", err)
	}

	repository := &fakePersonnelRepository{
		byID: map[domain.PersonnelID]domain.Personnel{
			firstPersonnel.ID():  firstPersonnel,
			secondPersonnel.ID(): secondPersonnel,
		},
	}

	service := NewSearchPersonnelService(repository)

	personnel, err := service.Execute(context.Background(), SearchPersonnelCommand{
		Query: "sergeant doe",
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(personnel) != 1 {
		t.Fatalf("expected 1 personnel record, got %d", len(personnel))
	}

	if personnel[0].ID() != firstPersonnel.ID() {
		t.Fatalf("expected personnel %s, got %s", firstPersonnel.ID(), personnel[0].ID())
	}
}
