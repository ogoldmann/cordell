package postgres

import (
	"testing"
	"time"
)

func TestTimeToTimestamptzConvertsTimeToValidUTCTimestamp(t *testing.T) {
	location := time.FixedZone("BRT", -3*60*60)
	value := time.Date(2026, 7, 14, 10, 30, 0, 0, location)

	timestamp := timeToTimestamptz(value)

	if !timestamp.Valid {
		t.Fatal("expected timestamp to be valid")
	}

	expected := value.UTC()
	if !timestamp.Time.Equal(expected) {
		t.Fatalf("expected timestamp time %v, got %v", expected, timestamp.Time)
	}

	if timestamp.Time.Location() != time.UTC {
		t.Fatalf("expected timestamp location UTC, got %v", timestamp.Time.Location())
	}
}
