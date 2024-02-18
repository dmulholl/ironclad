package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo/v4"
)

var restoreCmdHelptext = `
Usage: ironclad restore <entries>

  Restores one or more inactive entries to active status. Entries to restore
  should be specified by ID.

Arguments:
  <entries>                 List of entry IDs.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerRestoreCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("restore")
	cmdParser.Helptext = restoreCmdHelptext
	cmdParser.Callback = restoreCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func restoreCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing entry argument")
	}

	ids, err := cmdParser.ArgsAsInts()
	if err != nil {
		return fmt.Errorf("arguments must be integer IDs: %w", err)
	}

	filename, masterpass, db := loadDB(cmdParser)

	list := db.Inactive().FilterByID(ids...)
	if len(list) == 0 {
		return fmt.Errorf("no matching entries")
	}

	printCompactList(list, len(db.Inactive()), filepath.Base(filename))
	answer := input("  Restore the entries listed above? (y/n): ")
	if strings.ToLower(answer) == "y" {
		for _, entry := range list {
			db.SetActive(entry.Id)
		}
		saveDB(filename, masterpass, db)
		fmt.Println("  Entries restored.")
	} else {
		fmt.Println("  Restore aborted.")
	}
	printLineOfChar("─")

	return nil
}
