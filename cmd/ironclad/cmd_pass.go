package main

import (
	"fmt"
	"path/filepath"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/textutils"
)

var passCmdHelptext = `
Usage: ironclad pass <entry>

  Copies a stored password to the system clipboard.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
  -p, --print               Print the password to the standard output stream.
  -r, --readable            Add spaces to the password for readability.
`

func registerPassCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("pass")
	cmdParser.Helptext = passCmdHelptext
	cmdParser.Callback = passCmdCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("readable r")
	cmdParser.NewFlag("print p")
}

func passCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
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
		return fmt.Errorf("no matching entry")
	}

	if len(matchingEntries) > 1 {
		fmt.Println("The query string matches multiple entries:")
		printCompactList(matchingEntries, len(db.Active()), filepath.Base(filename))
		return nil
	}

	entry := matchingEntries[0]

	password := entry.GetPassword()
	if cmdParser.Found("readable") {
		password = textutils.AddSpacer(password, "  ", 4)
	}

	if cmdParser.Found("print") {
		fmt.Print(password)
		if stdoutIsTerminal() {
			fmt.Println()
		}
		return nil
	}

	return writeToClipboardWithTimeout(password)
}
