package web

import (
	"net/http"
	"testing"
)

func TestIsSafeMethod(t *testing.T) {
	if !isSafeMethod(http.MethodGet) {
		t.Fatal("expected GET to be safe")
	}

	if isSafeMethod(http.MethodPost) {
		t.Fatal("expected POST not to be safe")
	}
}

func TestSameSecret(t *testing.T) {
	if !sameSecret("csrf-token", "csrf-token") {
		t.Fatal("expected matching tokens")
	}

	if sameSecret("csrf-token", "other-token") {
		t.Fatal("expected different tokens not to match")
	}

	if sameSecret("", "csrf-token") {
		t.Fatal("expected empty token not to match")
	}
}
