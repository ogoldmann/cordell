package web

import (
	"net/http/httptest"
	"testing"
)

func TestOperatorInitials(t *testing.T) {
	tests := []struct {
		name        string
		displayName string
		want        string
	}{
		{name: "empty", displayName: "", want: "OP"},
		{name: "single", displayName: "Silva", want: "S"},
		{name: "rank and alias", displayName: "Sd Silva", want: "SS"},
		{name: "multi", displayName: "3º Sgt Oliveira", want: "3O"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := operatorInitials(test.displayName)
			if got != test.want {
				t.Fatalf("expected %q, got %q", test.want, got)
			}
		})
	}
}

func TestNavbarKeepsHomeLinkInvisibleOnHome(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)

	layout := privateLayoutData{
		ShowNavbarSearch: true,
		ThemeSelector:    newThemeSelectorView(),
		CurrentOperator: currentOperatorView{
			DisplayName:        "Sd Silva",
			RoleLabel:          "Administrador",
			CanManageOperators: true,
		},
	}

	navbar := newNavbarView(request, layout)

	var homeLink navbarLinkView
	for _, link := range navbar.Links {
		if link.Label == "Home" {
			homeLink = link
			break
		}
	}

	if !homeLink.Show {
		t.Fatal("expected Home link to remain rendered for layout stability")
	}

	if !homeLink.Invisible {
		t.Fatal("expected Home link to be invisible on dashboard")
	}
}

func TestNavbarShowsHomeLinkOutsideHome(t *testing.T) {
	request := httptest.NewRequest("GET", "/personnel", nil)

	layout := privateLayoutData{
		ShowNavbarSearch: true,
		ThemeSelector:    newThemeSelectorView(),
		CurrentOperator: currentOperatorView{
			DisplayName:        "Sd Silva",
			RoleLabel:          "Administrador",
			CanManageOperators: true,
		},
	}

	navbar := newNavbarView(request, layout)

	if !navbar.ShowHomeLink {
		t.Fatal("expected home link to be shown outside dashboard")
	}

	var homeLink navbarLinkView
	for _, link := range navbar.Links {
		if link.Label == "Home" {
			homeLink = link
			break
		}
	}

	if !homeLink.Show {
		t.Fatal("expected Home link to be rendered outside dashboard")
	}

	if homeLink.Invisible {
		t.Fatal("expected Home link to be visible outside dashboard")
	}
}
