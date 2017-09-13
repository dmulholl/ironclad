package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
)


var userHelp = fmt.Sprintf(`
Usage: %s user [FLAGS] [OPTIONS] ARGUMENTS

  Copy a stored username to the system clipboard or print it to stdout. This
  command will fall back on the email address if the username field is empty.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
  -p, --print               Print the username to stdout.
`, filepath.Base(os.Args[0]))


func userCallback(parser *args.ArgParser) {

    // Make sure we have at least one argument.
    if !parser.HasArgs() {
        exit("missing entry argument")
    }

    // Load the database.
    _, _, db := loadDB(parser)

    // Search for an entry corresponding to the supplied arguments.
    list := db.Active().FilterByAll(parser.GetArgs()...)
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
