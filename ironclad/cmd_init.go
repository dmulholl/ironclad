package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
	"github.com/dmulholl/ironclad/irondb"
)

var initHelp = fmt.Sprintf(`
Usage: %s init <file>

  Creates a new encrypted password database. You will be prompted to supply
  a master password for the database.

Arguments:
  <file>                    Filename for the new database.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerInitCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("init")
	cmdParser.Helptext = initHelp
	cmdParser.Callback = initCallback
}

func initCallback(cmdName string, cmdParser *argo.ArgParser) {
	if !cmdParser.HasArgs() {
		exit("you must supply a filename for the database")
	}
	filename := cmdParser.Arg(0)

	masterpass1 := inputPass("Enter the master password for the new database: ")
	masterpass2 := inputPass("                           Re-enter to confirm: ")
	if masterpass1 != masterpass2 {
		exit("the passwords do not match")
	}

	cachepass1 := inputPass("Enter the cache password for the new database: ")
	cachepass2 := inputPass("                          Re-enter to confirm: ")
	if cachepass1 != cachepass2 {
		exit("the passwords do not match")
	}

	db := irondb.New(cachepass1)
	setCachedFilename(filename)
	saveDB(filename, masterpass1, db)
}
