package web

import (
	"crypto/rand"
	"time"

	"cordell/internal/domain"

	"github.com/oklog/ulid/v2"
)

func newFormULID() (string, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func newFormTransactionID() (domain.CustodyTransactionID, error) {
	id, err := newFormULID()
	if err != nil {
		return "", err
	}

	return domain.CustodyTransactionID(id), nil
}

func newFormCorrectionID() (domain.CustodyCorrectionID, error) {
	id, err := newFormULID()
	if err != nil {
		return "", err
	}

	return domain.CustodyCorrectionID(id), nil
}
