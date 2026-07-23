package web

import (
	"strings"

	"cordell/internal/domain"
)

func militaryDisplayName(rank domain.Rank, alias string) string {
	alias = strings.TrimSpace(alias)

	rankLabel := rank.Abbreviation()
	if rankLabel == "" {
		rankLabel = rank.Label()
	}

	if alias == "" {
		return rankLabel
	}

	return rankLabel + " " + alias
}
