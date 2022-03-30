package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
)

var cachepassHelp = fmt.Sprintf(`
Usage: %s setcachepass

  Changes a database's cache password. This password is used to encrypt the
  master password while it's temporarily cached in memory.

  Note that if you set the cache password to an empty string you will not be
  prompted to enter it.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerSetCachePassCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("setcachepass")
	cmdParser.Helptext = cachepassHelp
	cmdParser.Callback = cachepassCallback
	cmdParser.NewStringOption("file f", "")
}

func cachepassCallback(cmdName string, cmdParser *argo.ArgParser) {
	filename, masterpass, db := loadDB(cmdParser)

	printLineOfChar("─")
	newCachePass := inputPass("Enter new cache password: ")
	confirmNewCachePass := inputPass("     Re-enter to confirm: ")
	printLineOfChar("─")

	if newCachePass == confirmNewCachePass {
		db.CachePass = newCachePass
		saveDB(filename, masterpass, db)
		setCachedPassword(filename, masterpass, newCachePass)
	} else {
		exit("passwords do not match")
	}
}
