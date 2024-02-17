package database

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// FilterActive returns a new EntryList containing active entries.
func (list EntryList) FilterActive() EntryList {
	active := EntryList{}

	for _, entry := range list {
		if entry.Active {
			active = append(active, entry)
		}
	}

	return active
}

// FilterInactive returns a new EntryList containing inactive entries.
func (list EntryList) FilterInactive() EntryList {
	inactive := EntryList{}

	for _, entry := range list {
		if !entry.Active {
			inactive = append(inactive, entry)
		}
	}

	return inactive
}

// FilterByTag returns a new EntryList containing entries which match the specified tag. Matches
// are case-insensitive.
func (list EntryList) FilterByTag(tag string) EntryList {
	tag = strings.ToLower(tag)
	matches := EntryList{}

	for _, entry := range list {
		for _, entrytag := range entry.Tags {
			if strings.ToLower(entrytag) == tag {
				matches = append(matches, entry)
			}
		}
	}

	return matches
}

// FilterByAny returns a new EntryList containing entries which match *any* of the specified query
// strings. A query string can be an entry ID or a case-insensitive substring of an entry title.
func (list EntryList) FilterByAny(queries ...string) EntryList {
	matches := EntryList{}

	for _, query := range queries {
		if i, err := strconv.ParseInt(query, 10, 32); err == nil {
			id := int(i)
			for _, entry := range list {
				if id == entry.Id && !matches.Contains(entry) {
					matches = append(matches, entry)
					break
				}
			}
		}

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

// FilterByAll returns a new EntryList containing entries which match *all* of the specified query
// strings, where each query string is a case-insensitive substring of the entry title. If a single
// query string is supplied, it will first be checked to see if it matches a valid entry ID.
func (list EntryList) FilterByAll(queries ...string) EntryList {
	matches := EntryList{}

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
	for i := range queries {
		queries[i] = strings.ToLower(queries[i])
	}

	// Check entry titles for a case-insensitive substring match on all query strings.
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

// FilterByID returns a new EntryList containing entries with matching IDs.
func (list EntryList) FilterByID(ids ...int) EntryList {
	matches := EntryList{}

	for _, id := range ids {
		for _, entry := range list {
			if id == entry.Id {
				matches = append(matches, entry)
				break
			}
		}
	}

	return matches
}

// Export exports the EntryList as formatted JSON.
func (list EntryList) Export() (string, error) {
	exports := []ExportEntry{}

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

	data, err := json.Marshal(exports)
	if err != nil {
		return "", fmt.Errorf("failed to marshall entries as JSON: %w", err)
	}

	var formatted bytes.Buffer
	json.Indent(&formatted, data, "", "  ")

	return formatted.String(), nil
}
