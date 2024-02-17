/*
Package ironcrypt is a symmetric encryption library built from secure, industry-standard components.
*/
package ironcrypt

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

// Key derives a secure encryption key from a password using the PBKDF2 key derivation algorithm
// with an SHA-256 hash.
func Key(password string, salt []byte, iterations, size int) []byte {
	return pbkdf2.Key([]byte(password), salt, iterations, size, sha256.New)
}
