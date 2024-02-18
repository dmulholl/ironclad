package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrRandBytes = errors.New("error generating random bytes")
)

// RandBytes generates an arbitrary-length slice of cryptographically-secure pseudorandom bytes.
func RandBytes(length int) ([]byte, error) {
	output := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, output)
	if err != nil {
		return nil, ErrRandBytes
	}
	return output, nil
}

// Key derives a secure encryption key from a password using the PBKDF2 key-derivation algorithm
// with an SHA-256 hash. The returned key will have keyLength bytes.
func Key(password string, salt []byte, iterationCount, keyLength int) []byte {
	return pbkdf2.Key([]byte(password), salt, iterationCount, keyLength, sha256.New)
}
