package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/janus/v2"
)

var importHelp = fmt.Sprintf(`
Usage: %s import <file>

  Import a previously-exported list of entries in JSON format.

Arguments:
  <file>                    File to import.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerImportCmd(parser *janus.ArgParser) {
	cmd := parser.NewCmd("import", importHelp, importCallback)
	cmd.NewString("file f")
}

func importCallback(parser *janus.ArgParser) {
	var input []byte
	var err error
	if parser.HasArgs() {
		input, err = ioutil.ReadFile(parser.GetArg(0))
		if err != nil {
			exit(err)
		}
	} else {
		exit("you must specify a file to import")
	}

	// V1 exports could sometimes get polluted with leading asterisks
	// from the password input.
	trimmed := strings.Trim(string(input), "*")
	input = []byte(trimmed)

	filename, masterpass, db := loadDB(parser)
	count, err := db.Import(input)
	if err != nil {
		exit(err)
	}
	saveDB(filename, masterpass, db)
	fmt.Printf("%d entries imported.\n", count)
}
