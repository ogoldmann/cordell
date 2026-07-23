package web

import "strconv"

func editCountLabel(editCount int) string {
	if editCount <= 0 {
		return ""
	}

	label := "Editado " + strconv.Itoa(editCount) + " vez"
	if editCount != 1 {
		label += "es"
	}

	return label
}
