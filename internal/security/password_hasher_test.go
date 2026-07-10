package security

import (
	"strings"
	"testing"
)

func testArgon2idParams() Argon2idParams {
	return Argon2idParams{
		Memory:      64,
		Iterations:  1,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func TestArgon2idPasswordHasherHashAndVerify(t *testing.T) {
	hasher := NewArgon2idPasswordHasher(testArgon2idParams())

	hash, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	if !strings.HasPrefix(hash, "$argon2id$v=19$") {
		t.Fatalf("expected argon2id hash prefix, got %s", hash)
	}

	matches, err := hasher.Verify("correct horse battery staple", hash)
	if err != nil {
		t.Fatalf("expected no error verifying password, got %v", err)
	}

	if !matches {
		t.Fatal("expected password to match hash")
	}
}

func TestArgon2idPasswordHasherRejectsWrongPassword(t *testing.T) {
	hasher := NewArgon2idPasswordHasher(testArgon2idParams())

	hash, err := hasher.Hash("correct horse battery staple")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}

	matches, err := hasher.Verify("wrong password", hash)
	if err != nil {
		t.Fatalf("expected no error verifying password, got %v", err)
	}

	if matches {
		t.Fatal("expected password not to match hash")
	}
}

func TestArgon2idPasswordHasherRejectsInvalidHash(t *testing.T) {
	hasher := NewArgon2idPasswordHasher(testArgon2idParams())

	_, err := hasher.Verify("password", "invalid-hash")
	if err != ErrInvalidPasswordHash {
		t.Fatalf("expected ErrInvalidPasswordHash, got %v", err)
	}
}
