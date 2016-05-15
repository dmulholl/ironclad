package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/ironclad/irondb"
    "github.com/dmulholland/clio/go/clio"
    "github.com/atotto/clipboard"
)


// Help text for the 'user' command.
var userHelptext = fmt.Sprintf(`
Usage: %s user [FLAGS] [OPTIONS] ARGUMENTS

  Print a username.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file.

Flags:
  -c, --clipboard           Write the username to the system clipboard.
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'user' command.
func userCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("Error: missing entry argument.")
    }

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

    // Load the database file.
    db, _, err := irondb.Load(password, filename)
    if err != nil {
        exit("Error:", err)
    }
    cacheLastPassword(password)
    cacheLastFilename(filename)

    // Search for an entry corresponding to the specified argument.
    entries := db.Lookup(parser.GetArgs()[0])
    if len(entries) == 0 {
        exit("Error: no matching entry.")
    } else if len(entries) > 1 {
        exit("Error: query matches multiple entries.")
    }

    // Print to the clipboard or stdout.
    if parser.GetFlag("clipboard") {
        if clipboard.Unsupported {
            exit("Error: clipboard not supported on this system.")
        }
        err := clipboard.WriteAll(entries[0].Username)
        if err != nil {
            exit("Error:", err)
        }
    } else {
        fmt.Print(entries[0].Username)
        if stdoutIsTerminal() {
            fmt.Println()
        }
    }
}
