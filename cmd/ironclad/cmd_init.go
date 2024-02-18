package main

import (
	"fmt"
	"path/filepath"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/config"
	"github.com/dmulholl/ironclad/internal/database"
)

var initCmdHelptext = `
Usage: ironclad init <file>

  Creates a new encrypted password database.

  You will be prompted to enter a master password which will be used to encrypt
  the database file, and a cache password which will be used to encrypt the
  master password while it's temporarily cached in memory.

Arguments:
  <file>                    Filename for the new database.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerInitCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("init")
	cmdParser.Helptext = initCmdHelptext
	cmdParser.Callback = initCmdCallback
}

func initCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing filename argument")
	}

	filename, err := filepath.Abs(cmdParser.Args[0])
	if err != nil {
		return fmt.Errorf("failed to resolve filename: %w", err)
	}

	masterpass1 := inputPass("Enter the master password for the new database: ")
	masterpass2 := inputPass("                           Re-enter to confirm: ")
	if masterpass1 != masterpass2 {
		return fmt.Errorf("the master passwords do not match")
	}

	cachepass1 := inputPass("Enter the cache password for the new database: ")
	cachepass2 := inputPass("                          Re-enter to confirm: ")
	if cachepass1 != cachepass2 {
		return fmt.Errorf("the cache passwords do not match")
	}

	db := database.New(cachepass1)
	saveDB(filename, masterpass1, db)

	err = config.Set("file", filename)
	if err != nil {
		return fmt.Errorf("failed to cache filename: %w", err)
	}

	return nil
}
