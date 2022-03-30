package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
	"github.com/dmulholl/ironclad/irondb"
)

var listHelp = fmt.Sprintf(`
Usage: %s list [entries]

  Prints a list of entries from a database. Entries to list can be specified by
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
  -h, --help                Print this command's help text and exit.
  -i, --inactive            List inactive entries.
  -v, --verbose             Use the verbose list format.
`, filepath.Base(os.Args[0]))

func registerListCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("list show")
	cmdParser.Helptext = listHelp
	cmdParser.Callback = listCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewStringOption("tag t", "")
	cmdParser.NewFlag("verbose v")
	cmdParser.NewFlag("inactive i")
}

func listCallback(cmdName string, cmdParser *argo.ArgParser) {
	filename, _, db := loadDB(cmdParser)

	// Default to displaying all active entries.
	var list irondb.EntryList
	var title string
	var count int
	if cmdParser.Found("inactive") {
		list = db.Inactive()
		title = "Inactive Entries"
		count = len(list)
	} else {
		list = db.Active()
		title = "All Entries"
		count = len(list)
	}

	// Do we have query strings to filter on?
	if cmdParser.HasArgs() {
		list = list.FilterByAny(cmdParser.Args()...)
		title = "Matching Entries"
	}

	// Are we filtering by tag?
	if cmdParser.StringValue("tag") != "" {
		list = list.FilterByTag(cmdParser.StringValue("tag"))
		title = "Matching Entries"
	}

	// Print the list of entries.
	if cmdParser.Found("verbose") || cmdName == "show" {
		printVerbose(list, count, title, filepath.Base(filename))
	} else {
		printCompact(list, count, filepath.Base(filename))
	}
}
