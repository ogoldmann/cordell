package web

import (
	"crypto/rand"
	"math/big"
)

type dashboardOperationalDockView struct {
	Groups []dashboardOperationalGroupView
}

type dashboardOperationalGroupView struct {
	Title   string
	Actions []dashboardOperationalActionView
}

type dashboardOperationalActionView struct {
	Label string
	URL   string
	Icon  iconView
}

func newDashboardOperationalDockView() dashboardOperationalDockView {
	return dashboardOperationalDockView{
		Groups: []dashboardOperationalGroupView{
			{
				Title: "CADASTRAR",
				Actions: []dashboardOperationalActionView{
					{
						Label: "Militar",
						URL:   "/personnel/new",
						Icon:  newIcon("users", "size-4"),
					},
					{
						Label: "Material",
						URL:   "/assets/new",
						Icon:  newIcon("package", "size-4"),
					},
				},
			},
			{
				Title: "REGISTRAR",
				Actions: []dashboardOperationalActionView{
					{
						Label: "Cautela",
						URL:   "/custody/checkouts/new",
						Icon:  newIcon("plus", "size-4"),
					},
					{
						Label: "Descautela",
						URL:   "/custody/returns/new",
						Icon:  newIcon("receipt-text", "size-4"),
					},
				},
			},
		},
	}
}

func randomDashboardWelcomePhrase() string {
	phrases := []string{
		"Segura, aqui é Pelotão de Segurança.",
		"Excelente, vamos ao trabalho.",
		"Excelente.",
		"Mais um dia. Tudo sob controle.",
		"Vamos colocar tudo em ordem.",
		"Hoje também, no controle.",
	}

	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(phrases))))
	if err != nil {
		return phrases[0]
	}

	return phrases[index.Int64()]
}
