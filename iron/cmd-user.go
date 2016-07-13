package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
)


// Help text for the 'user' command.
var userHelptext = fmt.Sprintf(`
Usage: %s user [FLAGS] [OPTIONS] ARGUMENTS

  Copy a stored username to the system clipboard or print it to stdout.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
      --help                Print this command's help text and exit.
  -p, --print               Print the username to stdout.
`, filepath.Base(os.Args[0]))


// Callback for the 'user' command.
func userCallback(parser *clio.ArgParser) {

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("Error: missing entry argument.")
    }

    // Load the database.
    _, _, db := loadDB(parser)

    // Search for an entry corresponding to the specified argument.
    list := db.Active().FilterProgressive(parser.GetArgs()[0])
    if len(list) == 0 {
        exit("no matching entry")
    } else if len(list) > 1 {
        exit("query matches multiple entries")
    }
    entry := list[0]

    // Print the username to stdout.
    if parser.GetFlag("print") {
        fmt.Print(entry.Username)
        if stdoutIsTerminal() {
            fmt.Println()
        }
        return
    }

    // Copy the username to the clipboard.
    writeToClipboard(entry.Username)
}
