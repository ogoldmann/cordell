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
}

func newDashboardOperationalDockView() dashboardOperationalDockView {
	return dashboardOperationalDockView{
		Groups: []dashboardOperationalGroupView{
			{
				Title: "CADASTRAR",
				Actions: []dashboardOperationalActionView{
					{Label: "Militar", URL: "/personnel/new"},
					{Label: "Material", URL: "/assets/new"},
				},
			},
			{
				Title: "REGISTRAR",
				Actions: []dashboardOperationalActionView{
					{Label: "Cautela", URL: "/custody/checkouts/new"},
					{Label: "Descautela", URL: "/custody/returns/new"},
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
