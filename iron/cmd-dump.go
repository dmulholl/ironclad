package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "encoding/json"
    "bytes"
    "github.com/dmulholland/ironclad/ironio"
)


// Help text for the 'dump' command.
var dumpHelptext = fmt.Sprintf(`
Usage: %s dump [FLAGS] [OPTIONS]

  Dump a database's internal JSON data store to stdout.

Options:
  -f, --file <str>          Database file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'dump' command.
func dumpCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

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

    // Load the JSON data store from the encrypted database file.
    data, _, err := ironio.Load(password, filename)
    if err != nil {
        exit("Error:", err)
    }
    cacheLastPassword(password)
    cacheLastFilename(filename)

    // Format the JSON for display.
    var formatted bytes.Buffer
    json.Indent(&formatted, data, "", "  ")

    // Print the formatted JSON to stdout.
    fmt.Println(formatted.String())
}
