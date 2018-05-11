package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
)


var urlHelp = fmt.Sprintf(`
Usage: %s url [FLAGS] [OPTIONS] ARGUMENTS

  Copy a stored url to the system clipboard or print it to stdout.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
  -p, --print               Print the url to stdout.
`, filepath.Base(os.Args[0]))


func registerUrl(parser *args.ArgParser) {
    cmd := parser.NewCmd("url", urlHelp, urlCallback)
    cmd.NewString("file f")
    cmd.NewFlag("print p")
}


func urlCallback(parser *args.ArgParser) {

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

    // Print the url to stdout.
    if parser.GetFlag("print") {
        fmt.Print(entry.Url)
        if stdoutIsTerminal() {
            fmt.Println()
        }
        return
    }

    // Copy the url to the clipboard.
    writeToClipboard(entry.Url)
}
