/*
    Package ironio provides read and write access to the content of encrypted
    files.
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


// Load reads data from an encrypted file.
func Load(password, filename string) (data []byte, err error) {

    // Load the file content.
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        return
    }

    // Split the file content into its component parts.
    salt := content[:SaltLength]
    ciphertext := content[SaltLength:]

    // Use the password and salt to regenerate the file encryption key.
    filekey := ironcrypt.Key([]byte(password), salt, PBKDFIterations)

    // Use the key to decrypt the ciphertext.
    plaintext, err := ironcrypt.Decrypt(filekey, ciphertext)
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


// Save writes data to disk as an encrypted file.
func Save(password, filename string, data []byte) (err error) {

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

    // Use the password and salt to generate a new file encryption key.
    filekey := ironcrypt.Key([]byte(password), salt, PBKDFIterations)

    // Encrypt the data using the key.
    ciphertext, err := ironcrypt.Encrypt(filekey, zipped.Bytes())
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
    err = os.Rename(filename + ".new", filename)
    if err != nil {
        return
    }

    return
}
