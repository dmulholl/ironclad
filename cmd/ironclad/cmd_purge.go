package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo/v4"
)

var purgeCmdHelptext = `
Usage: ironclad purge

  Purges all inactive entries from a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerPurgeCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("purge")
	cmdParser.Helptext = purgeCmdHelptext
	cmdParser.Callback = purgeCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func purgeCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	filename, masterpass, db := loadDB(cmdParser)

	list := db.Inactive()
	if len(list) == 0 {
		return fmt.Errorf("no inactive entries to purge")
	}

	printCompactList(list, len(list), filepath.Base(filename))
	answer := input("  Purge the entries listed above? (y/n): ")
	if strings.ToLower(answer) == "y" {
		db.PurgeInactive()
		fmt.Println("  Entries purged.")
		printLineOfChar("─")
	} else {
		fmt.Println("  Purge aborted.")
		printLineOfChar("─")
	}

	saveDB(filename, masterpass, db)
	return nil
}
