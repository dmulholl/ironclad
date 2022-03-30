package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
)

var userHelp = fmt.Sprintf(`
Usage: %s user <entry>

  Copies a stored username to the system clipboard or prints it to stdout. This
  command will fall back on the email address if the username field is empty.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
  -p, --print               Print the username to stdout.
`, filepath.Base(os.Args[0]))

func registerUserCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("user")
	cmdParser.Helptext = userHelp
	cmdParser.Callback = userCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("print p")
}

func userCallback(cmdName string, cmdParser *argo.ArgParser) {
	if !cmdParser.HasArgs() {
		exit("missing entry argument")
	}
	filename, _, db := loadDB(cmdParser)

	// Search for an entry corresponding to the supplied arguments.
	list := db.Active().FilterByAll(cmdParser.Args...)
	if len(list) == 0 {
		exit("no matching entry")
	} else if len(list) > 1 {
		println("Error: the query string matches multiple entries.")
		printCompact(list, len(db.Active()), filepath.Base(filename))
		os.Exit(1)
	}
	entry := list[0]

	// Return the email field if the username field is empty.
	user := entry.Username
	if user == "" {
		user = entry.Email
	}

	// Print the username to stdout.
	if cmdParser.Found("print") {
		fmt.Print(user)
		if stdoutIsTerminal() {
			fmt.Println()
		}
		return
	}

	// Copy the username to the clipboard.
	writeToClipboard(user)
}
