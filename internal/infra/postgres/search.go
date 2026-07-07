package postgres

import (
	"strings"
	"unicode"
)

const noRegistrationSearchMatchPattern = "CORDELLNOREGISTRATIONSEARCHMATCH"

func buildTextSearchPattern(query string) string {
	return "%" + escapeLikePattern(strings.TrimSpace(query)) + "%"
}

func buildRegistrationSearchPattern(query string) string {
	var builder strings.Builder

	for _, char := range strings.TrimSpace(query) {
		if unicode.IsDigit(char) || unicode.IsLetter(char) {
			builder.WriteRune(unicode.ToUpper(char))
		}
	}

	if builder.Len() == 0 {
		return noRegistrationSearchMatchPattern
	}

	return "%" + escapeLikePattern(builder.String()) + "%"
}

func escapeLikePattern(value string) string {
	var builder strings.Builder

	for _, char := range value {
		switch char {
		case '\\', '%', '_':
			builder.WriteRune('\\')
		}

		builder.WriteRune(char)
	}

	return builder.String()
}
