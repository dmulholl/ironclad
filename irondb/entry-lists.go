package irondb


import (
    "strings"
    "strconv"
    "encoding/json"
    "bytes"
)


// An EntryList is a slice of Entry pointers.
type EntryList []*Entry


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


// FilterByQuery filters an EntryList returning only those entries which match
// the specified query strings. Each query string can be an entry ID or a
// case-insensitive substring of an entry title.
func (list EntryList) FilterByQuery(queries ...string) EntryList {
    matches := make([]*Entry, 0)
    for _, query := range queries {

        // First, see if we can parse the query string as an integer ID.
        // If we can, look for an entry with a matching ID.
        if i, err := strconv.ParseInt(query, 10, 32); err == nil {
            id := int(i)
            for _, entry := range list {
                if id == entry.Id {
                    matches = append(matches, entry)
                    break
                }
            }
        }

        // Check for a case-insensitive substring match on the entry title.
        query = strings.ToLower(query)
        for _, entry := range list {
            if strings.Contains(strings.ToLower(entry.Title), query) {
                matches = append(matches, entry)
            }
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


// FilterProgressive filters an EntryList using progressively less restrictive
// criteria until it finds one or more matches or runs out of criteria. The
// query string may be (in order) an entry ID or a case-insensitive exact,
// prefix, or substring match for an entry title.
func (list EntryList) FilterProgressive(query string) EntryList {
    matches := make([]*Entry, 0)

    // First, see if we can parse the query string as an integer ID.
    // If we can, look for an entry with a matching ID.
    if i, err := strconv.ParseInt(query, 10, 32); err == nil {
        id := int(i)
        for _, entry := range list {
            if id == entry.Id {
                matches = append(matches, entry)
                break
            }
        }
    }
    if len(matches) > 0 {
        return matches
    }

    // No ID match so check for an exact match on an entry title.
    query = strings.ToLower(query)
    for _, entry := range list {
        if query == strings.ToLower(entry.Title) {
            matches = append(matches, entry)
        }
    }
    if len(matches) > 0 {
        return matches
    }

    // No exact match so check for a prefix match on an entry title.
    for _, entry := range list {
        if strings.HasPrefix(strings.ToLower(entry.Title), query) {
            matches = append(matches, entry)
        }
    }
    if len(matches) > 0 {
        return matches
    }

    // No prefix match so check for a substring match.
    for _, entry := range list {
        if strings.Contains(strings.ToLower(entry.Title), query) {
            matches = append(matches, entry)
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
            Title: entry.Title,
            Url: entry.Url,
            Username: entry.Username,
            Password: entry.Password,
            Email: entry.Email,
            Notes: entry.Notes,
            Tags: entry.Tags,
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
