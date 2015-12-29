/*
    Package ironio handles reading from and writing to encrypted files.
*/
package ironio


import (
    "github.com/dmulholland/ironclad/ironcrypt"
    "os"
    "io/ioutil"
    "compress/gzip"
    "bytes"
)


// Length of the key derivation salt in bytes.
const SaltLength = 32


// Number of iterations used by the key derivation function.
const PBKDFIterations = 10000


// Length of the encrypted 64-byte master key.
// Made up of: IV + 4 x blocks of data + 1 x block of padding + HMAC.
const EncryptedKeyLength = 16 + 64 + 16 + 32


// Load reads data from an encrypted file.
func Load(password, filename string) (data, key []byte, err error) {

    // Load the file content.
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return
    }

    // Split the file content into its component parts.
    salt := content[:SaltLength]
    encryptedkey := content[SaltLength:SaltLength + EncryptedKeyLength]
    ciphertext := content[SaltLength + EncryptedKeyLength:]

    // Use the password and salt to regenerate the file encryption key.
    filekey := ironcrypt.Key([]byte(password), salt, PBKDFIterations)

    // Use the file key to decrypt the master key.
    key, err = ironcrypt.Decrypt(filekey, encryptedkey)
    if err != nil {
        return
    }

    // Use the master key to decrypt the ciphertext.
    plaintext, err := ironcrypt.Decrypt(key, ciphertext)
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

    return unzipped, key, nil
}


// Save writes data to disk as an encrypted file.
func Save(data, key []byte, password, filename string) (err error) {

    // Zip the plaintext before encrypting.
    var zipped bytes.Buffer
    zipper := gzip.NewWriter(&zipped)
    _, err = zipper.Write(data)
    if err != nil {
        return
    }
    zipper.Close()

    // Encrypt the data store using the master key.
    ciphertext, err := ironcrypt.Encrypt(key, zipped.Bytes())
    if err != nil {
        return
    }

    // Generate a random salt.
    salt, err := ironcrypt.RandBytes(SaltLength)
    if err != nil {
        return
    }

    // Use the password and salt to generate a new encryption key.
    filekey := ironcrypt.Key([]byte(password), salt, PBKDFIterations)

    // Use the new key to encrypt the master key.
    encryptedkey, err := ironcrypt.Encrypt(filekey, key)
    if err != nil {
        return
    }

    // Write the salt, encrypted master key, and encrypted data store to file.
    file, err := os.Create(filename + ".new")
    if err != nil {
        return
    }
    _, err = file.Write(salt)
    if err != nil {
        return
    }
    _, err = file.Write(encryptedkey)
    if err != nil {
        return
    }
    _, err = file.Write(ciphertext)
    if err != nil {
        return
    }
    file.Close()

    // Delete any older file instance and rename the new file.
    if _, err := os.Stat(filename); err == nil {
        err = os.Remove(filename)
        if err != nil {
            return err
        }
    }
    err = os.Rename(filename + ".new", filename)
    if err != nil {
        return
    }

    return
}
