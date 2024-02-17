package database

import (
	"encoding/json"
	"fmt"
)

// DB represents an in-memory database of password records.
type DB struct {
	Version   int      `json:"version"`
	CachePass string   `json:"cachepass"`
	Entries   []*Entry `json:"entries"`
}

// New initializes a new database instance.
func New(cachepass string) *DB {
	return &DB{
		Version:   2,
		CachePass: cachepass,
		Entries:   make([]*Entry, 0),
	}
}

// FromJSON initializes a database instance from JSON.
func FromJSON(input []byte) (*DB, error) {
	db := &DB{}
	err := json.Unmarshal(input, db)
	return db, err
}

// ToJSON serializes a database instance to JSON.
func (db *DB) ToJSON() ([]byte, error) {
	return json.Marshal(db)
}

// Import adds entries to the database from JSON. The entries must have been exported using an
// EntryList's Export() method. Returns the number of new entrie.
func (db *DB) Import(data []byte) (int, error) {
	entries := []*ExportEntry{}

	err := json.Unmarshal(data, &entries)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshall input: %w", err)
	}

	for _, entry := range entries {
		newEntry := NewEntry()
		newEntry.Title = entry.Title
		newEntry.Url = entry.Url
		newEntry.Username = entry.Username
		newEntry.Passwords = entry.Passwords
		newEntry.Email = entry.Email
		newEntry.Tags = entry.Tags
		newEntry.Notes = entry.Notes
		db.Add(newEntry)
	}

	return len(entries), nil
}

// Add inserts a new entry into the database.
func (db *DB) Add(entry *Entry) {
	if len(db.Entries) == 0 {
		entry.Id = 1
	} else {
		entry.Id = db.Entries[len(db.Entries)-1].Id + 1
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

// Count returns the number of active entries in the database.
func (db *DB) Count() int {
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
	tagMap := map[string]EntryList{}

	for _, entry := range db.Active() {
		for _, tag := range entry.Tags {
			tagMap[tag] = append(tagMap[tag], entry)
		}
	}

	return tagMap
}
