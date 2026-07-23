package web

import "cordell/internal/domain"

type personnelOptionView struct {
	ID                string
	DisplayName       string
	FullName          string
	SectionShortLabel string
	Label             string
	Selected          bool
}

func newPersonnelOptionView(personnel domain.Personnel) personnelOptionView {
	displayName := militaryDisplayName(personnel.Rank(), personnel.Alias())
	sectionShortLabel := personnel.Section().Abbreviation()

	return personnelOptionView{
		ID:                string(personnel.ID()),
		DisplayName:       displayName,
		FullName:          personnel.FullName(),
		SectionShortLabel: sectionShortLabel,
		Label:             newPersonnelOptionLabel(displayName, personnel.FullName(), sectionShortLabel),
	}
}

func newPersonnelOptionLabel(displayName string, fullName string, sectionShortLabel string) string {
	label := displayName

	if fullName != "" {
		label += " — " + fullName
	}

	if sectionShortLabel != "" {
		label += " — " + sectionShortLabel
	}

	return label
}
