package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


// Help text for the 'purge' command.
var purgeHelp = fmt.Sprintf(`
Usage: %s purge [FLAGS] [OPTIONS]

  Purge deleted entries from a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'purge' command.
func purgeCallback(parser *clio.ArgParser) {

    // Load the database.
    filename, password, db := loadDB(parser)

    // Purge the database.
    db.Purge()

    // Save the updated database to disk.
    saveDB(filename, password, db)
}
