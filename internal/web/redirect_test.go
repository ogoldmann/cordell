package web

import "testing"

func TestSanitizeReturnToAcceptsLocalPath(t *testing.T) {
	got := sanitizeReturnTo("/personnel?q=doe")

	if got != "/personnel?q=doe" {
		t.Fatalf("expected /personnel?q=doe, got %s", got)
	}
}

func TestSanitizeReturnToRejectsExternalURL(t *testing.T) {
	got := sanitizeReturnTo("https://evil.example/login")

	if got != "/" {
		t.Fatalf("expected /, got %s", got)
	}
}

func TestSanitizeReturnToRejectsProtocolRelativeURL(t *testing.T) {
	got := sanitizeReturnTo("//evil.example/login")

	if got != "/" {
		t.Fatalf("expected /, got %s", got)
	}
}

func TestSanitizeReturnToRejectsBackslash(t *testing.T) {
	got := sanitizeReturnTo("/\\evil")

	if got != "/" {
		t.Fatalf("expected /, got %s", got)
	}
}
