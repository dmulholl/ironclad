package irondb


import (
    "encoding/json"
)


// DB represents an in-memory database of password records.
type DB struct {
    Entries []*Entry    `json:"entries"`
}


// New initializes a new database instance.
func New() (db *DB) {
    return &DB{Entries: make([]*Entry, 0)}
}


// FromJSON initializes a database instance from a serialized byte-slice.
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


// Import adds entries from a previously-exported byte-slice of JSON.
func (db *DB) Import(data []byte) error {

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
        entry.Passwords = export.Passwords
        entry.Email = export.Email
        entry.Tags = export.Tags
        entry.Notes = export.Notes
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


// SetActive sets an entry's active status to true.
func (db *DB) SetActive(id int) {
    for _, entry := range db.Entries {
        if entry.Id == id {
            entry.Active = true
        }
    }
}


// SetInactive sets an entry's active status to false.
func (db *DB) SetInactive(id int) {
    for _, entry := range db.Entries {
        if entry.Id == id {
            entry.Active = false
        }
    }
}


// Purge clears inactive entries from the database.
func (db *DB) PurgeInactive() {
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


// Inactive returns a list containing all the database's inactive entries.
func (db *DB) Inactive() EntryList {
    return db.All().FilterInactive()
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
