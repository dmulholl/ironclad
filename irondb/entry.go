package irondb


import (
    "github.com/dmulholland/ironclad/ironcrypt"
)


// An Entry object represents a single database record.
type Entry struct {
    Id int              `json:"id"`
    Active bool         `json:"active"`
    Title string        `json:"title"`
    Url string          `json:"url"`
    Username string     `json:"username"`
    Password []byte     `json:"password"`
    Email string        `json:"email"`
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


// GetPassword returns an entry's (decrypted) password.
func (entry *Entry) GetPassword(key []byte) (string, error) {
    decrypted, err := ironcrypt.Decrypt(key, entry.Password)
    if err != nil {
        return "", err
    }
    return string(decrypted), nil
}


// SetPassword sets an entry's password.
func (entry *Entry) SetPassword(key []byte, password string) error {
    encrypted, err := ironcrypt.Encrypt(key, []byte(password))
    if err != nil {
        return err
    }
    entry.Password = encrypted
    return nil
}
