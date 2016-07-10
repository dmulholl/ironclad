package main


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/irondb"
)


// Help text for the 'import' command.
var importHelptext = fmt.Sprintf(`
Usage: %s import [FLAGS] [OPTIONS] ARGUMENTS

  Import a list of entries in JSON format.

Arguments:
  <file>                    File to import.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'import' command.
func importCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("Error: you must supply the name of a file to import.")
    }

    // Determine the filename to use.
    filename = parser.GetStr("file")
    if filename == "" {
        if filename, found = fetchLastFilename(); !found {
            filename = input("Filepath: ")
        }
    }

    // Determine the password to use.
    password = parser.GetStr("db-password")
    if password == "" {
        if password, found = fetchLastPassword(); !found {
            password = input("Password: ")
        }
    }
    cacheLastPassword(password)
    cacheLastFilename(filename)

    // Load the database file.
    db, key, err := irondb.Load(password, filename)
    if err != nil {
        exit("Error:", err)
    }

    // Read the JSON input file.
    input, err := ioutil.ReadFile(parser.GetArgs()[0])
    if err != nil {
        exit("Error:", err)
    }

    // Import the entries into the database.
    db.Import(key, string(input))
    db.Save(key, password, filename)
}
