package app

import (
	"errors"

	"cordell/internal/domain"
)

type DuplicatePersonnelRegistrationIDError struct {
	ExistingPersonnelID domain.PersonnelID
}

func (e DuplicatePersonnelRegistrationIDError) Error() string {
	return domain.ErrDuplicateRegistrationID.Error()
}

func (e DuplicatePersonnelRegistrationIDError) Unwrap() error {
	return domain.ErrDuplicateRegistrationID
}

type DuplicateAssetNameError struct {
	ExistingAssetID domain.AssetID
}

func (e DuplicateAssetNameError) Error() string {
	return domain.ErrDuplicateAssetName.Error()
}

func (e DuplicateAssetNameError) Unwrap() error {
	return domain.ErrDuplicateAssetName
}

func ExistingPersonnelIDFromDuplicateError(err error) (domain.PersonnelID, bool) {
	var duplicateErr DuplicatePersonnelRegistrationIDError
	if errors.As(err, &duplicateErr) {
		return duplicateErr.ExistingPersonnelID, true
	}

	return "", false
}

func ExistingAssetIDFromDuplicateError(err error) (domain.AssetID, bool) {
	var duplicateErr DuplicateAssetNameError
	if errors.As(err, &duplicateErr) {
		return duplicateErr.ExistingAssetID, true
	}

	return "", false
}
