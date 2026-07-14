package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// RandomSessionTokenGenerator generates cryptographically random session tokens.
type RandomSessionTokenGenerator struct {
	byteLength int
}

// NewRandomSessionTokenGenerator creates a RandomSessionTokenGenerator.
func NewRandomSessionTokenGenerator(byteLength int) *RandomSessionTokenGenerator {
	return &RandomSessionTokenGenerator{
		byteLength: byteLength,
	}
}

// NewDefaultRandomSessionTokenGenerator creates a RandomSessionTokenGenerator with a secure default length.
func NewDefaultRandomSessionTokenGenerator() *RandomSessionTokenGenerator {
	return NewRandomSessionTokenGenerator(32)
}

// NewToken creates a new random session token.
func (g *RandomSessionTokenGenerator) NewToken() (string, error) {
	bytes := make([]byte, g.byteLength)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

// SHA256SessionTokenHasher hashes session tokens using SHA-256.
type SHA256SessionTokenHasher struct{}

// NewSHA256SessionTokenHasher creates a SHA256SessionTokenHasher.
func NewSHA256SessionTokenHasher() *SHA256SessionTokenHasher {
	return &SHA256SessionTokenHasher{}
}

// Hash hashes a raw session token.
func (h *SHA256SessionTokenHasher) Hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
