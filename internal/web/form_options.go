package web

import "cordell/internal/domain"

type personnelOptionView struct {
	ID          string
	DisplayName string
	FullName    string
}

func newPersonnelOptionView(personnel domain.Personnel) personnelOptionView {
	return personnelOptionView{
		ID:          string(personnel.ID()),
		DisplayName: militaryDisplayName(personnel.Rank(), personnel.Alias()),
		FullName:    personnel.FullName(),
	}
}
