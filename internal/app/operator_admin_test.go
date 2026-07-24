package app

import (
	"context"
	"testing"
	"time"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

func TestListOperatorsServiceExecute(t *testing.T) {
	createdAt := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	repository := &fakeOperatorRepository{
		summaries: []ports.OperatorSummary{
			{
				ID:             "operator-1",
				RegistrationID: mustRegistrationID(t, "52998224725"),
				Alias:          "silva",
				Rank:           domain.RankSergeant,
				Role:           domain.OperatorRoleAdmin,
				Active:         true,
				CreatedAt:      createdAt,
			},
			{
				ID:             "operator-2",
				RegistrationID: mustRegistrationID(t, "93541134780"),
				Alias:          "costa",
				Rank:           domain.RankCorporal,
				Role:           domain.OperatorRoleOperator,
				Active:         true,
				CreatedAt:      createdAt,
			},
		},
	}

	service := NewListOperatorsService(repository)

	operators, err := service.Execute(context.Background(), ListOperatorsCommand{
		Query:  "silva",
		Status: ports.RecordStatusFilterAll,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(operators) != 1 {
		t.Fatalf("expected 1 operator, got %d", len(operators))
	}

	if operators[0].RegistrationID.String() != "52998224725" {
		t.Fatalf("expected registration id 52998224725, got %s", operators[0].RegistrationID)
	}
}

func TestListOperatorsServiceDefaultsStatusToAll(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewListOperatorsService(repository)

	_, err := service.Execute(context.Background(), ListOperatorsCommand{
		Limit: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.lastFilters.Status != ports.RecordStatusFilterAll {
		t.Fatalf("expected all status filter, got %s", repository.lastFilters.Status)
	}
}

func TestListOperatorsServiceLimitsMaximum(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewListOperatorsService(repository)

	_, err := service.Execute(context.Background(), ListOperatorsCommand{
		Limit: maxListOperatorsLimit + 1,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repository.lastFilters.Limit != maxListOperatorsLimit {
		t.Fatalf("expected limit %d, got %d", maxListOperatorsLimit, repository.lastFilters.Limit)
	}
}

func TestDeactivateOperatorServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	clerk, err := buildOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid clerk, got %v", err)
	}

	session, err := domain.NewOperatorSession(
		"session-1",
		clerk.ID(),
		"hash:token",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin, clerk)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{
			session.TokenHash(): session,
		},
	}

	service := NewDeactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), DeactivateOperatorCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        clerk.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	deactivatedClerk := operatorRepository.byID[clerk.ID()]
	if deactivatedClerk.Active() {
		t.Fatal("expected clerk to be inactive")
	}

	if _, ok := sessionRepository.byTokenHash[session.TokenHash()]; ok {
		t.Fatal("expected clerk sessions to be deleted")
	}
}

func TestDeactivateOperatorServiceRejectsCurrentOperator(t *testing.T) {
	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewDeactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), DeactivateOperatorCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        admin.ID(),
	})
	if err != domain.ErrCannotDeactivateCurrentOperator {
		t.Fatalf("expected ErrCannotDeactivateCurrentOperator, got %v", err)
	}
}

func TestDeactivateOperatorServiceRejectsLastAdmin(t *testing.T) {
	currentAdmin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid current admin, got %v", err)
	}

	targetAdmin, err := buildOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid target admin, got %v", err)
	}

	inactiveCurrentAdmin, err := domain.ReconstituteOperator(
		currentAdmin.ID(),
		currentAdmin.RegistrationID(),
		currentAdmin.Alias(),
		currentAdmin.Rank(),
		currentAdmin.Role(),
		currentAdmin.PasswordHash(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive current admin, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(inactiveCurrentAdmin, targetAdmin)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewDeactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), DeactivateOperatorCommand{
		CurrentOperatorID: inactiveCurrentAdmin.ID(),
		OperatorID:        targetAdmin.ID(),
	})
	if err != domain.ErrCannotDeactivateLastAdmin {
		t.Fatalf("expected ErrCannotDeactivateLastAdmin, got %v", err)
	}
}

func TestChangeOperatorRoleServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	clerk, err := buildOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid clerk, got %v", err)
	}

	session, err := domain.NewOperatorSession(
		"session-1",
		clerk.ID(),
		"hash:token",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin, clerk)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{
			session.TokenHash(): session,
		},
	}

	service := NewChangeOperatorRoleService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), ChangeOperatorRoleCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        clerk.ID(),
		Role:              domain.OperatorRoleAdmin.String(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	changedClerk := operatorRepository.byID[clerk.ID()]
	if changedClerk.Role() != domain.OperatorRoleAdmin {
		t.Fatalf("expected clerk to become admin, got %s", changedClerk.Role())
	}

	if _, ok := sessionRepository.byTokenHash[session.TokenHash()]; ok {
		t.Fatal("expected clerk sessions to be deleted")
	}
}

func TestChangeOperatorRoleServiceRejectsCurrentOperator(t *testing.T) {
	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewChangeOperatorRoleService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), ChangeOperatorRoleCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        admin.ID(),
		Role:              domain.OperatorRoleOperator.String(),
	})
	if err != domain.ErrCannotChangeCurrentOperatorRole {
		t.Fatalf("expected ErrCannotChangeCurrentOperatorRole, got %v", err)
	}
}

func TestChangeOperatorRoleServiceRejectsDemotingLastAdmin(t *testing.T) {
	currentAdmin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid current admin, got %v", err)
	}

	targetAdmin, err := buildOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid target admin, got %v", err)
	}

	inactiveCurrentAdmin, err := domain.ReconstituteOperator(
		currentAdmin.ID(),
		currentAdmin.RegistrationID(),
		currentAdmin.Alias(),
		currentAdmin.Rank(),
		currentAdmin.Role(),
		currentAdmin.PasswordHash(),
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive current admin, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(inactiveCurrentAdmin, targetAdmin)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewChangeOperatorRoleService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), ChangeOperatorRoleCommand{
		CurrentOperatorID: inactiveCurrentAdmin.ID(),
		OperatorID:        targetAdmin.ID(),
		Role:              domain.OperatorRoleOperator.String(),
	})
	if err != domain.ErrCannotDemoteLastAdmin {
		t.Fatalf("expected ErrCannotDemoteLastAdmin, got %v", err)
	}
}

func TestChangeOperatorRoleServiceRejectsInvalidRole(t *testing.T) {
	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	clerk, err := buildOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid clerk, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin, clerk)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewChangeOperatorRoleService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), ChangeOperatorRoleCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        clerk.ID(),
		Role:              "root",
	})
	if err != domain.ErrInvalidOperatorRole {
		t.Fatalf("expected ErrInvalidOperatorRole, got %v", err)
	}
}

func TestResetOperatorPasswordServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$old-admin-hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	clerk, err := buildOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$old-clerk-hash")
	if err != nil {
		t.Fatalf("expected valid clerk, got %v", err)
	}

	session, err := domain.NewOperatorSession(
		"session-1",
		clerk.ID(),
		"hash:token",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin, clerk)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{
			session.TokenHash(): session,
		},
	}

	service := NewResetOperatorPasswordService(
		operatorRepository,
		sessionRepository,
		fakePasswordHasher{hash: "$argon2id$new-hash"},
	)

	err = service.Execute(context.Background(), ResetOperatorPasswordCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        clerk.ID(),
		Password:          "correct horse battery staple",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updatedClerk := operatorRepository.byID[clerk.ID()]
	if updatedClerk.PasswordHash() != "$argon2id$new-hash" {
		t.Fatalf("expected new hash, got %s", updatedClerk.PasswordHash())
	}

	if _, ok := sessionRepository.byTokenHash[session.TokenHash()]; ok {
		t.Fatal("expected clerk sessions to be deleted")
	}
}

func TestResetOperatorPasswordServiceRejectsCurrentOperator(t *testing.T) {
	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewResetOperatorPasswordService(
		operatorRepository,
		sessionRepository,
		fakePasswordHasher{hash: "$argon2id$new-hash"},
	)

	err = service.Execute(context.Background(), ResetOperatorPasswordCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        admin.ID(),
		Password:          "correct horse battery staple",
	})
	if err != domain.ErrCannotResetCurrentOperatorPassword {
		t.Fatalf("expected ErrCannotResetCurrentOperatorPassword, got %v", err)
	}
}

func TestResetOperatorPasswordServiceRejectsWeakPassword(t *testing.T) {
	admin, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleAdmin, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid admin, got %v", err)
	}

	clerk, err := buildOperator("operator-2", "93541134780", "costa", domain.RankCorporal, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid clerk, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(admin, clerk)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewResetOperatorPasswordService(
		operatorRepository,
		sessionRepository,
		fakePasswordHasher{hash: "$argon2id$new-hash"},
	)

	err = service.Execute(context.Background(), ResetOperatorPasswordCommand{
		CurrentOperatorID: admin.ID(),
		OperatorID:        clerk.ID(),
		Password:          "short",
	})
	if err != domain.ErrWeakOperatorPassword {
		t.Fatalf("expected ErrWeakOperatorPassword, got %v", err)
	}
}

func TestGetOperatorAdminServiceExecute(t *testing.T) {
	createdAt := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	repository := &fakeOperatorRepository{
		summaries: []ports.OperatorSummary{
			{
				ID:             "operator-1",
				RegistrationID: mustRegistrationID(t, "52998224725"),
				Alias:          "silva",
				Rank:           domain.RankSergeant,
				Role:           domain.OperatorRoleAdmin,
				Active:         true,
				CreatedAt:      createdAt,
			},
		},
	}

	service := NewGetOperatorAdminService(repository)

	operator, err := service.Execute(context.Background(), GetOperatorAdminCommand{
		OperatorID: "operator-1",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if operator.RegistrationID.String() != "52998224725" {
		t.Fatalf("expected registration id 52998224725, got %s", operator.RegistrationID)
	}
}

func TestGetOperatorAdminServiceReturnsNotFound(t *testing.T) {
	repository := &fakeOperatorRepository{}
	service := NewGetOperatorAdminService(repository)

	_, err := service.Execute(context.Background(), GetOperatorAdminCommand{
		OperatorID: "missing",
	})
	if err != ports.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestReactivateOperatorServiceExecute(t *testing.T) {
	now := time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)

	inactiveOperator, err := domain.ReconstituteOperator(
		"operator-1",
		mustRegistrationID(t, "52998224725"),
		"silva",
		domain.RankSergeant,
		domain.OperatorRoleOperator,
		"$argon2id$hash",
		false,
	)
	if err != nil {
		t.Fatalf("expected valid inactive operator, got %v", err)
	}

	session, err := domain.NewOperatorSession(
		"session-1",
		inactiveOperator.ID(),
		"hash:token",
		"csrf-token",
		now.Add(time.Hour),
		now,
	)
	if err != nil {
		t.Fatalf("expected valid session, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(inactiveOperator)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{
			session.TokenHash(): session,
		},
	}

	service := NewReactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), ReactivateOperatorCommand{
		OperatorID: inactiveOperator.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	reactivatedOperator := operatorRepository.byID[inactiveOperator.ID()]
	if !reactivatedOperator.Active() {
		t.Fatal("expected operator to be active")
	}

	if _, ok := sessionRepository.byTokenHash[session.TokenHash()]; ok {
		t.Fatal("expected stale sessions to be deleted")
	}
}

func TestReactivateOperatorServiceIsNoOpForActiveOperator(t *testing.T) {
	activeOperator, err := buildOperator("operator-1", "52998224725", "silva", domain.RankSergeant, domain.OperatorRoleOperator, "$argon2id$hash")
	if err != nil {
		t.Fatalf("expected valid active operator, got %v", err)
	}

	operatorRepository := newFakeOperatorRepository(activeOperator)

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewReactivateOperatorService(operatorRepository, sessionRepository)

	err = service.Execute(context.Background(), ReactivateOperatorCommand{
		OperatorID: activeOperator.ID(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestReactivateOperatorServiceReturnsNotFound(t *testing.T) {
	operatorRepository := newFakeOperatorRepository()

	sessionRepository := &fakeOperatorSessionRepository{
		byTokenHash: map[string]domain.OperatorSession{},
	}

	service := NewReactivateOperatorService(operatorRepository, sessionRepository)

	err := service.Execute(context.Background(), ReactivateOperatorCommand{
		OperatorID: "missing",
	})
	if err != ports.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
