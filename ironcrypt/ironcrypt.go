/*
    Package ironcrypt is an easy-to-use symmetric encryption library built
    from secure, industry-standard components.

    Encryption is performed using 256-bit AES in CBC mode. Message padding
    is handled transparently using the PKCS #7 padding scheme. Message
    authentication is handled transparently using the HMAC-SHA-256 algorithm.

    Encryption requires a 64-byte key. The first 32 bytes are used as the
    cipher key; the last 32 bytes are used as the authentication key.

    Note that this package acts as a wrapper around Go's native cryptographic
    libraries. It does not reimplement the cryptographic primitives it uses.

    The output of the Encrypt() function is structured as follows:

        |-- IV --|-- Payload --|-- HMAC --|
            16         16n          32

        IV: initialization vector
        Payload: encrypted plaintext
        HMAC: authentication code

    Lengths are given in bytes. The length of the payload will be a multiple of
    the AES block size (16 bytes). (Note that the padding algorithm appends a
    null block to the plaintext when the length of the plaintext is a multiple
    of the block size. This padding is automatically removed during decryption.)

    The Key() function can be used to derive a suitable 64-byte encryption key
    from a password. It uses the PBKDF2 scheme with an SHA-256 hash.
*/
package ironcrypt


import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/hmac"
    "crypto/sha256"
    "errors"
    "io"
    "golang.org/x/crypto/pbkdf2"
)


const (
    HMACSize = 32 // Output size of HMAC-SHA-256
    CKeySize = 32 // Cipher key size for AES-256
    AKeySize = 32 // Authentication key size for HMAC-SHA-256
    KeySize  = 64 // CKeySize + AKeySize
)


var (
    ErrInvalidPadding = errors.New("invalid padding length")
    ErrInvalidKeySize = errors.New("invalid key length")
    ErrIV = errors.New("error generating initialization vector")
    ErrInvalidCiphertext = errors.New("invalid ciphertext length")
    ErrInvalidHMAC = errors.New("invalid authentication code")
    ErrRandBytes = errors.New("error generating random bytes")
)


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


// Encrypt encrypts a plaintext slice using a 64-byte key.
func Encrypt(key, plaintext []byte) ([]byte, error) {

    // Check the key size.
    if len(key) != KeySize {
        return nil, ErrInvalidKeySize
    }

    // Create a random initialization vector.
    iv := make([]byte, aes.BlockSize)
    _, err := io.ReadFull(rand.Reader, iv)
    if err != nil {
        return nil, ErrIV
    }

    // Pad the plaintext so its length is a multiple of the AES block size.
    pplaintext := pad(plaintext)
    ciphertext := make([]byte, len(pplaintext))

    // Generate the ciphertext and prepend the initialization vector.
    c, _ := aes.NewCipher(key[:CKeySize])
    encrypter := cipher.NewCBCEncrypter(c, iv)
    encrypter.CryptBlocks(ciphertext, pplaintext)
    ciphertext = append(iv, ciphertext...)

    // Generate and append the HMAC.
    hash := hmac.New(sha256.New, key[CKeySize:])
    hash.Write(ciphertext)
    ciphertext = hash.Sum(ciphertext)

    return ciphertext, nil
}


// Decrypt decrypts a ciphertext slice using a 64-byte key.
func Decrypt(key, ciphertext []byte) ([]byte, error) {

    // Check the key size.
    if len(key) != KeySize {
        return nil, ErrInvalidKeySize
    }

    // HMAC-SHA-256 returns a 32-byte MAC so the overall message length
    // should be a multiple of the AES block size.
    if (len(ciphertext) % aes.BlockSize) != 0 {
        return nil, ErrInvalidCiphertext
    }

    // The ciphertext must include at least an IV block, a message block,
    // and a double-length block of HMAC.
    if len(ciphertext) < (4 * aes.BlockSize) {
        return nil, ErrInvalidCiphertext
    }

    // Strip the HMAC from the end of the ciphertext.
    macIndex := len(ciphertext) - HMACSize
    oldMac := ciphertext[macIndex:]
    ciphertext = ciphertext[:macIndex]
    pplaintext := make([]byte, macIndex - aes.BlockSize)

    // Check the message authentication codes.
    hash := hmac.New(sha256.New, key[CKeySize:])
    hash.Write(ciphertext)
    newMac := hash.Sum(nil)
    if !hmac.Equal(newMac, oldMac) {
        return nil, ErrInvalidHMAC
    }

    // Decrypt the ciphertext.
    c, _ := aes.NewCipher(key[:CKeySize])
    decrypter := cipher.NewCBCDecrypter(c, ciphertext[:aes.BlockSize])
    decrypter.CryptBlocks(pplaintext, ciphertext[aes.BlockSize:])

    // Strip the padding bytes from the plaintext.
    plaintext, err := unpad(pplaintext)
    if err != nil {
        return nil, ErrInvalidPadding
    }

    return plaintext, nil
}


// RandBytes generates an arbitrary-length slice of cryptographically-secure
// pseudorandom bytes.
func RandBytes(length int) ([]byte, error) {
    output := make([]byte, length)
    _, err := io.ReadFull(rand.Reader, output)
    if err != nil {
        return nil, ErrRandBytes
    }
    return output, nil
}


// Key generates a 64-byte encryption key from a password using the PBKDF2
// key derivation algorithm.
func Key(password, salt []byte, iterations int) []byte {
    return pbkdf2.Key([]byte(password), salt, iterations, KeySize, sha256.New)
}


// LenCipher returns the length in bytes of the ciphertext that will be
// returned for a given plaintext length.
func LenCipher(lenPlain int) (length int) {

    // The initialization vector is a single AES block.
    length += aes.BlockSize

    // The padding algorithm adds an extra null block to the plaintext if the
    // length of the plaintext is a multiple of the block size.
    if lenPlain % aes.BlockSize == 0 {
        length += lenPlain + aes.BlockSize
    } else {
        length += lenPlain
    }

    // Authentication code.
    length += HMACSize

    return length
}
