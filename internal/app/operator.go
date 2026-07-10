package app

import (
	"context"
	"strings"
	"unicode/utf8"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

const minOperatorPasswordLength = 15

// CreateOperatorCommand contains the input data required to create an operator.
type CreateOperatorCommand struct {
	Username string
	Password string
}

// CreateOperatorService handles operator creation.
type CreateOperatorService struct {
	operatorRepository ports.OperatorRepository
	idGenerator        ports.IDGenerator
	passwordHasher     ports.PasswordHasher
}

// NewCreateOperatorService creates a CreateOperatorService with its dependencies.
func NewCreateOperatorService(
	operatorRepository ports.OperatorRepository,
	idGenerator ports.IDGenerator,
	passwordHasher ports.PasswordHasher,
) *CreateOperatorService {
	return &CreateOperatorService{
		operatorRepository: operatorRepository,
		idGenerator:        idGenerator,
		passwordHasher:     passwordHasher,
	}
}

// Execute creates an operator account.
func (s *CreateOperatorService) Execute(ctx context.Context, cmd CreateOperatorCommand) (domain.Operator, error) {
	if strings.TrimSpace(cmd.Password) == "" {
		return domain.Operator{}, domain.ErrEmptyOperatorPassword
	}

	if utf8.RuneCountInString(cmd.Password) < minOperatorPasswordLength {
		return domain.Operator{}, domain.ErrWeakOperatorPassword
	}

	id, err := s.idGenerator.NewID()
	if err != nil {
		return domain.Operator{}, err
	}

	passwordHash, err := s.passwordHasher.Hash(cmd.Password)
	if err != nil {
		return domain.Operator{}, err
	}

	operator, err := domain.NewOperator(
		domain.OperatorID(id),
		cmd.Username,
		passwordHash,
	)
	if err != nil {
		return domain.Operator{}, err
	}

	if err := s.operatorRepository.Save(ctx, operator); err != nil {
		return domain.Operator{}, err
	}

	return operator, nil
}
