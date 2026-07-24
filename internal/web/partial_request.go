package web

import "net/http"

const partialRequestHeader = "X-Cordell-Partial"

func wantsPartialResponse(r *http.Request) bool {
	return r.Header.Get(partialRequestHeader) == "1" || r.URL.Query().Get("partial") == "1"
}
