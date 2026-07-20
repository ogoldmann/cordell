package web

import "strconv"

func editCountLabel(editCount int) string {
	if editCount <= 0 {
		return ""
	}

	label := "Edited " + strconv.Itoa(editCount) + " time"
	if editCount != 1 {
		label += "s"
	}

	return label
}
