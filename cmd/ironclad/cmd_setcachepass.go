package main

import (
	"fmt"

	"github.com/dmulholl/argo/v4"
)

var cachepassCmdHelptext = `
Usage: ironclad setcachepass

  Changes a database's cache password. This password is used to encrypt the
  master password while it's temporarily cached in memory.

  Note that if you set the cache password to an empty string you will not be
  prompted to enter it.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerSetCachePassCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("setcachepass")
	cmdParser.Helptext = cachepassCmdHelptext
	cmdParser.Callback = cachepassCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func cachepassCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	masterpass, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	printLineOfChar("─")

	newCachePass, err := inputMasked("Enter new cache password: ")
	if err != nil {
		return err
	}

	confirmNewCachePass, err := inputMasked("     Re-enter to confirm: ")
	if err != nil {
		return err
	}

	printLineOfChar("─")

	if newCachePass != confirmNewCachePass {
		return fmt.Errorf("passwords do not match")
	}

	db.CachePass = newCachePass

	if err := saveDB(filename, masterpass, db); err != nil {
		return err
	}

	return setCachedPassword(filename, masterpass, newCachePass)
}
