package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "encoding/json"
    "bytes"
)


// Help text for the 'dump' command.
var dumpHelptext = fmt.Sprintf(`
Usage: %s dump [FLAGS] [OPTIONS]

  Dump a database's internal JSON data store to stdout.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'dump' command.
func dumpCallback(parser *clio.ArgParser) {

    // Load the database.
    _, _, db := loadDB(parser)

    // Serialize the database as a byte-slice of JSON.
    data, err := db.ToJSON()
    if err != nil {
        exit(err)
    }

    // Format the JSON for display.
    var formatted bytes.Buffer
    json.Indent(&formatted, data, "", "  ")

    // Print the formatted JSON to stdout.
    fmt.Println(formatted.String())
}
