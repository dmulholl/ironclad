package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo/v4"
)

var exportCmdHelptext = `
Usage: ironclad export [entries]

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
`

func registerExportCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("export")
	cmdParser.Helptext = exportCmdHelptext
	cmdParser.Callback = exportCmdCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewStringOption("tag t", "")
	cmdParser.NewStringOption("out o", "passwords.json")
}

func exportCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	_, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	// Default to exporting all active entries.
	list := db.Active()

	// Do we have query strings to filter on?
	if len(cmdParser.Args) > 0 {
		list = list.FilterByAny(cmdParser.Args...)
	}

	// Are we filtering by tag?
	if cmdParser.StringValue("tag") != "" {
		list = list.FilterByTag(cmdParser.StringValue("tag"))
	}

	// Confirm export.
	printCompactList(list, db.Count(), filepath.Base(filename))

	answer, err := input("  Export the entries listed above? (y/n): ")
	if err != nil {
		return err
	}

	if strings.ToLower(answer) != "y" {
		fmt.Println("  Export aborted.")
		return nil
	}

	jsonstr, err := list.Export()
	if err != nil {
		return fmt.Errorf("failed to export entries: %w", err)
	}

	err = os.WriteFile(cmdParser.StringValue("out"), []byte(jsonstr), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
