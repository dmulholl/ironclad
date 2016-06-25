/*
    Package irondb implements an in-memory database of password records.
*/
package irondb


import (
    "encoding/json"
    "strings"
    "strconv"
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
        entry.Active = true
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


// FilterByTag returns a list of active entries matching the specified tag.
func (db *DB) FilterByTag(tag string) []*Entry {
    return FilterByTag(db.Active(), tag)
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


// Lookup searches the database for entries matching the specified query
// strings. A query string can be an entry ID or a case-insensitive substring
// of an entry title.
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


// Size returns the number of active entries in the database.
func (db *DB) Size() int {
    return len(db.Active())
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
