package web

import "net/http"

type privateLayoutData struct {
	CSRFToken        string
	CurrentOperator  currentOperatorView
	HasOperator      bool
	ShowNavbarSearch bool
	ThemeSelector    themeSelectorView
	Navbar           navbarView
	Breadcrumbs      []breadcrumbItemView
}

func newPrivateLayoutData(r *http.Request) privateLayoutData {
	data := privateLayoutData{
		CSRFToken:        csrfTokenFromContext(r),
		ShowNavbarSearch: true,
		ThemeSelector:    newCompactThemeSelectorView(),
	}

	operator, ok := currentOperatorFromContext(r.Context())
	if ok {
		data.CurrentOperator = newCurrentOperatorView(operator)
		data.HasOperator = true
	}

	data.Navbar = newNavbarView(r, data)

	return data
}

func (data *privateLayoutData) HideNavbarSearch() {
	data.ShowNavbarSearch = false
	data.Navbar.ShowSearch = false
}
