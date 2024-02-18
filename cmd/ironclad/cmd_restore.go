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

	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	masterpass, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	list := db.Inactive().FilterByID(ids...)
	if len(list) == 0 {
		return fmt.Errorf("no matching entries")
	}

	printCompactList(list, len(db.Inactive()), filepath.Base(filename))

	answer := input("  Restore the entries listed above? (y/n): ")
	if strings.ToLower(answer) != "y" {
		fmt.Println("  Operation aborted.")
		printLineOfChar("─")
		return nil
	}

	for _, entry := range list {
		db.SetActive(entry.Id)
	}

	if err := saveDB(filename, masterpass, db); err != nil {
		return err
	}

	fmt.Println("  Entries restored.")
	printLineOfChar("─")

	return nil
}
