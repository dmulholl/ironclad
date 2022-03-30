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

  Prints a list of entries from a database, showing only the entry title.

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
  -i, --inactive            List inactive entries.
`, filepath.Base(os.Args[0]))

func registerListCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("list")
	cmdParser.Helptext = listHelp
	cmdParser.Callback = listCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewStringOption("tag t", "")
	cmdParser.NewFlag("inactive i")
}

func listCallback(cmdName string, cmdParser *argo.ArgParser) {
	filename, _, db := loadDB(cmdParser)

	// Default to displaying all active entries.
	var list irondb.EntryList
	var totalCount int
	if cmdParser.Found("inactive") {
		list = db.Inactive()
		totalCount = len(list)
	} else {
		list = db.Active()
		totalCount = len(list)
	}

	// Do we have query strings to filter on?
	if cmdParser.HasArgs() {
		list = list.FilterByAny(cmdParser.Args()...)
	}

	// Are we filtering by tag?
	if cmdParser.StringValue("tag") != "" {
		list = list.FilterByTag(cmdParser.StringValue("tag"))
	}

	// Print the list of entries.
	printCompact(list, totalCount, filepath.Base(filename))
}
