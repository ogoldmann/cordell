package web

import "time"

func formatDateTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}

	return value.Format("2006-01-02 15:04:05 MST")
}
