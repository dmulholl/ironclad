package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
)


var importHelp = fmt.Sprintf(`
Usage: %s import [FLAGS] [OPTIONS] ARGUMENTS

  Import a list of entries in JSON format.

Arguments:
  <file>                    File to import.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func importCallback(parser *args.ArgParser) {

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("you must supply the name of a file to import")
    }

    // Load the database.
    filename, password, db := loadDB(parser)

    // Read the JSON input file.
    input, err := ioutil.ReadFile(parser.GetArgs()[0])
    if err != nil {
        exit(err)
    }

    // Import the entries into the database.
    err = db.Import(input)
    if err != nil {
        exit(err)
    }
    // Save the updated database to disk.
    saveDB(filename, password, db)
}
