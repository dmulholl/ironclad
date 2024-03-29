package main

import (
	"fmt"
	"os"

	"github.com/dmulholl/argo/v4"
)

var importCmdHelptext = `
Usage: ironclad import <file>

  Imports a previously-exported list of entries in JSON format.

Arguments:
  <file>                    File to import.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerImportCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("import")
	cmdParser.Helptext = importCmdHelptext
	cmdParser.Callback = importCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func importCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing filename argument")
	}

	input, err := os.ReadFile(cmdParser.Args[0])
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	masterpass, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	count, err := db.Import(input)
	if err != nil {
		return fmt.Errorf("failed to import entries: %w", err)
	}

	if err := saveDB(filename, masterpass, db); err != nil {
		return err
	}

	fmt.Printf("%d entries imported.\n", count)
	return nil
}
