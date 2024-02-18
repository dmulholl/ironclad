package main

import (
	"fmt"

	"github.com/dmulholl/argo/v4"
)

var masterpassCmdHelptext = `
Usage: ironclad setmasterpass

  Changes a database's master password.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerSetMasterPassCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("setmasterpass")
	cmdParser.Helptext = masterpassCmdHelptext
	cmdParser.Callback = masterpassCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func masterpassCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	_, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	printLineOfChar("─")
	newMasterPass := inputPass("Enter new master password: ")
	confirmNewMasterPass := inputPass("      Re-enter to confirm: ")
	printLineOfChar("─")

	if newMasterPass != confirmNewMasterPass {
		return fmt.Errorf("passwords do not match")
	}

	if err := saveDB(filename, newMasterPass, db); err != nil {
		return err
	}

	setCachedPassword(filename, newMasterPass, db.CachePass)
	return nil
}
