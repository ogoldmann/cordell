package web

type pageActionView struct {
	Label string
	URL   string
}

type pageHeaderView struct {
	Eyebrow       string
	Title         string
	Description   string
	PrimaryAction *pageActionView
}

type emptyStateView struct {
	Title         string
	Description   string
	PrimaryAction *pageActionView
}

type formActionsView struct {
	PrimaryLabel              string
	SecondaryLabel            string
	SecondaryURL              string
	SaveAndCreateAnotherLabel string
	ShowSaveAndCreateAnother  bool
}

type feedbackMessageView struct {
	Kind        string
	Title       string
	Message     string
	ActionLabel string
	ActionURL   string
}

type sectionHeaderView struct {
	Eyebrow     string
	Title       string
	Description string
}

type detailFieldView struct {
	Label string
	Value string
}

func newPageAction(label string, url string) *pageActionView {
	if label == "" || url == "" {
		return nil
	}

	return &pageActionView{
		Label: label,
		URL:   url,
	}
}

func newPageHeader(eyebrow string, title string, description string, primaryAction *pageActionView) pageHeaderView {
	return pageHeaderView{
		Eyebrow:       eyebrow,
		Title:         title,
		Description:   description,
		PrimaryAction: primaryAction,
	}
}

func newEmptyState(title string, description string, primaryAction *pageActionView) emptyStateView {
	return emptyStateView{
		Title:         title,
		Description:   description,
		PrimaryAction: primaryAction,
	}
}

func newFormActions(primaryLabel string, secondaryLabel string, secondaryURL string) formActionsView {
	return formActionsView{
		PrimaryLabel:   primaryLabel,
		SecondaryLabel: secondaryLabel,
		SecondaryURL:   secondaryURL,
	}
}

func newFormActionsWithSaveAndCreateAnother(
	primaryLabel string,
	secondaryLabel string,
	secondaryURL string,
	saveAndCreateAnotherLabel string,
) formActionsView {
	return formActionsView{
		PrimaryLabel:              primaryLabel,
		SecondaryLabel:            secondaryLabel,
		SecondaryURL:              secondaryURL,
		SaveAndCreateAnotherLabel: saveAndCreateAnotherLabel,
		ShowSaveAndCreateAnother:  true,
	}
}

func newErrorFeedback(message string) *feedbackMessageView {
	if message == "" {
		return nil
	}

	return &feedbackMessageView{
		Kind:    "error",
		Title:   "Verifique as informações",
		Message: message,
	}
}

func newErrorFeedbackWithAction(message string, actionLabel string, actionURL string) *feedbackMessageView {
	if message == "" {
		return nil
	}

	return &feedbackMessageView{
		Kind:        "error",
		Title:       "Verifique as informações",
		Message:     message,
		ActionLabel: actionLabel,
		ActionURL:   actionURL,
	}
}

func newSuccessFeedback(message string) *feedbackMessageView {
	if message == "" {
		return nil
	}

	return &feedbackMessageView{
		Kind:    "success",
		Title:   "Operação concluída",
		Message: message,
	}
}

func newSectionHeader(eyebrow string, title string, description string) sectionHeaderView {
	return sectionHeaderView{
		Eyebrow:     eyebrow,
		Title:       title,
		Description: description,
	}
}

func newDetailField(label string, value string) detailFieldView {
	return detailFieldView{
		Label: label,
		Value: value,
	}
}
