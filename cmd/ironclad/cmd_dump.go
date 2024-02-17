package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/dmulholl/argo/v4"
)

var dumpCmdHelptext = `
Usage: ironclad dump

  Dumps a database's internal JSON data store to stdout.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerDumpCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("dump")
	cmdParser.Helptext = dumpCmdHelptext
	cmdParser.Callback = dumpCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func dumpCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	_, _, db := loadDB(cmdParser)

	data, err := db.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize database as JSON: %w", err)
	}

	var formatted bytes.Buffer

	err = json.Indent(&formatted, data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize database as JSON: %w", err)
	}

	fmt.Println(formatted.String())
	return nil
}
