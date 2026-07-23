package web

import "testing"

func TestNewPersonnelOptionLabelIncludesSection(t *testing.T) {
	got := newPersonnelOptionLabel("Sd John", "John Doe", "S1")

	want := "Sd John — John Doe — S1"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}
