package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo"
)

var purgeHelp = fmt.Sprintf(`
Usage: %s purge

  Purges all inactive entries from a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerPurgeCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("purge")
	cmdParser.Helptext = purgeHelp
	cmdParser.Callback = purgeCallback
	cmdParser.NewStringOption("file f", "")
}

func purgeCallback(cmdName string, cmdParser *argo.ArgParser) {
	filename, masterpass, db := loadDB(cmdParser)

	list := db.Inactive()
	if len(list) == 0 {
		exit("no inactive entries to purge")
	}

	printCompact(list, len(list), filepath.Base(filename))
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
}
