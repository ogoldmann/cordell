package domain

import "strings"

// PersonnelID uniquely identifies a personnel record.
type PersonnelID string

// Personnel represents a person who can receive assets under custody.
type Personnel struct {
	id       PersonnelID
	fullName string
	active   bool
}

// NewPersonnel creates an active Personnel with validated required fields.
func NewPersonnel(id PersonnelID, fullName string) (Personnel, error) {
	if id == "" {
		return Personnel{}, ErrEmptyPersonnelID
	}

	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return Personnel{}, ErrEmptyPersonnelName
	}

	return Personnel{
		id:       id,
		fullName: fullName,
		active:   true,
	}, nil
}

// ID returns the personnel identifier.
func (p Personnel) ID() PersonnelID {
	return p.id
}

// FullName returns the personnel full name.
func (p Personnel) FullName() string {
	return p.fullName
}

// Active reports whether the personnel record is active.
func (p Personnel) Active() bool {
	return p.active
}
