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

  Export a list of entries in JSON format. Entries can be specified by ID or
  by title. If no entries are specified, all entries will be exported.

Arguments:
  [entries]                 List of entries to export by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'export' command.
func exportCallback(parser *clio.ArgParser) {

    // Load the database.
    db, password, _ := loadDB(parser)

    // Assemble a list of entries to export.
    var entries []*irondb.Entry
    if parser.HasArgs() {
        entries = db.Lookup(parser.GetArgs()...)
    } else {
        entries = db.Active()
    }

    // Create the JSON dump.
    dump, err := irondb.Export(entries, db.Key(password))
    if err != nil {
        exit(err)
    }

    // Print the JSON to stdout.
    fmt.Println(dump)
}
