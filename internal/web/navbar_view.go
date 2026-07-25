package web

import (
	"net/http"
	"strings"
)

type navbarView struct {
	BrandName      string
	SectionName    string
	DeveloperLabel string
	GithubLabel    string
	GithubURL      string
	ShowHomeLink   bool
	ShowSearch     bool
	Search         searchBarView
	Links          []navbarLinkView
	Operator       navbarOperatorView
	ThemeSelector  themeSelectorView
	CSRFToken      string
}

type navbarLinkView struct {
	Label  string
	URL    string
	Icon   iconView
	Active bool
	Show   bool
}

type navbarOperatorView struct {
	DisplayName string
	RoleLabel   string
	Initials    string
}

func newNavbarView(r *http.Request, layout privateLayoutData) navbarView {
	path := r.URL.Path
	isHome := path == "/"

	operator := navbarOperatorView{
		DisplayName: layout.CurrentOperator.DisplayName,
		RoleLabel:   layout.CurrentOperator.RoleLabel,
		Initials:    operatorInitials(layout.CurrentOperator.DisplayName),
	}

	showHomeLink := !isHome

	return navbarView{
		BrandName:      brandName,
		SectionName:    brandSectionName,
		DeveloperLabel: brandDeveloper,
		GithubLabel:    brandGithubLabel,
		GithubURL:      brandGithubURL,
		ShowHomeLink:   showHomeLink,
		ShowSearch:     layout.ShowNavbarSearch,
		Search:         newNavbarSearchBar(),
		Links: []navbarLinkView{
			{
				Label:  "Home",
				URL:    "/",
				Icon:   newIcon("home", "size-4"),
				Active: isHome,
				Show:   showHomeLink,
			},
			{
				Label:  "Militares",
				URL:    "/personnel",
				Icon:   newIcon("users", "size-4"),
				Active: strings.HasPrefix(path, "/personnel"),
				Show:   true,
			},
			{
				Label:  "Materiais",
				URL:    "/assets",
				Icon:   newIcon("package", "size-4"),
				Active: strings.HasPrefix(path, "/assets"),
				Show:   true,
			},
			{
				Label:  "Transações",
				URL:    "/custody/transactions",
				Icon:   newIcon("clipboard-list", "size-4"),
				Active: strings.HasPrefix(path, "/custody"),
				Show:   true,
			},
			{
				Label:  "Admin",
				URL:    "/admin/operators",
				Icon:   newIcon("shield", "size-4"),
				Active: strings.HasPrefix(path, "/admin"),
				Show:   layout.CurrentOperator.CanManageOperators,
			},
		},
		Operator:      operator,
		ThemeSelector: layout.ThemeSelector,
		CSRFToken:     layout.CSRFToken,
	}
}

func operatorInitials(displayName string) string {
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		return "OP"
	}

	parts := strings.Fields(displayName)
	if len(parts) == 1 {
		value := []rune(parts[0])
		if len(value) == 0 {
			return "OP"
		}
		return strings.ToUpper(string(value[0]))
	}

	first := []rune(parts[0])
	last := []rune(parts[len(parts)-1])

	if len(first) == 0 || len(last) == 0 {
		return "OP"
	}

	return strings.ToUpper(string(first[0]) + string(last[0]))
}
