package web

import (
	"strings"

	"cordell/internal/domain"
)

func militaryDisplayName(rank domain.Rank, alias string) string {
	alias = strings.TrimSpace(alias)
	if alias == "" {
		return rank.Label()
	}

	return rank.Label() + " " + alias
}
