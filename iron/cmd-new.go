package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/ironcrypt"
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
        exit("Error: you must supply a filename.")
    }
    filename := parser.GetArgs()[0]

    // Prompt for a password if none has been supplied.
    password := parser.GetStrOpt("db-password")
    if password == "" {
        password = input("Password: ")
    }

    // Create a new in-memory database.
    db := irondb.New()

    // Generate a random master key.
    key, err := ironcrypt.RandBytes(ironcrypt.KeySize)
    if err != nil {
        exit("Error:", err)
    }

    // Save the new database as an encrypted file.
    err = db.Save(key, password, filename)
    if err != nil {
        exit("Error:", err)
    }

    // Cache the password and filename.
    cacheLastPassword(password)
    cacheLastFilename(filename)
}
