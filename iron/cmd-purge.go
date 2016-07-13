package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
)


// Help text for the 'purge' command.
var purgeHelptext = fmt.Sprintf(`
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
    password, filename, db := loadDB(parser)

    // Purge the database.
    db.Purge()

    // Save the updated database to disk.
    saveDB(password, filename, db)
}
