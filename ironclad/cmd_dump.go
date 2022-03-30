package main

import "github.com/dmulholl/janus/v2"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

var dumpHelp = fmt.Sprintf(`
Usage: %s dump

  Dump a database's internal JSON data store to stdout.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerDumpCmd(parser *janus.ArgParser) {
	cmd := parser.NewCmd("dump", dumpHelp, dumpCallback)
	cmd.NewString("file f")
}

func dumpCallback(parser *janus.ArgParser) {

	// Load the database.
	_, _, db := loadDB(parser)

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
