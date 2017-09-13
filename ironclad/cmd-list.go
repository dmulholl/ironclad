package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
)


var listHelp = fmt.Sprintf(`
Usage: %s list [FLAGS] [OPTIONS] [ARGUMENTS]

  Print a list of entries from a database. Entries to list can be specified by
  ID or by title. (Titles are checked for a case-insensitive substring match.)

  If no arguments are specified, all the entries in the database will be
  listed.

  The 'list' command has an alias, 'show', which is equivalent to:

    list --verbose

Arguments:
  [entries]                 Entries to list by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.
  -t, --tag <str>           Filter entries using the specified tag.

Flags:
  -h, --help                Print this command's help text and exit.
  -v, --verbose             Use the verbose list format.
`, filepath.Base(os.Args[0]))


func listCallback(parser *args.ArgParser) {

    // Load the database.
    _, _, db := loadDB(parser)

    // Default to displaying all active entries.
    list := db.Active()
    title := "All Entries"

    // Do we have query strings to filter on?
    if parser.HasArgs() {
        list = list.FilterByAny(parser.GetArgs()...)
        title = "Matching Entries"
    }

    // Are we filtering by tag?
    if parser.GetString("tag") != "" {
        list = list.FilterByTag(parser.GetString("tag"))
        title = "Matching Entries"
    }

    // Print the list of entries.
    if parser.GetFlag("verbose") || parser.GetParent().GetCmdName() == "show" {
        printVerbose(list, db.Size(), title)
    } else {
        printCompact(list, db.Size())
    }
}
