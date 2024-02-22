package ioutils

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"

	"github.com/dmulholl/ironclad/internal/crypto"
	"github.com/dmulholl/ironclad/internal/crypto/aes"
)

// Length of the key derivation salt in bytes.
const SaltLength = 32

// Number of iterations used by the key derivation function.
const PBKDFIterations = 100000

// Encrypt compresses and encrypts plaintext using the specified password.
func Encrypt(password string, plaintext []byte) ([]byte, error) {
	// Compress the plaintext before encrypting.
	var zipped bytes.Buffer
	zipper := gzip.NewWriter(&zipped)

	_, err := zipper.Write(plaintext)
	if err != nil {
		return nil, fmt.Errorf("failed to zip plaintext: %w", err)
	}

	err = zipper.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close zipper: %w", err)
	}

	// Generate a random salt.
	salt, err := crypto.RandBytes(SaltLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random salt: %w", err)
	}

	// Use the password and salt to generate a file encryption key.
	key := crypto.Key(password, salt, PBKDFIterations, aes.KeySize)

	// Encrypt the data using the key.
	ciphertext, err := aes.Encrypt(zipped.Bytes(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt zipped plaintext: %w", err)
	}

	output := make([]byte, len(salt)+len(ciphertext))
	copy(output, salt)
	copy(output[len(salt):], ciphertext)
	return output, nil
}

// Decrypt decrypts output from the Encrypt() function.
func Decrypt(password string, data []byte) ([]byte, error) {
	// Split the data into its component parts.
	salt := data[:SaltLength]
	ciphertext := data[SaltLength:]

	// Use the password and salt to regenerate the file encryption key.
	key := crypto.Key(password, salt, PBKDFIterations, aes.KeySize)

	// Use the key to decrypt the ciphertext.
	plaintext, err := aes.Decrypt(ciphertext, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	// Unzip the plaintext.
	unzipper, err := gzip.NewReader(bytes.NewBuffer(plaintext))
	if err != nil {
		return nil, fmt.Errorf("failed to unzip plaintext: %w", err)
	}
	defer unzipper.Close()

	unzipped, err := io.ReadAll(unzipper)
	if err != nil {
		return nil, fmt.Errorf("failed to unzip plaintext: %w", err)
	}

	return unzipped, nil
}

// Load reads and decrypts data from an encrypted file.
func Load(filename, password string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return Decrypt(password, content)
}

// Save compresses and writes plaintext as an encrypted file.
func Save(filename, password string, plaintext []byte) error {
	content, err := Encrypt(password, plaintext)
	if err != nil {
		return err
	}

	// Write the encrypted content to a temporary output file.
	file, err := os.Create(filename + ".ironclad.temp")
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	file.Close()

	// Delete any existing file instance.
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		if err != nil {
			return fmt.Errorf("failed to delete existing file: %w", err)
		}
	}

	err = os.Rename(filename+".ironclad.temp", filename)
	if err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}
