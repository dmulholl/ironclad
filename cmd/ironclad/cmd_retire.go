package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo/v4"
)

var retireCmdHelptext = `
Usage: ironclad retire <entries>

  Retires one or more entries from the database. Entries to be retired should
  be specified by ID.

  Retired entries are marked as inactive and do not appear in normal listings
  but their data remains in the database. Inactive entries can be stripped
  from the database using the 'purge' command.

  You can view inactive entries using the 'list --inactive' command.

Arguments:
  <entries>                 List of entry IDs.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerRetireCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("retire")
	cmdParser.Helptext = retireCmdHelptext
	cmdParser.Callback = retireCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func retireCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing entry argument")
	}

	ids, err := cmdParser.ArgsAsInts()
	if err != nil {
		return fmt.Errorf("arguments must be integer IDs: %w", err)
	}

	filename, masterpass, db := loadDB(cmdParser)

	list := db.Active().FilterByID(ids...)
	if len(list) == 0 {
		return fmt.Errorf("no matching entries")
	}

	printCompact(list, db.Count(), filepath.Base(filename))
	answer := input("  Retire the entries listed above? (y/n): ")
	if strings.ToLower(answer) == "y" {
		for _, entry := range list {
			db.SetInactive(entry.Id)
		}
		saveDB(filename, masterpass, db)
		fmt.Println("  Entries retired.")
	} else {
		fmt.Println("  Operation aborted.")
	}
	printLineOfChar("─")

	return nil
}
