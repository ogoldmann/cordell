package domain

import "testing"

func TestNewRegistrationIDAcceptsFormattedCPF(t *testing.T) {
	registrationID, err := NewRegistrationID("529.982.247-25")
	if err != nil {
		t.Fatalf("expected valid registration id, got %v", err)
	}

	if registrationID.String() != "52998224725" {
		t.Fatalf("expected normalized registration id 52998224725, got %s", registrationID.String())
	}
}

func TestNewRegistrationIDRejectsNonNumericValue(t *testing.T) {
	_, err := NewRegistrationID("abc")
	if err != ErrInvalidRegistrationID {
		t.Fatalf("expected ErrInvalidRegistrationID, got %v", err)
	}
}
