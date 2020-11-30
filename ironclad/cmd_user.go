package main


import "github.com/dmulholl/janus/v2"


import (
    "fmt"
    "os"
    "path/filepath"
)


var userHelp = fmt.Sprintf(`
Usage: %s user <entry>

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


func registerUserCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("user", userHelp, userCallback)
    cmd.NewString("file f")
    cmd.NewFlag("print p")
}


func userCallback(parser *janus.ArgParser) {
    if !parser.HasArgs() {
        exit("missing entry argument")
    }
    filename, _, db := loadDB(parser)

    // Search for an entry corresponding to the supplied arguments.
    list := db.Active().FilterByAll(parser.GetArgs()...)
    if len(list) == 0 {
        exit("no matching entry")
    } else if len(list) > 1 {
        println("Error: the query string matches multiple entries.")
        printCompact(list, len(db.Active()), filepath.Base(filename))
        os.Exit(1)
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
