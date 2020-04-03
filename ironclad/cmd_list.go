package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
)


import (
    "github.com/dmulholl/ironclad/irondb"
)


var listHelp = fmt.Sprintf(`
Usage: %s list [FLAGS] [OPTIONS] [ARGUMENTS]

  Print a list of entries from a database. Entries to list can be specified by
  ID or by title. (Titles are checked for a case-insensitive substring match.)

  If no arguments are specified, all the entries in the database will be
  listed.

  The 'list' command has an alias, 'show', which is equivalent to:

    list --verbose

  This will display the full details for each entry listed.

Arguments:
  [entries]                 Entries to list by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.
  -t, --tag <str>           Filter entries using the specified tag.

Flags:
  -d, --deleted             List deleted (i.e. inactive) entries.
  -h, --help                Print this command's help text and exit.
  -v, --verbose             Use the verbose list format.
`, filepath.Base(os.Args[0]))


func registerListCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("list show", listHelp, listCallback)
    cmd.NewString("file f")
    cmd.NewString("tag t")
    cmd.NewFlag("verbose v")
    cmd.NewFlag("deleted d")
}


func listCallback(parser *janus.ArgParser) {
    filename, _, db := loadDB(parser)

    // Default to displaying all active entries.
    var list irondb.EntryList
    var title string
    var count int
    if parser.GetFlag("deleted") {
        list = db.Inactive()
        title = "Deleted Entries"
        count = len(list)
    } else {
        list = db.Active()
        title = "All Entries"
        count = len(list)
    }

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
        printVerbose(list, count, title, filepath.Base(filename))
    } else {
        printCompact(list, count, filepath.Base(filename))
    }
}
