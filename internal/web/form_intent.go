package web

import "net/http"

const (
	formIntentSave                 = "save"
	formIntentSaveAndCreateAnother = "save_and_create_another"
)

func formIntent(r *http.Request) string {
	intent := r.FormValue("intent")
	if intent == "" {
		return formIntentSave
	}

	return intent
}

func wantsSaveAndCreateAnother(r *http.Request) bool {
	return formIntent(r) == formIntentSaveAndCreateAnother
}
