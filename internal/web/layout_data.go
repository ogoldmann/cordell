package web

import "net/http"

type privateLayoutData struct {
	CSRFToken string
}

func newPrivateLayoutData(r *http.Request) privateLayoutData {
	return privateLayoutData{
		CSRFToken: csrfTokenFromContext(r),
	}
}
