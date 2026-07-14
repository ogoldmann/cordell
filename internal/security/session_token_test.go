package security

import "testing"

func TestRandomSessionTokenGeneratorNewToken(t *testing.T) {
	generator := NewRandomSessionTokenGenerator(32)

	firstToken, err := generator.NewToken()
	if err != nil {
		t.Fatalf("expected no error generating first token, got %v", err)
	}

	secondToken, err := generator.NewToken()
	if err != nil {
		t.Fatalf("expected no error generating second token, got %v", err)
	}

	if firstToken == "" {
		t.Fatal("expected token not to be empty")
	}

	if firstToken == secondToken {
		t.Fatal("expected tokens to be different")
	}
}

func TestSHA256SessionTokenHasherHash(t *testing.T) {
	hasher := NewSHA256SessionTokenHasher()

	firstHash := hasher.Hash("token")
	secondHash := hasher.Hash("token")
	otherHash := hasher.Hash("other-token")

	if firstHash == "" {
		t.Fatal("expected hash not to be empty")
	}

	if firstHash != secondHash {
		t.Fatal("expected hash to be deterministic")
	}

	if firstHash == otherHash {
		t.Fatal("expected different token to produce different hash")
	}
}
