/*
    Package irondb implements an in-memory database of password records.
*/
package irondb


import (
    "encoding/json"
    "strings"
    "strconv"
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


// New initializes and returns an empty database instance.
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


// FromJSON initializes a new database instance from a byte-slice of JSON.
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


// Import adds entries from an exported byte-slice of JSON.
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


// Delete removes an entry from the database.
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









// Active returns a list of active entries.
func (db *DB) Active() []*Entry {
    entries := make([]*Entry, 0)
    for _, entry := range db.Entries {
        if entry.Active {
            entries = append(entries, entry)
        }
    }
    return entries
}



// TagMap returns a map of tags to entry-lists.
func (db *DB) TagMap() map[string][]*Entry {
    tags := make(map[string][]*Entry)
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


// Lookup searches the database for entries matching the specified query strings
// Each query string can be an entry ID or a case-insensitive substring of an
// entry title.
func (db *DB) Lookup(queries ...string) []*Entry {

    // List of entries to return.
    matches := make([]*Entry, 0)

    // We only want to look for active entries.
    active := db.Active()

    for _, query := range queries {

        // String comparisons will be case-insensitive.
        query = strings.ToLower(query)

        // First, see if we can parse the query string as an integer ID.
        if i, err := strconv.ParseInt(query, 10, 32); err == nil {
            id := int(i)
            for _, entry := range active {
                if id == entry.Id {
                    matches = append(matches, entry)
                    break
                }
            }
        }

        // Check for a case-insensitive substring match on the entry title.
        for _, entry := range active {
            if strings.Contains(strings.ToLower(entry.Title), query) {
                matches = append(matches, entry)
            }
        }
    }

    return matches
}


// LookupUnique searches the database for a single entry matching the query
// string. The query string may be (in order) an entry ID or a
// (case-insensitive) exact, prefix, or substring match for an entry title.
// This function returns a slice of entries; zero or multiple matches may be
// interpreted by the caller as error conditions.
func (db *DB) LookupUnique(query string) []*Entry {

    // List of entries to return.
    matches := make([]*Entry, 0)

    // We only want to look for active entries.
    active := db.Active()

    // String comparisons will be case-insensitive.
    query = strings.ToLower(query)

    // First, see if we can parse the query string as an integer ID.
    if i, err := strconv.ParseInt(query, 10, 32); err == nil {
        id := int(i)
        for _, entry := range active {
            if id == entry.Id {
                matches = append(matches, entry)
                return matches
            }
        }
    }

    // Check for an exact match on the entry title.
    for _, entry := range active {
        if query == strings.ToLower(entry.Title) {
            matches = append(matches, entry)
        }
    }
    if len(matches) > 0 {
        return matches
    }

    // No exact match so check for a prefix match on the entry title.
    for _, entry := range active {
        if strings.HasPrefix(strings.ToLower(entry.Title), query) {
            matches = append(matches, entry)
        }
    }
    if len(matches) > 0 {
        return matches
    }

    // No exact or prefix match so check for a substring match.
    for _, entry := range active {
        if strings.Contains(strings.ToLower(entry.Title), query) {
            matches = append(matches, entry)
        }
    }

    return matches
}


// FilterByTag filters a list of entries by the specified tag.
func FilterByTag(entries []*Entry, tag string) []*Entry {
    matches := make([]*Entry, 0)
    searchtag := strings.ToLower(tag)
    for _, entry := range entries {
        for _, entrytag := range entry.Tags {
            if strings.ToLower(entrytag) == searchtag {
                matches = append(matches, entry)
            }
        }
    }
    return matches
}
