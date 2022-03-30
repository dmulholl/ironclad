package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo"
)

var exportHelp = fmt.Sprintf(`
Usage: %s export [entries]

  Exports a list of entries in JSON format. Entries can be specified by ID or
  by title. (Titles are checked for a case-insensitive substring match.) If no
  entries are specified, all entries will be exported.

Arguments:
  [entries]                 List of entries to export by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.
  -o, --out <str>           Output filename. Defaults to 'passwords.json'.
  -t, --tag <str>           Filter entries using the specified tag.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerExportCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("export")
	cmdParser.Helptext = exportHelp
	cmdParser.Callback = exportCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewStringOption("tag t", "")
	cmdParser.NewStringOption("out o", "passwords.json")
}

func exportCallback(cmdName string, cmdParser *argo.ArgParser) {

	// Load the database.
	filename, _, db := loadDB(cmdParser)

	// Default to exporting all active entries.
	list := db.Active()

	// Do we have query strings to filter on?
	if cmdParser.HasArgs() {
		list = list.FilterByAny(cmdParser.Args()...)
	}

	// Are we filtering by tag?
	if cmdParser.StringValue("tag") != "" {
		list = list.FilterByTag(cmdParser.StringValue("tag"))
	}

	// Confirm export.
	printCompact(list, db.Size(), filepath.Base(filename))
	answer := input("  Export the entries listed above? (y/n): ")
	if strings.ToLower(answer) != "y" {
		fmt.Println("  Export aborted.")
		os.Exit(0)
	}

	// Create the JSON dump.
	jsonstr, err := list.Export()
	if err != nil {
		exit(err)
	}

	// Write to file.
	err = ioutil.WriteFile(cmdParser.StringValue("out"), []byte(jsonstr), 0644)
	if err != nil {
		exit(err)
	}
}
