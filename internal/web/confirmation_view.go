package web

type confirmationPageData struct {
	privateLayoutData
	Title             string
	Kicker            string
	Heading           string
	Description       string
	Warning           string
	ConfirmLabel      string
	CancelLabel       string
	ConfirmAction     string
	CancelURL         string
	ConfirmationStyle string
}
