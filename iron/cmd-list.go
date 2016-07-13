package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
)


// Help text for the 'list' command.
var listHelptext = fmt.Sprintf(`
Usage: %s list [FLAGS] [OPTIONS] [ARGUMENTS]

  Print a list of entries from a database. Entries to list can be specified by
  ID or by title. (Titles are checked for a case-insensitive substring match.)

  If no arguments are specified, all the entries in the database will be
  listed.

  The 'list' command has an alias, 'show', which is equivalent to:

    list --verbose --cleartext

Arguments:
  [entries]                 Entries to list by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.
  -t, --tag <str>           Filter entries using the specified tag.

Flags:
  -c, --cleartext           Print passwords in cleartext.
      --help                Print this command's help text and exit.
  -v, --verbose             Use the verbose list format.
`, filepath.Base(os.Args[0]))


// Callback for the 'list' command.
func listCallback(parser *clio.ArgParser) {

    // Load the database.
    password, _, db := loadDB(parser)

    // Has the 'show' alias been used?
    if parser.GetParent().GetCmdName() == "show" {
        parser.SetFlag("verbose", true)
        parser.SetFlag("cleartext", true)
    }

    // Default to displaying all active entries.
    list := db.Active()
    title := "All Entries"

    // Do we have query strings to filter on?
    if parser.HasArgs() {
        list = list.FilterByQuery(parser.GetArgs()...)
        title = "Matching Entries"
    }

    // Are we filtering by tag?
    if parser.GetStr("tag") != "" {
        list = list.FilterByTag(parser.GetStr("tag"))
        title = "Matching Entries"
    }

    // Print the list of entries.
    if parser.GetFlag("verbose") {
        cleartext := parser.GetFlag("cleartext")
        printVerbose(list, db.Size(), db.Key(password), title, cleartext)
    } else {
        printCompact(list, db.Size())
    }
}
