package ids

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestULIDGeneratorNewID(t *testing.T) {
	generator := NewULIDGenerator()

	id, err := generator.NewID()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(id) != ulid.EncodedSize {
		t.Fatalf("expected ULID length %d, got %d", ulid.EncodedSize, len(id))
	}

	if _, err := ulid.ParseStrict(id); err != nil {
		t.Fatalf("expected valid ULID, got %v", err)
	}
}

func TestULIDGeneratorCreatesDifferentIDs(t *testing.T) {
	generator := NewULIDGenerator()

	firstID, err := generator.NewID()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	secondID, err := generator.NewID()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if firstID == secondID {
		t.Fatal("expected different ULIDs")
	}
}
