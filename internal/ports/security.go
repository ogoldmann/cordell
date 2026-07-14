package ports

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password string, encodedHash string) (bool, error)
}

// SessionTokenGenerator generates raw session tokens.
type SessionTokenGenerator interface {
	NewToken() (string, error)
}

// SessionTokenHasher hashes raw session tokens for server-side storage.
type SessionTokenHasher interface {
	Hash(token string) string
}
