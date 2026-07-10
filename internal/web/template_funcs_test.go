package web

import "testing"

func TestQueryEscape(t *testing.T) {
	escaped := queryEscape("sergeant doe")

	if escaped != "sergeant+doe" {
		t.Fatalf("expected sergeant+doe, got %s", escaped)
	}
}

func TestQueryEscapeEscapesSpecialCharacters(t *testing.T) {
	escaped := queryEscape("529.982 doe & radio")

	expected := "529.982+doe+%26+radio"
	if escaped != expected {
		t.Fatalf("expected %s, got %s", expected, escaped)
	}
}
