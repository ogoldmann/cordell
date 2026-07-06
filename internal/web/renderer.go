package web

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed views/*.html views/pages/*.html
var templatesFS embed.FS

// Renderer renders HTML templates.
type Renderer struct {
	templates *template.Template
}

// NewRenderer parses and creates the HTML template renderer.
func NewRenderer() (*Renderer, error) {
	templates, err := template.ParseFS(
		templatesFS,
		"views/*.html",
		"views/pages/*.html",
	)
	if err != nil {
		return nil, err
	}

	return &Renderer{
		templates: templates,
	}, nil
}

// Render writes a named HTML template response.
func (r *Renderer) Render(w http.ResponseWriter, status int, name string, data any) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	return r.templates.ExecuteTemplate(w, name, data)
}
