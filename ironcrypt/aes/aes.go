/*
    Package ironcrypt/aes is a symmetric encryption library built from secure,
    industry-standard components.

    Encryption is performed using 256-bit AES in CBC mode. Message padding
    is handled transparently using the PKCS #7 padding scheme. Message
    authentication is handled transparently using the HMAC-SHA-256 algorithm.

    Encryption requires a 32-byte key.

    Note that this package acts as a wrapper around Go's native cryptographic
    libraries. It does not reimplement the cryptographic primitives it uses.

    The output of the Encrypt() function is structured as follows:

        |-- IV --|-- Payload --|-- HMAC --|
            16         16n          32

        IV: initialization vector
        Payload: encrypted plaintext
        HMAC: authentication code

    Lengths are given in bytes. The length of the payload will be a multiple
    of the AES block size (16 bytes). (Note that the padding algorithm
    appends a null block to the plaintext when the length of the plaintext is
    a multiple of the block size. This padding is automatically removed
    during decryption.)
*/
package aes


import (
    "io"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/hmac"
    "crypto/sha256"
    "errors"
)


const (
    KeySize  = 32 // Size in bytes of encryption key.
    HMACSize = 32 // Size in bytes of authentication code.
)


var (
    ErrInvalidPadding = errors.New("invalid padding length")
    ErrInvalidKeySize = errors.New("invalid key size")
    ErrInvalidIV = errors.New("error generating initialization vector")
    ErrInvalidCipherLength = errors.New("invalid ciphertext length")
    ErrInvalidHMAC = errors.New("authentication failed")
)


// Encrypt encrypts a plaintext slice using a 32-byte key.
func Encrypt(plaintext, key []byte) ([]byte, error) {

    // Check the key size.
    if len(key) != KeySize {
        return nil, ErrInvalidKeySize
    }

    // Create a random initialization vector.
    iv := make([]byte, aes.BlockSize)
    _, err := io.ReadFull(rand.Reader, iv)
    if err != nil {
        return nil, ErrInvalidIV
    }

    // Pad the plaintext so its length is a multiple of the AES block size.
    pplaintext := pad(plaintext)
    ciphertext := make([]byte, len(pplaintext))

    // Generate the ciphertext and prepend the initialization vector.
    c, _ := aes.NewCipher(key)
    encrypter := cipher.NewCBCEncrypter(c, iv)
    encrypter.CryptBlocks(ciphertext, pplaintext)
    ciphertext = append(iv, ciphertext...)

    // Generate and append the HMAC.
    hash := hmac.New(sha256.New, key)
    hash.Write(ciphertext)
    ciphertext = hash.Sum(ciphertext)

    return ciphertext, nil
}


// Decrypt decrypts a ciphertext slice using a 32-byte key.
func Decrypt(ciphertext, key []byte) ([]byte, error) {

    // Check the key size.
    if len(key) != KeySize {
        return nil, ErrInvalidKeySize
    }

    // HMAC-SHA-256 returns a 32-byte MAC so the overall message length
    // should be a multiple of the AES block size.
    if (len(ciphertext) % aes.BlockSize) != 0 {
        return nil, ErrInvalidCipherLength
    }

    // The ciphertext must include at least an IV block, a message block,
    // and a double-length block of HMAC.
    if len(ciphertext) < (4 * aes.BlockSize) {
        return nil, ErrInvalidCipherLength
    }

    // Strip the HMAC from the end of the ciphertext.
    index := len(ciphertext) - HMACSize
    cipherhmac := ciphertext[index:]
    ciphertext = ciphertext[:index]
    pplaintext := make([]byte, index - aes.BlockSize)

    // Check the message authentication codes.
    hash := hmac.New(sha256.New, key)
    hash.Write(ciphertext)
    computedhmac := hash.Sum(nil)
    if !hmac.Equal(computedhmac, cipherhmac) {
        return nil, ErrInvalidHMAC
    }

    // Decrypt the ciphertext.
    c, _ := aes.NewCipher(key)
    decrypter := cipher.NewCBCDecrypter(c, ciphertext[:aes.BlockSize])
    decrypter.CryptBlocks(pplaintext, ciphertext[aes.BlockSize:])

    // Strip the padding bytes from the plaintext.
    plaintext, err := unpad(pplaintext)
    if err != nil {
        return nil, ErrInvalidPadding
    }

    return plaintext, nil
}


// Pad a message so its length is a multiple of the AES block size.
func pad(message []byte) []byte {
    lenPadding := aes.BlockSize - (len(message) % aes.BlockSize)
    for i := 0; i < lenPadding; i++ {
        message = append(message, byte(lenPadding))
    }
    return message
}


// Strip padding bytes from a message.
func unpad(message []byte) ([]byte, error) {
    if len(message) == 0 {
        return nil, ErrInvalidPadding
    }

    lenPadding := message[len(message) - 1]
    if lenPadding == 0 || lenPadding > aes.BlockSize {
        return nil, ErrInvalidPadding
    }

    for i := len(message) - 1; i > len(message) - int(lenPadding) - 1; i-- {
        if message[i] != lenPadding {
            return nil, ErrInvalidPadding
        }
    }

    return message[:len(message) - int(lenPadding)], nil
}
