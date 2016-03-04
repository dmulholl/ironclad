package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "strings"
    "github.com/dmulholland/ironclad/irondb"
)


// Help text for the 'add' command.
var addHelptext = fmt.Sprintf(`
Usage: %s add [FLAGS] [OPTIONS]

  Add a new entry to a database.

Options:
  -f, --file <str>          Database file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'add' command.
func addCallback(parser *clio.ArgParser) {

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

    // Load the database.
    db, key, err := irondb.Load(password, filename)
    if err != nil {
        exit("Error:", err)
    }

    // Create a new Entry object to add to the database.
    entry := irondb.NewEntry()

    // Print header.
    line("-")
    fmt.Println("  Add Entry")
    line("-")

    // Grab user input.
    entry.Title    = input("  Title:      ")
    entry.Url      = input("  URL:        ")
    entry.Username = input("  Username:   ")

    // Get the password.
    entry.SetPassword(key, input("  Password:   "))

    // Split tags on commas.
    line("-")
    tagstring := input("  Enter a comma-separated list of tags for this entry:\n> ")
    tagslice := strings.Split(tagstring, ",")
    for _, tag := range tagslice {
        tag = strings.TrimSpace(tag)
        if tag != "" {
            entry.Tags = append(entry.Tags, tag)
        }
    }

    // Do we need to launch a text editor to add notes?
    line("-")
    notesquery := input("  Add a note to this entry? (y/n): ")
    if len(notesquery) > 0 && strings.ToLower(notesquery)[0] == 'y' {
        entry.Notes = inputViaEditor("add-note", "")
    } else {
        entry.Notes = ""
    }

    // Add the new entry to the database.
    db.Add(entry)

    // Save the updated database to disk.
    db.Save(key, password, filename)

    // Footer.
    line("-")

    // Cache the password and filename.
    cacheLastPassword(password)
    cacheLastFilename(filename)
}
