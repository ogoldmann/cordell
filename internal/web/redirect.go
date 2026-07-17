package web

import (
	"net/http"
	"strings"
)

const defaultPostLoginRedirectPath = "/"

func currentRequestPath(r *http.Request) string {
	path := r.URL.EscapedPath()
	if path == "" {
		path = "/"
	}

	if r.URL.RawQuery != "" {
		path += "?" + r.URL.RawQuery
	}

	return path
}

func sanitizeReturnTo(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultPostLoginRedirectPath
	}

	if !strings.HasPrefix(value, "/") {
		return defaultPostLoginRedirectPath
	}

	if strings.HasPrefix(value, "//") {
		return defaultPostLoginRedirectPath
	}

	if strings.Contains(value, "\\") {
		return defaultPostLoginRedirectPath
	}

	return value
}
