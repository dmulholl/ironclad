package main

import (
	"path/filepath"

	"github.com/dmulholl/argo/v4"
)

var showCmdHelptext = `
Usage: ironclad show [entries]

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
  -p, --show-passwords      Show each entry's password in clear text.
`

func registerShowCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("show")
	cmdParser.Helptext = showCmdHelptext
	cmdParser.Callback = showCmdCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewStringOption("tag t", "")
	cmdParser.NewFlag("inactive i")
	cmdParser.NewFlag("show-passwords p")
	cmdParser.NewFlag("show-notes n")
}

func showCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	filename, _, db := loadDB(cmdParser)

	// Default to displaying active entries.
	list := db.Active()
	title := "All Entries"
	totalCount := len(list)

	if cmdParser.Found("inactive") {
		list = db.Inactive()
		title = "Inactive Entries"
		totalCount = len(list)
	}

	if len(cmdParser.Args) > 0 {
		list = list.FilterByAny(cmdParser.Args...)
		title = "Matching Entries"
	}

	if cmdParser.StringValue("tag") != "" {
		list = list.FilterByTag(cmdParser.StringValue("tag"))
		title = "Matching Entries"
	}

	printVerbose(
		list,
		totalCount,
		cmdParser.Found("show-passwords"),
		cmdParser.Found("show-notes"),
		title,
		filepath.Base(filename),
	)

	return nil
}
