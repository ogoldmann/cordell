package web

import (
	"log/slog"
	"net/http"
)

type notFoundPageData struct {
	privateLayoutData
	UsePrivateLayout bool
	Title            string
	Heading          string
	Description      string
	HomeURL          string
	HomeLabel        string
	SearchURL        string
	SearchLabel      string
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	s.renderNotFound(w, r)
}

func (s *Server) renderNotFound(w http.ResponseWriter, r *http.Request) {
	data := notFoundPageData{
		Title:       "Página não encontrada",
		Heading:     "Página não encontrada",
		Description: "A página solicitada não existe, foi movida ou você não possui um caminho válido para acessá-la.",
		HomeURL:     "/",
		HomeLabel:   "Ir para o painel",
		SearchURL:   "/search",
		SearchLabel: "Pesquisar no Cordell",
	}

	if _, ok := currentOperatorFromContext(r.Context()); ok {
		data.privateLayoutData = newPrivateLayoutData(r)
		data.UsePrivateLayout = true
	}

	if err := s.renderer.Render(w, http.StatusNotFound, "not_found.html", data); err != nil {
		s.logger.Error("failed to render not found page", slog.String("error", err.Error()))
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (s *Server) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	data := notFoundPageData{
		Title:       "Ação não permitida",
		Heading:     "Ação não permitida",
		Description: "O caminho existe, mas não aceita este tipo de requisição.",
		HomeURL:     "/",
		HomeLabel:   "Ir para o painel",
		SearchURL:   "/search",
		SearchLabel: "Pesquisar no Cordell",
	}

	if _, ok := currentOperatorFromContext(r.Context()); ok {
		data.privateLayoutData = newPrivateLayoutData(r)
		data.UsePrivateLayout = true
	}

	if err := s.renderer.Render(w, http.StatusMethodNotAllowed, "not_found.html", data); err != nil {
		s.logger.Error("failed to render method not allowed page", slog.String("error", err.Error()))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
