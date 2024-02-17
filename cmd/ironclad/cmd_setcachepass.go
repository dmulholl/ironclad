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
	filename, masterpass, db := loadDB(cmdParser)

	printLineOfChar("─")
	newCachePass := inputPass("Enter new cache password: ")
	confirmNewCachePass := inputPass("     Re-enter to confirm: ")
	printLineOfChar("─")

	if newCachePass != confirmNewCachePass {
		return fmt.Errorf("passwords do not match")
	}

	db.CachePass = newCachePass
	saveDB(filename, masterpass, db)
	setCachedPassword(filename, masterpass, newCachePass)

	return nil
}