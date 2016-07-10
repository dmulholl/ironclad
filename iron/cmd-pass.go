package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/ironclad/irondb"
    "github.com/dmulholland/clio/go/clio"
    "github.com/atotto/clipboard"
)


// Help text for the 'pass' command.
var passHelptext = fmt.Sprintf(`
Usage: %s pass [FLAGS] [OPTIONS] ARGUMENTS

  Copy a password to the system clipboard. The password can additionally be
  printed to stdout.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
  -p, --print               Print the password to stdout.
  -r, --readable            Add spaces to the password for readability.
`, filepath.Base(os.Args[0]))


// Callback for the 'pass' command.
func passCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("Error: missing entry argument.")
    }

    // Determine the filename to use.
    filename = parser.GetStr("file")
    if filename == "" {
        if filename, found = fetchLastFilename(); !found {
            filename = input("Filename: ")
        }
    }

    // Determine the password to use.
    password = parser.GetStr("db-password")
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
    cacheLastPassword(password)
    cacheLastFilename(filename)

    // Search for an entry corresponding to the specified argument.
    entries := db.LookupUnique(parser.GetArgs()[0])
    if len(entries) == 0 {
        exit("Error: no matching entry.")
    } else if len(entries) > 1 {
        exit("Error: query matches multiple entries.")
    }

    // Decrypt the stored password.
    decrypted, err := entries[0].GetPassword(key)
    if err != nil {
        exit("Error:", err)
    }

    // Add spaces if required.
    if parser.GetFlag("readable") {
        decrypted = addSpaces(decrypted)
    }

    // Print the password to stdout.
    if parser.GetFlag("print") {
        fmt.Print(decrypted)
        if stdoutIsTerminal() {
            fmt.Println()
        }
    }

    // Copy the password to the clipboard.
    if clipboard.Unsupported {
        exit("Error: clipboard not supported on this system.")
    }
    err = clipboard.WriteAll(decrypted)
    if err != nil {
        exit("Error:", err)
    }
}
