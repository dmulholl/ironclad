package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/irondb"
)


// Help text for the 'export' command.
var exportHelptext = fmt.Sprintf(`
Usage: %s export [FLAGS] [OPTIONS] [ARGUMENTS]

  Export a list of entries in JSON format.

Arguments:
  [entry ...]               List of entries to export by ID or title.

Options:
  -f, --file <str>          Database file.

Flags:
  --help                    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'export' command.
func exportCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Determine the filename to use.
    filename = parser.GetStringOption("file")
    if filename == "" {
        if filename, found = fetchLastFilename(); !found {
            filename = input("Filepath: ")
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

    // Assemble a list of entries to export.
    var entries []*irondb.Entry

    if parser.HasArgs() {
        entries = db.Lookup(parser.GetArgs()...)
    } else {
        entries = db.Active()
    }

    // Create the JSON dump.
    dump, err := irondb.Export(entries, key)
    if err != nil {
        exit("Error:", err)
    }

    // Print the JSON to stdout.
    fmt.Println(dump)

    // Cache the password and filename.
    cacheLastPassword(password)
    cacheLastFilename(filename)
}
