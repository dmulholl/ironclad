/*
    Package irondb implements an in-memory database of password records.
*/
package irondb


import (
    "encoding/json"
    "strings"
    "strconv"
    "bytes"
    "github.com/dmulholland/ironclad/ironio"
)


// DB represents an in-memory database of password records.
type DB struct {
    entries []*Entry
}


// New returns a new database.
func New() *DB {
    return &DB{ entries: make([]*Entry, 0) }
}


// Load loads a saved database from an encrypted file.
func Load(password, filename string) (db *DB, key []byte, err error) {

    // Load the JSON data store from the encrypted file.
    data, key, err := ironio.Load(password, filename)
    if err != nil {
        return db, key, err
    }

    // Unmarshal the stored JSON.
    db = New()
    err = json.Unmarshal(data, &db.entries)

    return db, key, err
}


// Save saves a database to an encrypted file.
func (db *DB) Save(key []byte, password, filename string) error {

    // Generate a JSON dump of the database.
    data, err := json.Marshal(db.entries)
    if err != nil {
        return err
    }

    // Save the JSON dump as an encrypted file.
    return ironio.Save(data, key, password, filename)
}


// Export exports a list of entries in JSON format. Passwords are unencrypted.
func (db *DB) Export(key []byte, queries ...string) (dump string, err error) {

    var entries []*Entry

    // If no query strings have been specified, export all active entries.
    if len(queries) == 0 {
        entries = db.Active()
    } else {
        entries = db.Lookup(queries...)
    }

    // Create a list of entries with unencrypted passwords.
    clones := make([]Entry, 0)
    for _, original := range entries {
        clone := *original
        clone.Password, err = original.GetPassword(key)
        if err != nil {
            return "", err
        }
        clones = append(clones, clone)
    }

    // Generate a JSON dump of the list.
    data, err := json.Marshal(clones)
    if err != nil {
        return "", err
    }

    // Format the JSON for display.
    var formatted bytes.Buffer
    json.Indent(&formatted, data, "", "  ")

    return formatted.String(), nil
}


// Import adds entries from an exported JSON dump to the database.
func (db *DB) Import(key []byte, dump string) error {

    entries := make([]*Entry, 0)
    err := json.Unmarshal([]byte(dump), &entries)
    if err != nil {
        return err
    }

    for _, entry := range entries {
        err = entry.SetPassword(key, entry.Password)
        if err != nil {
            return err
        }
        db.Add(entry)
    }

    return nil
}


// Active returns a list of active entries.
func (db *DB) Active() []*Entry {
    entries := make([]*Entry, 0)
    for _, entry := range db.entries {
        if entry.Active {
            entries = append(entries, entry)
        }
    }
    return entries
}


// ByTag returns a list of active entries associated with the specified tag.
func (db *DB) ByTag(tag string) []*Entry {
    entries := make([]*Entry, 0)
    for _, entry := range db.entries {
        if entry.Active {
            for _, t := range entry.Tags {
                if strings.ToLower(t) == strings.ToLower(tag) {
                    entries = append(entries, entry)
                }
            }
        }
    }
    return entries
}


// Add inserts a new entry into the database.
func (db *DB) Add(entry *Entry) {
    if len(db.entries) == 0 {
        entry.Id = 1
    } else {
        entry.Id = db.entries[len(db.entries) - 1].Id + 1
    }
    db.entries = append(db.entries, entry)
}


// Delete removes an entry from the database.
func (db *DB) Delete(id int) {
    for _, entry := range db.entries {
        if entry.Id == id {
            entry.Active = false
        }
    }
}


// Purge clears deleted entries from the database.
func (db *DB) Purge() {
    entries := db.entries
    db.entries = make([]*Entry, 0)
    for _, entry := range entries {
        if entry.Active {
            db.Add(entry)
        }
    }
}


// Lookup searches the database for entries matching the query string or
// strings. A query string can be an ID or a case-insensitive title or
// title-prefix.
func (db *DB) Lookup(queries ...string) []*Entry {

    // List of entries to return.
    matches := make([]*Entry, 0)

    // We only want to look for active entries.
    active := db.Active()

    for _, query := range queries {

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

        // Check for a case-insensitive prefix match on the entry title.
        query = strings.ToLower(query)
        for _, entry := range active {
            if strings.HasPrefix(strings.ToLower(entry.Title), query) {
                matches = append(matches, entry)
            }
        }
    }

    return matches
}
