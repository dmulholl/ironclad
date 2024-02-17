/*
Package fileio provides read/write access to the content of encrypted files.
*/
package fileio

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"os"

	"github.com/dmulholl/ironclad/internal/ironcrypt"
	"github.com/dmulholl/ironclad/internal/ironcrypt/aes"
)

// Length of the key derivation salt in bytes.
const SaltLength = 32

// Number of iterations used by the key derivation function.
const PBKDFIterations = 100000

// Load reads data from an encrypted file.
func Load(filename, password string) (data []byte, err error) {

	// Load the file content.
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	// Split the file content into its component parts.
	salt := content[:SaltLength]
	ciphertext := content[SaltLength:]

	// Use the password and salt to regenerate the file encryption key.
	key := ironcrypt.Key(password, salt, PBKDFIterations, aes.KeySize)

	// Use the key to decrypt the ciphertext.
	plaintext, err := aes.Decrypt(ciphertext, key)
	if err != nil {
		return
	}

	// Unzip the plaintext.
	unzipper, err := gzip.NewReader(bytes.NewBuffer(plaintext))
	if err != nil {
		return
	}
	unzipped, err := ioutil.ReadAll(unzipper)
	if err != nil {
		return
	}
	unzipper.Close()

	return unzipped, nil
}

// Save compresses and writes a slice of data to disk as an encrypted file.
func Save(filename, password string, data []byte) (err error) {

	// Zip the plaintext before encrypting.
	var zipped bytes.Buffer
	zipper := gzip.NewWriter(&zipped)
	_, err = zipper.Write(data)
	if err != nil {
		return
	}
	zipper.Close()

	// Generate a random salt.
	salt, err := ironcrypt.RandBytes(SaltLength)
	if err != nil {
		return
	}

	// Use the password and salt to generate a file encryption key.
	key := ironcrypt.Key(password, salt, PBKDFIterations, aes.KeySize)

	// Encrypt the data using the key.
	ciphertext, err := aes.Encrypt(zipped.Bytes(), key)
	if err != nil {
		return
	}

	// Write the salt and ciphertext to a temporary output file.
	file, err := os.Create(filename + ".new")
	if err != nil {
		return
	}
	_, err = file.Write(salt)
	if err != nil {
		return
	}
	_, err = file.Write(ciphertext)
	if err != nil {
		return
	}
	file.Close()

	// Delete any older file instance and rename the temporary file.
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	}
	err = os.Rename(filename+".new", filename)
	if err != nil {
		return
	}

	return
}
