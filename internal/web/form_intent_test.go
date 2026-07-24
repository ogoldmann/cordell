package web

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestFormIntentDefaultsToSave(t *testing.T) {
	request := httptest.NewRequest(http.MethodPost, "/", nil)

	if got := formIntent(request); got != formIntentSave {
		t.Fatalf("expected default intent %q, got %q", formIntentSave, got)
	}
}

func TestWantsSaveAndCreateAnother(t *testing.T) {
	form := url.Values{}
	form.Set("intent", formIntentSaveAndCreateAnother)

	request := httptest.NewRequest(
		http.MethodPost,
		"/",
		strings.NewReader(form.Encode()),
	)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err := request.ParseForm(); err != nil {
		t.Fatalf("expected form parse, got %v", err)
	}

	if !wantsSaveAndCreateAnother(request) {
		t.Fatal("expected save and create another intent")
	}
}
