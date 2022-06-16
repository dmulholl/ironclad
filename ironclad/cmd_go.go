package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
	"github.com/pkg/browser"
)

var goHelp = fmt.Sprintf(`
Usage: %s go <entry>

  Opens the URL for the specified entry in the default browser.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerGoCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("go")
	cmdParser.Helptext = goHelp
	cmdParser.Callback = goCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("print p")
}

func goCallback(cmdName string, cmdParser *argo.ArgParser) {
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

	err := browser.OpenURL(entry.Url)
	if err != nil {
		fmt.Printf("Error: {}.\n", err)
		os.Exit(1)
	}
}
