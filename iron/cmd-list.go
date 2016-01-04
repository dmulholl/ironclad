package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/irondb"
    "strings"
    "github.com/tonnerre/golang-text"
    "github.com/mitchellh/go-wordwrap"
)


// Help text for the 'list' command.
var listHelptext = fmt.Sprintf(`
Usage: %s list [FLAGS] [OPTIONS] [ARGUMENTS]

  List the entries in a database.

Arguments:
  [entry ...]               Entries to list by ID or title. Default: all.

Options:
  -f, --file <str>          Database file.
  -t, --tag <str>           List entries by tag.

Flags:
      --help                Print this command's help text and exit.
  -v, --verbose             Use the verbose list format.
`, filepath.Base(os.Args[0]))


// Callback for the 'list' command.
func listCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Determine the filename to use.
    filename = parser.GetStringOption("file")
    if filename == "" {
        if filename, found = fetchLastFilename(); !found {
            filename = input("Filename: ")
        }
    }

    // Determine the password to use.
    password = parser.GetStringOption("db-password")
    if password == "" {
        if password, found = fetchLastPassword(); !found {
            password = input("Password: ")
        }
    }

    // Load the database file.
    db, key, err := irondb.Load(password, filename)
    if err != nil {
        exit("Error:", err)
    }

    // Assemble a list of entries.
    var entries []*irondb.Entry
    var title string

    if parser.HasArgs() {
        entries = db.Lookup(parser.GetArgs()...)
        title = "Matching Entries"
    } else if parser.GetStringOption("tag") != "" {
        entries = db.ByTag(parser.GetStringOption("tag"))
        title = "Entries Tagged: " + parser.GetStringOption("tag")
    } else {
        entries = db.Active()
        title = "All Entries"
    }

    // Print the list of entries.
    if parser.GetFlag("verbose") {
        printVerboseList(entries, key, title)
    } else {
        printCompactList(entries)
    }

    // Cache the password and filename.
    cacheLastPassword(password)
    cacheLastFilename(filename)
}


// Print a compact listing.
func printCompactList(entries []*irondb.Entry) {

    // Bail if we have no entries to display.
    if len(entries) == 0 {
        line("-")
        fmt.Println("  No Entries")
        line("-")
        return
    }

    // Header.
    line("-")
    fmt.Println("  ID  |  TITLE")
    line("-")

    // Print the entry listing.
    for _, entry := range entries {
        fmt.Printf("%4d  |  %s\n", entry.Id, entry.Title)
    }

    // Footer.
    line("-")
    if len(entries) == 1 {
        fmt.Println("  1 Entry")
    } else {
        fmt.Printf("  %d Entries\n", len(entries))
    }
    line("-")
}


// Print a verbose listing.
func printVerboseList(entries []*irondb.Entry, key []byte, title string) {

    // Bail if we have no entries to display.
    if len(entries) == 0 {
        line("-")
        fmt.Println("  No Entries")
        line("-")
        return
    }

    // Header.
    line("-")
    fmt.Println("  " + title)
    line("-")

    // Print the entry listing.
    for _, entry := range entries {
        fmt.Printf("  ID:       %d\n", entry.Id)
        fmt.Printf("  Title:    %s\n", entry.Title)
        fmt.Printf("  URL:      %s\n", entry.Url)
        fmt.Printf("  Username: %s\n", entry.Username)

        password, err := entry.GetPassword(key)
        if err != nil {
            exit("Error:")
        }

        fmt.Printf("  Password: %s\n", password)
        fmt.Printf("  Tags:     %s\n", strings.Join(entry.Tags, ", "))

        if entry.Notes != "" {
            iline("-")
            wrapped := wordwrap.WrapString(entry.Notes, 76)
            indented := text.Indent(wrapped, "  ")
            fmt.Println(strings.Trim(indented, "\r\n"))
        }

        line("-")
    }

    // Footer.
    if len(entries) == 1 {
        fmt.Println("  1 Entry")
    } else {
        fmt.Printf("  %d Entries\n", len(entries))
    }
    line("-")
}
