package web

import (
	"errors"

	"cordell/internal/app"
	"cordell/internal/domain"
)

func personnelFeedbackFromError(err error) *feedbackMessageView {
	if err == nil {
		return nil
	}

	message := duplicatePersonnelErrorMessage(err)

	if existingID, ok := app.ExistingPersonnelIDFromDuplicateError(err); ok {
		return newErrorFeedbackWithAction(
			message,
			"Abrir militar existente",
			"/personnel/"+string(existingID),
		)
	}

	return newErrorFeedback(message)
}

func assetFeedbackFromError(err error) *feedbackMessageView {
	if err == nil {
		return nil
	}

	message := duplicateAssetErrorMessage(err)

	if existingID, ok := app.ExistingAssetIDFromDuplicateError(err); ok {
		return newErrorFeedbackWithAction(
			message,
			"Abrir material existente",
			"/assets/"+string(existingID),
		)
	}

	return newErrorFeedback(message)
}

func duplicatePersonnelErrorMessage(err error) string {
	if errors.Is(err, domain.ErrDuplicateRegistrationID) {
		return "Esta identidade já está cadastrada."
	}

	return humanizePersonnelError(err)
}

func duplicateAssetErrorMessage(err error) string {
	if errors.Is(err, domain.ErrDuplicateAssetName) {
		return "Este material já está cadastrado."
	}

	return humanizeAssetError(err)
}
