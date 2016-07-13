package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/irondb"
)


// Help text for the 'new' command.
var newHelptext = fmt.Sprintf(`
Usage: %s new [FLAGS] ARGUMENTS

  Create a new password database.

Arguments:
  <file>                    Filename for the new database.

Flags:
  --help                    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'new' command.
func newCallback(parser *clio.ArgParser) {

    // Check that a filename argument has been supplied.
    if !parser.HasArgs() {
        exit("you must supply a filename")
    }
    filename := parser.GetArgs()[0]

    // Prompt for a password if none has been supplied.
    password := parser.GetStr("db-password")
    if password == "" {
        password = input("Password: ")
    }

    // Initialize a new database.
    db, err := irondb.New()
    if err != nil {
        exit(err)
    }

    // Cache the password and filename.
    setCachedPassword(password)
    setCachedFilename(filename)

    // Save the new database to disk.
    saveDB(password, filename, db)
}
