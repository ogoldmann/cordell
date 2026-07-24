package web

import "net/http"

type privateLayoutData struct {
	CSRFToken        string
	CurrentOperator  currentOperatorView
	HasOperator      bool
	ShowNavbarSearch bool
	ThemeSelector    themeSelectorView
	Breadcrumbs      []breadcrumbItemView
}

func newPrivateLayoutData(r *http.Request) privateLayoutData {
	data := privateLayoutData{
		CSRFToken:        csrfTokenFromContext(r),
		ShowNavbarSearch: true,
		ThemeSelector:    newThemeSelectorView(),
	}

	operator, ok := currentOperatorFromContext(r.Context())
	if ok {
		data.CurrentOperator = newCurrentOperatorView(operator)
		data.HasOperator = true
	}

	return data
}
