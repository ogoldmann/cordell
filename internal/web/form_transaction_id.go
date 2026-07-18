package web

import (
	"crypto/rand"
	"time"

	"cordell/internal/domain"

	"github.com/oklog/ulid/v2"
)

func newFormTransactionID() (domain.CustodyTransactionID, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now()), rand.Reader)
	if err != nil {
		return "", err
	}

	return domain.CustodyTransactionID(id.String()), nil
}
