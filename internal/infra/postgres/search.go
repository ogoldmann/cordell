package postgres

import (
	"strings"
	"unicode"
)

const noDigitSearchMatchPattern = "__cordell_no_digit_search_match__"

func buildTextSearchPatterns(query string) []string {
	tokens := searchTokens(query)
	patterns := make([]string, 0, len(tokens))

	for _, token := range tokens {
		patterns = append(patterns, "%"+escapeLikePattern(token)+"%")
	}

	return patterns
}

func buildDigitSearchPatterns(query string) []string {
	tokens := searchTokens(query)
	patterns := make([]string, 0, len(tokens))

	for _, token := range tokens {
		patterns = append(patterns, buildDigitSearchPattern(token))
	}

	return patterns
}

func searchTokens(query string) []string {
	return strings.Fields(strings.TrimSpace(query))
}

func buildDigitSearchPattern(token string) string {
	var builder strings.Builder

	for _, char := range token {
		if unicode.IsDigit(char) {
			builder.WriteRune(char)
		}
	}

	if builder.Len() == 0 {
		return noDigitSearchMatchPattern
	}

	return "%" + builder.String() + "%"
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
