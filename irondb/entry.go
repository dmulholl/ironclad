package irondb


import (
    "encoding/base64"
    "github.com/dmulholland/ironclad/ironcrypt"
)


// An Entry object represents a single database record.
type Entry struct {
    Id int              `json:"id"`
    Active bool         `json:"active"`
    Title string        `json:"title"`
    Url string          `json:"url"`
    Username string     `json:"username"`
    Password string     `json:"password"`
    Tags []string       `json:"tags"`
    Notes string        `json:"notes"`
}


// NewEntry returns a new Entry object.
func NewEntry() *Entry {
    return &Entry{
        Active: true,
        Tags: make([]string, 0),
    }
}


// GetPassword returns the entry's (decrypted) password.
func (entry *Entry) GetPassword(key []byte) (string, error) {
    encrypted, err := base64.StdEncoding.DecodeString(entry.Password)
    if err != nil {
        return "", err
    }
    decrypted, err := ironcrypt.Decrypt(key, encrypted)
    if err != nil {
        return "", err
    }
    return string(decrypted), nil
}


// SetPassword sets the entry's password.
func (entry *Entry) SetPassword(key []byte, password string) error {
    encrypted, err := ironcrypt.Encrypt(key, []byte(password))
    if err != nil {
        return err
    }
    entry.Password = base64.StdEncoding.EncodeToString(encrypted)
    return nil
}
