package main

import (
	"fmt"
	"path/filepath"

	"github.com/dmulholl/argo/v4"
	"github.com/pkg/browser"
)

var goCmdHelptext = `
Usage: ironclad go <entry>

  Opens the URL for the specified entry in the default browser.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerGoCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("go")
	cmdParser.Helptext = goCmdHelptext
	cmdParser.Callback = goCmdCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("print p")
}

func goCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("mising entry argument")
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
		return fmt.Errorf("no matching entry")
	}

	if len(matchingEntries) > 1 {
		fmt.Println("The query string matches multiple entries:")
		printCompactList(matchingEntries, len(db.Active()), filepath.Base(filename))
		return nil
	}

	err = browser.OpenURL(matchingEntries[0].Url)
	if err != nil {
		return err
	}

	return nil
}
