package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
	"github.com/dmulholl/ironclad/irondb"
)

var showHelp = fmt.Sprintf(`
Usage: %s show [entries]

  Prints a list of entries from a database, showing entry content.

  Entries to list can be specified by ID or by title. (Titles are checked for
  a case-insensitive substring match.)

  If no arguments are specified, all the entries in the database will be
  listed.

Arguments:
  [entries]                 Entries to list by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.
  -t, --tag <str>           Filter entries using the specified tag.

Flags:
  -h, --help                Print this command's help text and exit.
  -i, --inactive            Show inactive entries.
  -n, --show-notes          Show each entry's notes.
  -p, --show-password       Show each entry's password in clear text.
`, filepath.Base(os.Args[0]))

func registerShowCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("show")
	cmdParser.Helptext = showHelp
	cmdParser.Callback = showCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewStringOption("tag t", "")
	cmdParser.NewFlag("inactive i")
	cmdParser.NewFlag("show-password p")
	cmdParser.NewFlag("show-notes n")
}

func showCallback(cmdName string, cmdParser *argo.ArgParser) {
	filename, _, db := loadDB(cmdParser)

	// Default to displaying all active entries.
	var list irondb.EntryList
	var title string
	var totalCount int
	if cmdParser.Found("inactive") {
		list = db.Inactive()
		title = "Inactive Entries"
		totalCount = len(list)
	} else {
		list = db.Active()
		title = "All Entries"
		totalCount = len(list)
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
	printVerbose(
		list,
		totalCount,
		cmdParser.Found("show-password"),
		cmdParser.Found("show-notes"),
		title,
		filepath.Base(filename))
}
