package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


var purgeHelp = fmt.Sprintf(`
Usage: %s purge [FLAGS] [OPTIONS]

  Purge deleted entries from a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func purgeCallback(parser *clio.ArgParser) {
    filename, password, db := loadDB(parser)
    db.Purge()
    saveDB(filename, password, db)
}
