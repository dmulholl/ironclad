package main

import (
	"fmt"
	"path/filepath"

	"github.com/dmulholl/argo/v4"
)

var urlCmdHelptext = `
Usage: ironclad url <entry>

  Copies a stored url to the system clipboard or prints it to stdout.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
  -p, --print               Print the url to stdout.
`

func registerUrlCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("url")
	cmdParser.Helptext = urlCmdHelptext
	cmdParser.Callback = urlCmdCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("print p")
}

func urlCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing entry argument")
	}

	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	_, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	matchingEntries := db.Active().FilterByAll(cmdParser.Args...)
	if len(matchingEntries) == 0 {
		return fmt.Errorf("no matching entries")
	}

	if len(matchingEntries) > 1 {
		fmt.Println("The query string matches multiple entries:")
		printCompactList(matchingEntries, len(db.Active()), filepath.Base(filename))
		return nil
	}

	entry := matchingEntries[0]

	if cmdParser.Found("print") {
		fmt.Print(entry.Url)
		if stdoutIsTerminal() {
			fmt.Println()
		}
		return nil
	}

	writeToClipboard(entry.Url)
	return nil
}
