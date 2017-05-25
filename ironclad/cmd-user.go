package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


// Help text for the 'user' command.
var userHelp = fmt.Sprintf(`
Usage: %s user [FLAGS] [OPTIONS] ARGUMENTS

  Copy a stored username to the system clipboard or print it to stdout. This
  command will fall back on the email address if the username field is empty.

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

    // Return the email field if the username field is empty.
    user := entry.Username
    if user == "" {
        user = entry.Email
    }

    // Print the username to stdout.
    if parser.GetFlag("print") {
        fmt.Print(user)
        if stdoutIsTerminal() {
            fmt.Println()
        }
        return
    }

    // Copy the username to the clipboard.
    writeToClipboard(user)
}
