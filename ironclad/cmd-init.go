package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
)


import (
    "github.com/dmulholland/ironclad/irondb"
)


var initHelp = fmt.Sprintf(`
Usage: %s init [FLAGS] ARGUMENTS

  Create a new encrypted password database. You will be prompted to supply
  a master password.

Arguments:
  <file>                    Filename for the new database.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerInit(parser *args.ArgParser) {
    parser.NewCmd("init", initHelp, initCallback)
}


func initCallback(parser *args.ArgParser) {

    // Check that a filename argument has been supplied.
    if !parser.HasArgs() {
        exit("you must supply a filename for the database")
    }
    filename := parser.GetArgs()[0]

    // Prompt for a master password for the new database.
    password := inputPass("Master Password: ")

    // Initialize a new database.
    db := irondb.New()

    // Cache the filename. We don't cache the password when creating a
    // new database file in case the user has accidentally mistyped it.
    setCachedFilename(filename)

    // Save the new database to disk.
    saveDB(filename, password, db)
}
