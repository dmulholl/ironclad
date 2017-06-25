package irondb


// An Entry object represents a single database record.
type Entry struct {
    Id int              `json:"id"`
    Active bool         `json:"active"`
    Title string        `json:"title"`
    Url string          `json:"url"`
    Username string     `json:"username"`
    Passwords []string  `json:"passwords"`
    Email string        `json:"email"`
    Tags []string       `json:"tags"`
    Notes string        `json:"notes"`
}


// An ExportEntry object represents a database record prepared for export.
type ExportEntry struct {
    Title string        `json:"title"`
    Url string          `json:"url"`
    Username string     `json:"username"`
    Passwords []string  `json:"passwords"`
    Email string        `json:"email"`
    Tags []string       `json:"tags"`
    Notes string        `json:"notes"`
}


// NewEntry initializes a new Entry object.
func NewEntry() *Entry {
    return &Entry{
        Active: true,
        Passwords: make([]string, 0),
        Tags: make([]string, 0),
    }
}


// GetPassword returns the newest password from the passwords list.
func (entry *Entry) GetPassword() string {
    if len(entry.Passwords) > 0 {
        return entry.Passwords[len(entry.Passwords) - 1]
    }
    return ""
}


// SetPassword appends a new password to the passsword list.
func (entry *Entry) SetPassword(password string) {
    entry.Passwords = append(entry.Passwords, password)
}
