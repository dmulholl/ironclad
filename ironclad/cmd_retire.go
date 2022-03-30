package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo"
)

var retireHelp = fmt.Sprintf(`
Usage: %s retire <entries>

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
`, filepath.Base(os.Args[0]))

func registerRetireCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("retire")
	cmdParser.Helptext = retireHelp
	cmdParser.Callback = retireCallback
	cmdParser.NewStringOption("file f", "")
}

func retireCallback(cmdName string, cmdParser *argo.ArgParser) {

	// Check that at least one entry argument has been supplied.
	if !cmdParser.HasArgs() {
		exit("you must specify at least one entry to retire")
	}
	filename, masterpass, db := loadDB(cmdParser)

	// Grab the entries to retire.
	list := db.Active().FilterByIDString(cmdParser.Args()...)
	if len(list) == 0 {
		exit("no matching entries")
	}

	// Print a listing and request confirmation.
	printCompact(list, db.Size(), filepath.Base(filename))
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
}
