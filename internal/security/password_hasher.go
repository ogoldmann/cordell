package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const argon2idVersion = 19

// ErrInvalidPasswordHash is returned when a password hash cannot be parsed.
var ErrInvalidPasswordHash = errors.New("invalid password hash")

// Argon2idParams contains Argon2id password hashing parameters.
type Argon2idParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

// DefaultArgon2idParams returns the default password hashing parameters.
func DefaultArgon2idParams() Argon2idParams {
	return Argon2idParams{
		Memory:      19 * 1024,
		Iterations:  2,
		Parallelism: 1,
		SaltLength:  16,
		KeyLength:   32,
	}
}

// Argon2idPasswordHasher hashes and verifies passwords using Argon2id.
type Argon2idPasswordHasher struct {
	params Argon2idParams
}

// NewArgon2idPasswordHasher creates an Argon2idPasswordHasher.
func NewArgon2idPasswordHasher(params Argon2idParams) *Argon2idPasswordHasher {
	return &Argon2idPasswordHasher{
		params: params,
	}
}

// NewDefaultArgon2idPasswordHasher creates an Argon2idPasswordHasher with default parameters.
func NewDefaultArgon2idPasswordHasher() *Argon2idPasswordHasher {
	return NewArgon2idPasswordHasher(DefaultArgon2idParams())
}

// Hash hashes a password.
func (h *Argon2idPasswordHasher) Hash(password string) (string, error) {
	salt, err := randomBytes(h.params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		[]byte(password),
		salt,
		h.params.Iterations,
		h.params.Memory,
		h.params.Parallelism,
		h.params.KeyLength,
	)

	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2idVersion,
		h.params.Memory,
		h.params.Iterations,
		h.params.Parallelism,
		encodedSalt,
		encodedHash,
	), nil
}

// Verify verifies a password against an encoded hash.
func (h *Argon2idPasswordHasher) Verify(password string, encodedHash string) (bool, error) {
	params, salt, expectedHash, err := decodeArgon2idHash(encodedHash)
	if err != nil {
		return false, err
	}

	actualHash := argon2.IDKey(
		[]byte(password),
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		uint32(len(expectedHash)),
	)

	return subtle.ConstantTimeCompare(expectedHash, actualHash) == 1, nil
}

func randomBytes(length uint32) ([]byte, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func decodeArgon2idHash(encodedHash string) (Argon2idParams, []byte, []byte, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return Argon2idParams{}, nil, nil, ErrInvalidPasswordHash
	}

	if parts[1] != "argon2id" {
		return Argon2idParams{}, nil, nil, ErrInvalidPasswordHash
	}

	version, err := parseVersion(parts[2])
	if err != nil {
		return Argon2idParams{}, nil, nil, err
	}

	if version != argon2idVersion {
		return Argon2idParams{}, nil, nil, ErrInvalidPasswordHash
	}

	params, err := parseParams(parts[3])
	if err != nil {
		return Argon2idParams{}, nil, nil, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return Argon2idParams{}, nil, nil, ErrInvalidPasswordHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return Argon2idParams{}, nil, nil, ErrInvalidPasswordHash
	}

	return params, salt, hash, nil
}

func parseVersion(value string) (int, error) {
	version, found := strings.CutPrefix(value, "v=")
	if !found {
		return 0, ErrInvalidPasswordHash
	}

	parsedVersion, err := strconv.Atoi(version)
	if err != nil {
		return 0, ErrInvalidPasswordHash
	}

	return parsedVersion, nil
}

func parseParams(value string) (Argon2idParams, error) {
	values := strings.Split(value, ",")

	if len(values) != 3 {
		return Argon2idParams{}, ErrInvalidPasswordHash
	}

	params := Argon2idParams{}

	for _, item := range values {
		key, rawValue, found := strings.Cut(item, "=")
		if !found {
			return Argon2idParams{}, ErrInvalidPasswordHash
		}

		switch key {
		case "m":
			memory, err := strconv.ParseUint(rawValue, 10, 32)
			if err != nil {
				return Argon2idParams{}, ErrInvalidPasswordHash
			}
			params.Memory = uint32(memory)

		case "t":
			iterations, err := strconv.ParseUint(rawValue, 10, 32)
			if err != nil {
				return Argon2idParams{}, ErrInvalidPasswordHash
			}
			params.Iterations = uint32(iterations)

		case "p":
			parallelism, err := strconv.ParseUint(rawValue, 10, 8)
			if err != nil {
				return Argon2idParams{}, ErrInvalidPasswordHash
			}
			params.Parallelism = uint8(parallelism)

		default:
			return Argon2idParams{}, ErrInvalidPasswordHash
		}
	}

	if params.Memory == 0 || params.Iterations == 0 || params.Parallelism == 0 {
		return Argon2idParams{}, ErrInvalidPasswordHash
	}

	return params, nil
}
