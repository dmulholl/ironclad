package main


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
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

    // Check that at least one entry argument has been supplied.
    if !parser.HasArgs() {
        exit("you must specify at least one entry argument")
    }

    // Load the database.
    db, password, filename := loadDB(parser)

    // Grab the entries to delete.
    entries := db.Lookup(parser.GetArgs()...)
    if len(entries) == 0 {
        exit("no matching entries")
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

    // Save the updated database to disk.
    saveDB(db, password, filename)
}
