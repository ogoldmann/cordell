package web

const (
	searchBarVariantDefault = "default"
	searchBarVariantHero    = "hero"
	searchBarVariantCompact = "compact"
)

type searchBarView struct {
	ID                 string
	Name               string
	Value              string
	Label              string
	Placeholder        string
	Variant            string
	Autofocus          bool
	ButtonLabel        string
	ButtonAriaLabel    string
	ShowButtonText     bool
	ShellClass         string
	InputClass         string
	ButtonClass        string
	IconClass          string
	MobilePanelEnabled bool
}

func newSearchBar(id string, value string, placeholder string) searchBarView {
	return searchBarView{
		ID:              id,
		Name:            "q",
		Value:           value,
		Label:           "Pesquisar",
		Placeholder:     placeholder,
		Variant:         searchBarVariantDefault,
		ButtonLabel:     "Pesquisar",
		ButtonAriaLabel: "Pesquisar",
		ShowButtonText:  false,
		ShellClass:      "cordell-search cordell-search-default",
		InputClass:      "cordell-search-input",
		ButtonClass:     "cordell-search-button",
		IconClass:       "size-5",
	}
}

func newHeroSearchBar(id string, value string, placeholder string) searchBarView {
	view := newSearchBar(id, value, placeholder)
	view.Variant = searchBarVariantHero
	view.Autofocus = true
	view.ShellClass = "cordell-search cordell-search-hero"
	view.IconClass = "size-6"

	return view
}

func newCompactSearchBar(id string, value string, placeholder string) searchBarView {
	view := newSearchBar(id, value, placeholder)
	view.Variant = searchBarVariantCompact
	view.ShellClass = "cordell-search cordell-search-compact"
	view.IconClass = "size-4"

	return view
}

func newSearchBarWithLabel(id string, value string, placeholder string, label string) searchBarView {
	view := newSearchBar(id, value, placeholder)
	view.Label = label

	return view
}
