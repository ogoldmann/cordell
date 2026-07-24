package web

type iconView struct {
	Name       string
	Class      string
	Label      string
	Decorative bool
}

func newIcon(name string, class string) iconView {
	return iconView{
		Name:       name,
		Class:      class,
		Decorative: true,
	}
}

func newLabeledIcon(name string, class string, label string) iconView {
	return iconView{
		Name:       name,
		Class:      class,
		Label:      label,
		Decorative: false,
	}
}
