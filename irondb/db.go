/*
    Package irondb implements an in-memory database of password records.
*/
package irondb


import (
    "encoding/json"
    "github.com/dmulholland/ironclad/ironcrypt"
)


// Length of the key derivation salt in bytes.
const SaltLength = 32


// Number of iterations used by the key derivation function.
const PBKDFIterations = 10000


// DB represents an in-memory database of password records.
type DB struct {
    Salt []byte         `json:"salt"`
    Entries []*Entry    `json:"entries"`
}


// New initializes a new database instance.
func New() (db *DB, err error) {

    // Generate a random salt.
    salt, err := ironcrypt.RandBytes(SaltLength)
    if err != nil {
        return
    }

    // Initialize a new database instance.
    db = &DB{
        Salt: salt,
        Entries: make([]*Entry, 0)}

    return db, err
}


// FromJSON initializes a new database instance from a serialized byte-slice.
func FromJSON(data []byte) (db *DB, err error) {
    db = &DB{}
    err = json.Unmarshal(data, db)
    return db, err
}


// ToJSON serializes a database instance as a byte-slice of JSON.
func (db *DB) ToJSON() (data []byte, err error) {
    data, err = json.Marshal(db)
    if err != nil {
        return nil, err
    }
    return data, nil
}


// Key returns the encryption key corresponding to the specified password.
// This key is used to encrypt passwords *within* the database.
func (db *DB) Key(password string) []byte {
    return ironcrypt.Key([]byte(password), db.Salt, PBKDFIterations)
}


// Import adds entries from a previously-exported byte-slice of JSON.
func (db *DB) Import(key, data []byte) error {

    exports := make([]*ExportEntry, 0)
    err := json.Unmarshal(data, &exports)
    if err != nil {
        return err
    }

    for _, export := range exports {
        entry := NewEntry()
        entry.Title = export.Title
        entry.Url = export.Url
        entry.Username = export.Username
        entry.Email = export.Email
        entry.Tags = export.Tags
        entry.Notes = export.Notes
        err = entry.SetPassword(key, export.Password)
        if err != nil {
            return err
        }
        db.Add(entry)
    }

    return nil
}


// Add inserts a new entry into the database.
func (db *DB) Add(entry *Entry) {
    if len(db.Entries) == 0 {
        entry.Id = 1
    } else {
        entry.Id = db.Entries[len(db.Entries) - 1].Id + 1
    }
    db.Entries = append(db.Entries, entry)
}


// Delete sets an entry's active status to false.
func (db *DB) Delete(id int) {
    for _, entry := range db.Entries {
        if entry.Id == id {
            entry.Active = false
        }
    }
}


// Purge clears deleted entries from the database.
func (db *DB) Purge() {
    entries := db.Entries
    db.Entries = make([]*Entry, 0)
    for _, entry := range entries {
        if entry.Active {
            db.Add(entry)
        }
    }
}


// Size returns the number of active entries in the database.
func (db *DB) Size() int {
    return len(db.Active())
}


// All returns a list containing all the database's entries.
func (db *DB) All() EntryList {
    return db.Entries
}


// Active returns a list containing all the database's active entries.
func (db *DB) Active() EntryList {
    return db.All().FilterActive()
}


// TagMap returns a map of tags to entry-lists.
func (db *DB) TagMap() map[string]EntryList {
    tags := make(map[string]EntryList)
    for _, entry := range db.Active() {
        for _, tag := range entry.Tags {
            if _, ok := tags[tag]; !ok {
                tags[tag] = make([]*Entry, 0)
            }
            tags[tag] = append(tags[tag], entry)
        }
    }
    return tags
}
