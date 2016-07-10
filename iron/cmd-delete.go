package main


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/irondb"
)


// Help text for the 'delete' command.
var deleteHelptext = fmt.Sprintf(`
Usage: %s delete [FLAGS] [OPTIONS] ARGUMENTS

  Delete one or more entries from a database.

Arguments:
  <entries>                 List of entries to delete by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'delete' command.
func deleteCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Check that at least one entry argument has been supplied.
    if !parser.HasArgs() {
        exit("Error: you must specify at least one entry argument.")
    }

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

    // Grab the entries to delete.
    entries := db.Lookup(parser.GetArgs()...)
    if len(entries) == 0 {
        exit("Error: no matching entries.")
    }

    // Print a listing and request confirmation.
    printCompactList(entries, db.Size())
    confirm := input("  Delete the entries listed above? (y/n): ")
    if strings.ToLower(confirm)[0] == 'y' {
        for _, entry := range entries {
            db.Delete(entry.Id)
        }
        line("-")
        fmt.Println("  Entries deleted.")
    } else {
        line("-")
        fmt.Println("  Deletion aborted.")
    }

    // Save the altered database.
    db.Save(key, password, filename)
}
