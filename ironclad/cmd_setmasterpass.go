package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
)

var masterpassHelp = fmt.Sprintf(`
Usage: %s setmasterpass

  Changes a database's master password.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerSetMasterPassCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("setmasterpass")
	cmdParser.Helptext = masterpassHelp
	cmdParser.Callback = masterpassCallback
	cmdParser.NewStringOption("file f", "")
}

func masterpassCallback(cmdName string, cmdParser *argo.ArgParser) {
	filename, _, db := loadDB(cmdParser)

	printLineOfChar("─")
	newMasterPass := inputPass("Enter new master password: ")
	confirmNewMasterPass := inputPass("      Re-enter to confirm: ")
	printLineOfChar("─")

	if newMasterPass == confirmNewMasterPass {
		saveDB(filename, newMasterPass, db)
		setCachedPassword(filename, newMasterPass, db.CachePass)
	} else {
		exit("passwords do not match")
	}
}
