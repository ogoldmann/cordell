package web

import "testing"

func TestCurrentBreadcrumb(t *testing.T) {
	item := currentBreadcrumb("Militares")

	if item.Label != "Militares" {
		t.Fatalf("expected label Militares, got %q", item.Label)
	}

	if !item.Current {
		t.Fatal("expected current breadcrumb")
	}

	if item.URL != "" {
		t.Fatalf("expected empty URL, got %q", item.URL)
	}
}

func TestPersonnelBreadcrumb(t *testing.T) {
	item := personnelBreadcrumb()

	if item.Label != "Militares" {
		t.Fatalf("expected Militares, got %q", item.Label)
	}

	if item.URL != "/personnel" {
		t.Fatalf("expected /personnel, got %q", item.URL)
	}
}
