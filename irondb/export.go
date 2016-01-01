package irondb


import (
    "encoding/json"
    "bytes"
)


// An ExportEntry object represents a database record prepared for export.
type ExportEntry struct {
    Title string        `json:"title"`
    Url string          `json:"url"`
    Username string     `json:"username"`
    Password string     `json:"password"`
    Tags []string       `json:"tags"`
    Notes string        `json:"notes"`
}


// Export exports a list of entries in JSON format. Passwords are unencrypted.
func Export(entries []*Entry, key []byte) (dump string, err error) {

    // Assemble a list of ExportEntry objects with unencrypted passwords.
    exports := make([]ExportEntry, 0)
    for _, entry := range entries {
        export := ExportEntry{
            Title: entry.Title,
            Url: entry.Url,
            Username: entry.Username,
            Notes: entry.Notes,
            Tags: entry.Tags,
        }
        export.Password, err = entry.GetPassword(key)
        if err != nil {
            return "", err
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
