package web

type themeOptionView struct {
	Value       string
	Label       string
	Description string
	Icon        iconView
}

type themeSelectorView struct {
	Options []themeOptionView
	Compact bool
}

func newThemeSelectorView() themeSelectorView {
	return themeSelectorView{
		Options: []themeOptionView{
			{
				Value:       "light",
				Label:       "Claro",
				Description: "Interface clara e limpa.",
				Icon:        newIcon("sun", "size-4"),
			},
			{
				Value:       "dark",
				Label:       "Escuro",
				Description: "Interface escura e operacional.",
				Icon:        newIcon("moon", "size-4"),
			},
			{
				Value:       "sepia",
				Label:       "Sépia",
				Description: "Interface quente e confortável.",
				Icon:        newIcon("monitor", "size-4"),
			},
		},
	}
}

func newCompactThemeSelectorView() themeSelectorView {
	view := newThemeSelectorView()
	view.Compact = true

	return view
}
