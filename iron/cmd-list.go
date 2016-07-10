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

  Print a list of entries from a database. Entries to list can be specified by
  ID or by title. (Titles are checked for a case-insensitive substring match.)

  If no arguments are specified, all the entries in the database will be
  listed.

  The 'list' command has an alias, 'show', which is equivalent to:

    list --verbose --cleartext

Arguments:
  [entries]                 Entries to list by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.
  -t, --tag <str>           Filter entries using the specified tag.

Flags:
  -c, --cleartext           Print passwords in cleartext.
      --help                Print this command's help text and exit.
  -v, --verbose             Use the verbose list format.
`, filepath.Base(os.Args[0]))


// Callback for the 'list' command.
func listCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Determine the filename to use.
    filename = parser.GetStr("file")
    if filename == "" {
        if filename, found = fetchLastFilename(); !found {
            filename = input("Filename: ")
        }
    }

    // Determine the password to use.
    password = parser.GetStr("db-password")
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
    cacheLastPassword(password)
    cacheLastFilename(filename)

    // Has the 'show' alias been used?
    if parser.GetParent().GetCmdName() == "show" {
        parser.SetFlag("verbose", true)
        parser.SetFlag("cleartext", true)
    }

    // Assemble a list of entries.
    var entries []*irondb.Entry
    var title string

    if parser.HasArgs() {
        entries = db.Lookup(parser.GetArgs()...)
        title = fmt.Sprintf("%d Matching Entries", len(entries))
    } else {
        entries = db.Active()
        title = "All Entries"
    }

    // Filter by tag.
    if parser.GetStr("tag") != "" {
        entries = irondb.FilterByTag(entries, parser.GetStr("tag"))
        title = fmt.Sprintf("%d Matching Entries", len(entries))
    }

    // Print the list of entries.
    if parser.GetFlag("verbose") {
        clearflag := parser.GetFlag("cleartext")
        printVerboseList(entries, db.Size(), key, title, clearflag)
    } else {
        printCompactList(entries, db.Size())
    }
}


// Print a compact listing.
func printCompactList(entries []*irondb.Entry, dbsize int) {

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
    fmt.Printf("%4d/%d Entries\n", len(entries), dbsize)
    line("-")
}


// Print a verbose listing.
func printVerboseList(
    entries []*irondb.Entry, dbsize int, key []byte, title string, clear bool) {

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

        if entry.Url != "" {
            fmt.Printf("  URL:      %s\n", entry.Url)
        }

        if entry.Username != "" {
            fmt.Printf("  Username: %s\n", entry.Username)
        }

        password, err := entry.GetPassword(key)
        if err != nil {
            exit("Error:", err)
        }

        if clear {
            fmt.Printf("  Password: %s\n", password)
        } else {
            fmt.Printf("  Password: %s\n", stars(len([]rune(password))))
        }

        if entry.Email != "" {
            fmt.Printf("  Email:    %s\n", entry.Email)
        }

        if len(entry.Tags) > 0 {
            fmt.Printf("  Tags:     %s\n", strings.Join(entry.Tags, ", "))
        }

        if entry.Notes != "" {
            iline("~")
            wrapped := wordwrap.WrapString(entry.Notes, 76)
            indented := text.Indent(wrapped, "  ")
            fmt.Println(strings.Trim(indented, "\r\n"))
        }

        line("-")
    }

    // Footer.
    fmt.Printf("  %d/%d Entries\n", len(entries), dbsize)
    line("-")
}
