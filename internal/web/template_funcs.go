package web

import (
	"html/template"
	"net/url"
)

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"queryEscape": queryEscape,
	}
}

func queryEscape(value string) string {
	return url.QueryEscape(value)
}
