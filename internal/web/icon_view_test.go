package web

import (
	"bytes"
	"html/template"
	"strings"
	"testing"
)

func TestNewIcon(t *testing.T) {
	icon := newIcon("search", "size-5")

	if icon.Name != "search" {
		t.Fatalf("expected search icon, got %q", icon.Name)
	}

	if icon.Class != "size-5" {
		t.Fatalf("expected size-5 class, got %q", icon.Class)
	}

	if !icon.Decorative {
		t.Fatal("expected decorative icon")
	}
}

func TestNewLabeledIcon(t *testing.T) {
	icon := newLabeledIcon("search", "size-5", "Pesquisar")

	if icon.Name != "search" {
		t.Fatalf("expected search icon, got %q", icon.Name)
	}

	if icon.Label != "Pesquisar" {
		t.Fatalf("expected label Pesquisar, got %q", icon.Label)
	}

	if icon.Decorative {
		t.Fatal("expected non-decorative icon")
	}
}

func TestIconTemplateRendersSearchSVG(t *testing.T) {
	templates, err := template.New("").Funcs(templateFuncs()).ParseFS(
		templatesFS,
		"views/*.html",
		"views/pages/*.html",
	)
	if err != nil {
		t.Fatalf("expected templates to parse, got %v", err)
	}

	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "icon", newIcon("search", "size-5")); err != nil {
		t.Fatalf("expected icon template to render, got %v", err)
	}

	body := buffer.String()

	if !strings.Contains(body, "<svg") {
		t.Fatalf("expected rendered icon to contain svg, got %q", body)
	}

	if !strings.Contains(body, "size-5") {
		t.Fatalf("expected rendered icon to contain size-5, got %q", body)
	}
}
