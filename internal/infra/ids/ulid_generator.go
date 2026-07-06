package ids

import (
	"crypto/rand"
	"time"

	"cordell/internal/ports"

	"github.com/oklog/ulid/v2"
)

// ULIDGenerator creates lexicographically sortable ULID identifiers.
type ULIDGenerator struct{}

// NewULIDGenerator creates a ULIDGenerator.
func NewULIDGenerator() *ULIDGenerator {
	return &ULIDGenerator{}
}

// NewID creates a new ULID string using the current UTC time and cryptographic randomness.
func (g *ULIDGenerator) NewID() (string, error) {
	id, err := ulid.New(ulid.Timestamp(time.Now().UTC()), rand.Reader)
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

var _ ports.IDGenerator = (*ULIDGenerator)(nil)
