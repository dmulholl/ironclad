package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


// Help text for the 'pass' command.
var passHelp = fmt.Sprintf(`
Usage: %s pass [FLAGS] [OPTIONS] ARGUMENTS

  Copy a stored password to the system clipboard or print it to stdout.

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

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("missing entry argument")
    }

    // Load the database.
    _, _, db := loadDB(parser)

    // Search for an entry corresponding to the specified argument.
    list := db.Active().FilterProgressive(parser.GetArg(0))
    if len(list) == 0 {
        exit("no matching entry")
    } else if len(list) > 1 {
        exit("query matches multiple entries")
    }
    entry := list[0]

    // Add spaces if required.
    password := entry.Password
    if parser.GetFlag("readable") {
        password = addSpaces(password)
    }

    // Print the password to stdout.
    if parser.GetFlag("print") {
        fmt.Print(password)
        if stdoutIsTerminal() {
            fmt.Println()
        }
        return
    }

    // Copy the password to the clipboard.
    writeToClipboard(password)
}
