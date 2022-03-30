package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
)

var dumpHelp = fmt.Sprintf(`
Usage: %s dump

  Dumps a database's internal JSON data store to stdout.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerDumpCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("dump")
	cmdParser.Helptext = dumpHelp
	cmdParser.Callback = dumpCallback
	cmdParser.NewStringOption("file f", "")
}

func dumpCallback(cmdName string, cmdParser *argo.ArgParser) {
	// Load the database.
	_, _, db := loadDB(cmdParser)

	// Serialize the database as a byte-slice of JSON.
	data, err := db.ToJSON()
	if err != nil {
		exit(err)
	}

	// Format the JSON for display.
	var formatted bytes.Buffer
	json.Indent(&formatted, data, "", "  ")

	// Print the formatted JSON to stdout.
	fmt.Println(formatted.String())
}
