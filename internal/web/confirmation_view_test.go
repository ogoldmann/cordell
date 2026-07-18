package web

import "testing"

func TestConfirmationPageDataCanBeBuilt(t *testing.T) {
	data := confirmationPageData{
		Title:             "Deactivate asset",
		Kicker:            "Asset lifecycle",
		Heading:           "Deactivate Radio?",
		Description:       "This asset will be removed from normal checkout workflows.",
		Warning:           "Deactivation does not settle current custody.",
		ConfirmLabel:      "Deactivate asset",
		CancelLabel:       "Cancel",
		ConfirmAction:     "/assets/asset-1/deactivate",
		CancelURL:         "/assets/asset-1",
		ConfirmationStyle: "warning",
	}

	if data.ConfirmAction == "" {
		t.Fatal("expected confirm action")
	}

	if data.CancelURL == "" {
		t.Fatal("expected cancel URL")
	}
}
