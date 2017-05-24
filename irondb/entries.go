package irondb


// An Entry object represents a single database record.
type Entry struct {
    Id int              `json:"id"`
    Active bool         `json:"active"`
    Title string        `json:"title"`
    Url string          `json:"url"`
    Username string     `json:"username"`
    Password string     `json:"password"`
    Email string        `json:"email"`
    Tags []string       `json:"tags"`
    Notes string        `json:"notes"`
}


// An ExportEntry object represents a database record prepared for export.
type ExportEntry struct {
    Title string        `json:"title"`
    Url string          `json:"url"`
    Username string     `json:"username"`
    Password string     `json:"password"`
    Email string        `json:"email"`
    Tags []string       `json:"tags"`
    Notes string        `json:"notes"`
}


// NewEntry initializes a new Entry object.
func NewEntry() *Entry {
    return &Entry{
        Active: true,
        Tags: make([]string, 0),
    }
}
