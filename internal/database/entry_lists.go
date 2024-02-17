package database

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// An EntryList is a slice of Entry pointers.
type EntryList []*Entry

// Contains returns true if the EntryList contains the specified entry.
func (list EntryList) Contains(entry *Entry) bool {
	for _, current := range list {
		if entry == current {
			return true
		}
	}
	return false
}

// FilterActive filters an EntryList returning only those entries which are
// active.
func (list EntryList) FilterActive() EntryList {
	active := make([]*Entry, 0)
	for _, entry := range list {
		if entry.Active {
			active = append(active, entry)
		}
	}
	return active
}

// FilterInactive filters an EntryList returning only those entries which are
// inactive.
func (list EntryList) FilterInactive() EntryList {
	inactive := make([]*Entry, 0)
	for _, entry := range list {
		if !entry.Active {
			inactive = append(inactive, entry)
		}
	}
	return inactive
}

// FilterByTag filters an EntryList returning only those entries which match
// the specified tag. Matches are case-insensitive.
func (list EntryList) FilterByTag(tag string) EntryList {
	matches := make([]*Entry, 0)
	target := strings.ToLower(tag)
	for _, entry := range list {
		for _, entrytag := range entry.Tags {
			if strings.ToLower(entrytag) == target {
				matches = append(matches, entry)
			}
		}
	}
	return matches
}

// FilterByAny filters an EntryList returning entries which match *any* of the
// specified query strings. Each query string can be an entry ID or a
// case-insensitive substring of an entry title.
func (list EntryList) FilterByAny(queries ...string) EntryList {
	matches := EntryList(make([]*Entry, 0))
	for _, query := range queries {

		// First, see if we can parse the query string as an integer ID.
		// If we can, look for an entry with a matching ID.
		if i, err := strconv.ParseInt(query, 10, 32); err == nil {
			id := int(i)
			for _, entry := range list {
				if id == entry.Id && !matches.Contains(entry) {
					matches = append(matches, entry)
					break
				}
			}
		}

		// Check for a case-insensitive substring match on the entry title.
		query = strings.ToLower(query)
		for _, entry := range list {
			if strings.Contains(strings.ToLower(entry.Title), query) {
				if !matches.Contains(entry) {
					matches = append(matches, entry)
				}
			}
		}
	}

	return matches
}

// FilterByAll filters an EntryList returning entries which match *all* of the
// specified query strings, where each query string is a case-insensitive
// substring of the entry title. If a single query string is supplied it will
// first be checked to see if it matches a valid entry ID.
func (list EntryList) FilterByAll(queries ...string) EntryList {
	matches := make([]*Entry, 0)

	// If we have a single query string, check if it matches an entry ID.
	if len(queries) == 1 {
		if i, err := strconv.ParseInt(queries[0], 10, 32); err == nil {
			id := int(i)
			for _, entry := range list {
				if id == entry.Id {
					matches = append(matches, entry)
					return matches
				}
			}
		}
	}

	// Convert all query strings to lower case.
	for i, _ := range queries {
		queries[i] = strings.ToLower(queries[i])
	}

	// Check entry titles for a case-insensitive substring match on all
	// query strings.
	for _, entry := range list {
		entrytitle := strings.ToLower(entry.Title)
		match := true
		for _, query := range queries {
			if !strings.Contains(entrytitle, query) {
				match = false
				break
			}
		}
		if match {
			matches = append(matches, entry)
		}
	}

	return matches
}

// FilterByIDString filters an EntryList returning only those entries which
// match the specified query strings where each query string is an entry ID.
func (list EntryList) FilterByIDString(queries ...string) EntryList {
	matches := make([]*Entry, 0)
	for _, query := range queries {
		if i, err := strconv.ParseInt(query, 10, 32); err == nil {
			id := int(i)
			for _, entry := range list {
				if id == entry.Id {
					matches = append(matches, entry)
					break
				}
			}
		}
	}
	return matches
}

// Export exports a list of entries as a JSON string.
func (list EntryList) Export() (string, error) {

	// Assemble a list of ExportEntry objects.
	exports := make([]ExportEntry, 0)
	for _, entry := range list {
		export := ExportEntry{
			Title:     entry.Title,
			Url:       entry.Url,
			Username:  entry.Username,
			Passwords: entry.Passwords,
			Email:     entry.Email,
			Notes:     entry.Notes,
			Tags:      entry.Tags,
		}
		exports = append(exports, export)
	}

	// Generate a JSON dump of the list.
	data, err := json.Marshal(exports)
	if err != nil {
		return "", err
	}

	// Format the JSON for display.
	var formatted bytes.Buffer
	json.Indent(&formatted, data, "", "  ")

	return formatted.String(), nil
}
