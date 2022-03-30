package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo"
)

var restoreHelp = fmt.Sprintf(`
Usage: %s restore <entries>

  Restores one or more inactive entries to active status. Entries to restore
  should be specified by ID.

Arguments:
  <entries>                 List of entry IDs.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerRestoreCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("restore")
	cmdParser.Helptext = restoreHelp
	cmdParser.Callback = restoreCallback
	cmdParser.NewStringOption("file f", "")
}

func restoreCallback(cmdName string, cmdParser *argo.ArgParser) {
	if !cmdParser.HasArgs() {
		exit("you must specify at least one entry to restore")
	}
	filename, masterpass, db := loadDB(cmdParser)

	// Grab the entries to restore.
	list := db.Inactive().FilterByIDString(cmdParser.Args...)
	if len(list) == 0 {
		exit("no matching entries")
	}

	// Print a listing and request confirmation.
	printCompact(list, len(db.Inactive()), filepath.Base(filename))
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
}
