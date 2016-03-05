package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "strings"
    "github.com/dmulholland/ironclad/irondb"
)


// Help text for the 'edit' command.
var editHelptext = fmt.Sprintf(`
Usage: %s edit [FLAGS] [OPTIONS] ARGUMENTS

  Edit an existing database entry.

Arguments:
  <entry>                   Entry to edit by ID or title.

Options:
  -f, --file <str>          Database file.

Flags:
      --help                Print this command's help text and exit.
      --notes               Edit the entry's notes.
      --password            Edit the entry's password.
      --tags                Edit the entry's tags.
      --title               Edit the entry's title.
      --url                 Edit the entry's url.
      --username            Edit the entry's username.
`, filepath.Base(os.Args[0]))


// Callback for the 'edit' command.
func editCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("Error: missing entry argument.")
    }

    // Determine the filename to use.
    filename = parser.GetStrOpt("file")
    if filename == "" {
        if filename, found = fetchLastFilename(); !found {
            filename = input("Filename: ")
        }
    }

    // Determine the password to use.
    password = parser.GetStrOpt("db-password")
    if password == "" {
        if password, found = fetchLastPassword(); !found {
            password = input("Password: ")
        }
    }

    // Load the database.
    db, key, err := irondb.Load(password, filename)
    if err != nil {
        exit("Error:", err)
    }

    // Search for an entry corresponding to the specified argument.
    entries := db.Lookup(parser.GetArgs()[0])
    if len(entries) == 0 {
        exit("Error: no matching entry.")
    } else if len(entries) > 1 {
        exit("Error: query matches multiple entries.")
    }
    entry := entries[0]

    // Check that we have at least one field to edit.
    if !parser.GetFlag("title") && !parser.GetFlag("url") &&
        !parser.GetFlag("username") && !parser.GetFlag("password") &&
        !parser.GetFlag("tags") && !parser.GetFlag("notes") {
        exit("Error: you must specify at least one field to edit.")
    }

    // Header.
    line("-")
    fmt.Println("  Editing Entry: " + entry.Title)
    line("-")

    if parser.GetFlag("title") {
        fmt.Println("  TITLE")
        fmt.Println("  Old value: " + entry.Title)
        entry.Title = input("  New value: ")
        line("-")
    }

    if parser.GetFlag("url") {
        fmt.Println("  URL")
        fmt.Println("  Old value: " + entry.Url)
        entry.Url = input("  New value: ")
        line("-")
    }

    if parser.GetFlag("username") {
        fmt.Println("  USERNAME")
        fmt.Println("  Old value: " + entry.Username)
        entry.Username = input("  New value: ")
        line("-")
    }

    if parser.GetFlag("password") {
        fmt.Println("  PASSWORD")
        oldpass, err := entry.GetPassword(key)
        if err != nil {
            exit("Error: ", err)
        }
        fmt.Println("  Old value: " + oldpass)
        err = entry.SetPassword(key, input("  New value: "))
        if err != nil {
            exit("Error: ", err)
        }
        line("-")
    }

    if parser.GetFlag("tags") {
        fmt.Println("  TAGS")
        fmt.Println("  Old value: " + strings.Join(entry.Tags, ", "))
        tagstring := input("  New value: ")
        tagslice := strings.Split(tagstring, ",")
        entry.Tags = make([]string, 0)
        for _, tag := range tagslice {
            tag = strings.TrimSpace(tag)
            if tag != "" {
                entry.Tags = append(entry.Tags, tag)
            }
        }
        line("-")
    }

    if parser.GetFlag("notes") {
        entry.Notes = inputViaEditor("edit-note", entry.Notes)
    }

    // Save the updated database to disk.
    db.Save(key, password, filename)

    // Footer.
    fmt.Println("  Entry updated.")

    // Cache the password and filename.
    cacheLastPassword(password)
    cacheLastFilename(filename)
}
