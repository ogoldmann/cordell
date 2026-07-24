package app

import (
	"strings"
	"unicode"

	"cordell/internal/domain"
	"cordell/internal/ports"
)

// ExpandedSearchQuery contains the original query and expanded token groups.
type ExpandedSearchQuery struct {
	Original string
	Tokens   []ExpandedSearchToken
}

// ExpandedSearchToken contains one original search token and its equivalent terms.
type ExpandedSearchToken struct {
	Original string
	Terms    []string
}

// ExpandSearchQuery expands military and operational catalog tokens in a query.
func ExpandSearchQuery(query string) ExpandedSearchQuery {
	query = strings.TrimSpace(query)
	if query == "" {
		return ExpandedSearchQuery{}
	}

	rawTokens := searchTokens(query)
	expandedTokens := make([]ExpandedSearchToken, 0, len(rawTokens))

	for index, token := range rawTokens {
		terms := expandSearchToken(token)

		if index+1 < len(rawTokens) {
			nextToken := rawTokens[index+1]
			combinedToken := token + " " + nextToken

			for _, term := range expandSearchToken(combinedToken) {
				terms = appendUniqueSearchTerm(terms, term)
			}
		}

		expandedTokens = append(expandedTokens, ExpandedSearchToken{
			Original: token,
			Terms:    uniqueSearchTerms(terms),
		})
	}

	return ExpandedSearchQuery{
		Original: query,
		Tokens:   expandedTokens,
	}
}

// ToPortsSearchQuery converts an app search expansion into a repository-neutral search query.
func (q ExpandedSearchQuery) ToPortsSearchQuery() ports.SearchQuery {
	tokens := make([]ports.SearchToken, 0, len(q.Tokens))

	for _, token := range q.Tokens {
		tokens = append(tokens, ports.SearchToken{
			Original: token.Original,
			Terms:    append([]string(nil), token.Terms...),
		})
	}

	return ports.SearchQuery{
		Original: q.Original,
		Tokens:   tokens,
	}
}

// CanonicalSearchQuery replaces known operational tokens with their internal catalog values.
func CanonicalSearchQuery(query string) string {
	query = strings.TrimSpace(query)
	if query == "" {
		return ""
	}

	tokens := searchTokens(query)
	canonicalTokens := make([]string, 0, len(tokens))

	for index := 0; index < len(tokens); index++ {
		token := tokens[index]

		if index+1 < len(tokens) {
			combined := token + " " + tokens[index+1]
			if canonical, ok := canonicalCatalogValue(combined); ok {
				canonicalTokens = append(canonicalTokens, canonical)
				index++
				continue
			}
		}

		if canonical, ok := canonicalCatalogValue(token); ok {
			canonicalTokens = append(canonicalTokens, canonical)
			continue
		}

		canonicalTokens = append(canonicalTokens, token)
	}

	return strings.Join(canonicalTokens, " ")
}

func searchTokens(query string) []string {
	return strings.FieldsFunc(query, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func expandSearchToken(token string) []string {
	normalizedToken := normalizeSearchToken(token)

	termsByKey := map[string]string{}
	addSearchTerm(termsByKey, token)
	addSearchTerm(termsByKey, normalizedToken)

	for _, rank := range domain.RankOptions() {
		values := []string{
			rank.Value.String(),
			rank.Label,
			rank.Abbreviation,
		}
		values = append(values, rank.SearchTerms...)

		if searchTokenMatchesCatalogValue(normalizedToken, values...) {
			addSearchTerm(termsByKey, rank.Value.String())
			addSearchTerm(termsByKey, rank.Label)
			addSearchTerm(termsByKey, rank.Abbreviation)

			for _, searchTerm := range rank.SearchTerms {
				addSearchTerm(termsByKey, searchTerm)
			}
		}
	}

	for _, section := range domain.PersonnelSectionOptions() {
		values := []string{
			section.Value.String(),
			section.Label,
			section.Abbreviation,
		}
		values = append(values, section.SearchTerms...)

		if searchTokenMatchesCatalogValue(normalizedToken, values...) {
			addSearchTerm(termsByKey, section.Value.String())
			addSearchTerm(termsByKey, section.Label)
			addSearchTerm(termsByKey, section.Abbreviation)

			for _, searchTerm := range section.SearchTerms {
				addSearchTerm(termsByKey, searchTerm)
			}
		}
	}

	terms := make([]string, 0, len(termsByKey))
	for _, term := range termsByKey {
		terms = append(terms, term)
	}

	return terms
}

func searchTokenMatchesCatalogValue(token string, values ...string) bool {
	for _, value := range values {
		normalizedValue := normalizeSearchToken(value)

		if normalizedValue == token {
			return true
		}

		if strings.Contains(normalizedValue, token) {
			return true
		}
	}

	return false
}

func canonicalCatalogValue(value string) (string, bool) {
	normalizedValue := normalizeSearchToken(value)

	for _, rank := range domain.RankOptions() {
		values := []string{
			rank.Value.String(),
			rank.Label,
			rank.Abbreviation,
		}
		values = append(values, rank.SearchTerms...)

		for _, candidate := range values {
			if normalizeSearchToken(candidate) == normalizedValue {
				return rank.Value.String(), true
			}
		}
	}

	for _, section := range domain.PersonnelSectionOptions() {
		values := []string{
			section.Value.String(),
			section.Label,
			section.Abbreviation,
		}
		values = append(values, section.SearchTerms...)

		for _, candidate := range values {
			if normalizeSearchToken(candidate) == normalizedValue {
				return section.Value.String(), true
			}
		}
	}

	return "", false
}

func addSearchTerm(terms map[string]string, value string) {
	value = strings.TrimSpace(value)
	if value == "" {
		return
	}

	key := normalizeSearchToken(value)
	if key == "" {
		return
	}

	terms[key] = value
}

func appendUniqueSearchTerm(terms []string, term string) []string {
	term = strings.TrimSpace(term)
	if term == "" {
		return terms
	}

	normalizedTerm := normalizeSearchToken(term)
	for _, existing := range terms {
		if normalizeSearchToken(existing) == normalizedTerm {
			return terms
		}
	}

	return append(terms, term)
}

func uniqueSearchTerms(terms []string) []string {
	result := make([]string, 0, len(terms))
	for _, term := range terms {
		result = appendUniqueSearchTerm(result, term)
	}

	return result
}

func normalizeSearchToken(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))

	replacer := strings.NewReplacer(
		"º", "",
		"ª", "",
		"-", "",
		"_", "",
		"/", "",
		".", "",
		",", "",
		"(", "",
		")", "",
	)

	value = replacer.Replace(value)
	value = strings.Join(strings.Fields(value), "")

	return value
}
