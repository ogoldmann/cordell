package web

import "time"

func formatDateTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Local().Format("02/01/2006 15:04")
}
