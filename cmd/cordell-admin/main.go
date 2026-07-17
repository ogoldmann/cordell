package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"syscall"

	"cordell/internal/app"
	"cordell/internal/config"
	"cordell/internal/domain"
	"cordell/internal/infra/ids"
	"cordell/internal/infra/postgres"
	"cordell/internal/infra/postgres/db"
	"cordell/internal/security"

	"golang.org/x/term"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	if err := run(context.Background(), os.Args[1:]); err != nil {
		logger.Error("admin command failed", "error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		printUsage()
		return nil
	}

	switch args[0] {
	case "create-operator":
		return runCreateOperator(ctx, args[1:])
	case "help", "-h", "--help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runCreateOperator(ctx context.Context, args []string) error {
	flags := flag.NewFlagSet("create-operator", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	registrationID := flags.String("registration-id", "", "operator registration id")
	alias := flags.String("alias", "", "operator alias")
	rank := flags.String("rank", "", "operator rank")
	role := flags.String("role", domain.OperatorRoleAdmin.String(), "operator role: admin or operator")

	if err := flags.Parse(args); err != nil {
		return err
	}

	password, err := promptPassword()
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	pool, err := postgres.OpenPool(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	queries := db.New(pool)

	operatorRepository := postgres.NewOperatorRepository(queries)
	idGenerator := ids.NewULIDGenerator()
	passwordHasher := security.NewDefaultArgon2idPasswordHasher()

	service := app.NewCreateOperatorService(
		operatorRepository,
		idGenerator,
		passwordHasher,
	)

	operator, err := service.Execute(ctx, app.CreateOperatorCommand{
		RegistrationID: *registrationID,
		Alias:          *alias,
		Rank:           *rank,
		Role:           *role,
		Password:       password,
	})
	if err != nil {
		return humanizeCreateOperatorError(err)
	}

	fmt.Fprintf(
		os.Stdout,
		"Operator created successfully: %s %s (%s, %s)\n",
		operator.Rank().Label(),
		operator.Alias(),
		operator.RegistrationID().String(),
		operator.Role().Label(),
	)

	return nil
}

func promptPassword() (string, error) {
	fmt.Fprint(os.Stderr, "Password: ")

	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}

	fmt.Fprint(os.Stderr, "Confirm password: ")

	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}

	password := string(passwordBytes)
	confirmPassword := string(confirmBytes)

	if password != confirmPassword {
		return "", errors.New("passwords do not match")
	}

	if strings.TrimSpace(password) == "" {
		return "", domain.ErrEmptyOperatorPassword
	}

	return password, nil
}

func humanizeCreateOperatorError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmptyRegistrationID):
		return errors.New("registration id is required")
	case errors.Is(err, domain.ErrInvalidRegistrationID):
		return errors.New("registration id is invalid")
	case errors.Is(err, domain.ErrDuplicateRegistrationID):
		return errors.New("registration id is already registered")
	case errors.Is(err, domain.ErrEmptyOperatorAlias):
		return errors.New("alias is required")
	case errors.Is(err, domain.ErrInvalidOperatorRank):
		return errors.New("rank is required")
	case errors.Is(err, domain.ErrEmptyOperatorPassword):
		return errors.New("password is required")
	case errors.Is(err, domain.ErrWeakOperatorPassword):
		return errors.New("password must have at least 15 characters")
	case errors.Is(err, domain.ErrEmptyOperatorRole):
		return errors.New("role is required")
	case errors.Is(err, domain.ErrInvalidOperatorRole):
		return errors.New("role must be either admin or operator")
	default:
		return err
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Cordell admin commands:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  create-operator -registration-id <id> -alias <alias> -rank <rank> [-role admin|operator]")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Examples:")
	fmt.Fprintln(os.Stderr, "  cordell-admin create-operator -registration-id 52998224725 -alias silva -rank sergeant -role admin")
	fmt.Fprintln(os.Stderr, "  go run ./cmd/cordell-admin create-operator -registration-id 93541134780 -alias costa -rank corporal -role operator")
}
