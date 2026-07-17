package domain

import (
	"strings"
	"unicode"
)

// RegistrationID identifies personnel through an external registration document.
type RegistrationID string

// NewRegistrationID creates a validated RegistrationID.
//
// The current validation mechanism accepts Brazilian CPF-like identifiers,
// but the domain concept remains registration_id.
func NewRegistrationID(value string) (RegistrationID, error) {
	if strings.TrimSpace(value) == "" {
		return "", ErrEmptyRegistrationID
	}

	normalized := NormalizeRegistrationID(value)
	if normalized == "" {
		return "", ErrInvalidRegistrationID
	}

	if !isValidCPF(normalized) {
		return "", ErrInvalidRegistrationID
	}

	return RegistrationID(normalized), nil
}

// String returns the registration identifier value.
func (id RegistrationID) String() string {
	return string(id)
}

// NormalizeRegistrationID keeps only digits from a registration identifier.
func NormalizeRegistrationID(value string) string {
	var builder strings.Builder

	for _, char := range strings.TrimSpace(value) {
		if unicode.IsDigit(char) {
			builder.WriteRune(char)
		}
	}

	return builder.String()
}

func isValidCPF(value string) bool {
	if len(value) != 11 {
		return false
	}

	if hasOnlyRepeatedDigits(value) {
		return false
	}

	firstDigit := calculateCPFCheckDigit(value[:9], 10)
	if firstDigit != int(value[9]-'0') {
		return false
	}

	secondDigit := calculateCPFCheckDigit(value[:10], 11)
	return secondDigit == int(value[10]-'0')
}

func hasOnlyRepeatedDigits(value string) bool {
	for i := 1; i < len(value); i++ {
		if value[i] != value[0] {
			return false
		}
	}

	return true
}

func calculateCPFCheckDigit(value string, weight int) int {
	sum := 0

	for i := 0; i < len(value); i++ {
		digit := int(value[i] - '0')
		sum += digit * (weight - i)
	}

	remainder := sum % 11
	if remainder < 2 {
		return 0
	}

	return 11 - remainder
}
