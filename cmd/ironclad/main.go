/*
Ironclad: a command line tool for creating and managing encrypted password databases.
*/
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dmulholl/argo/v4"
)

const version = "2.7.0-rc3"

var helptext = `
Usage: ironclad [command]

  A utility for creating and managing encrypted password databases.

Flags:
  -h, --help        Print the application's help text.
  -v, --version     Print the application's version number.

Basic Commands:
  add               Add a new entry to a password database.
  edit              Edit an existing database entry.
  gen               Generate a new random password.
  go                Open an entry's URL in the default browser.
  init              Initialize a new password database.
  list              List database entries.
  pass              Copy a password to the clipboard.
  retire            Mark one or more entries as inactive.
  show              Show entry content.
  url               Copy a url to the clipboard.
  user              Copy a username to the clipboard.

Additional Commands:
  config            Set or print a configuration option.
  decrypt           Decrypt a file.
  dump              Dump a database's internal JSON data store.
  encrypt           Encrypt a file.
  export            Export entries from a database.
  import            Import entries into a database.
  purge             Purge inactive entries from a database.
  restore           Restore inactive entries to active status.
  setcachepass      Change a database's cache password.
  setmasterpass     Change a database's master password.
  tags              List database tags.

Aliases:
  new               Alias for the 'add' command.

Command Help:
  help <command>    Print the specified command's help text and exit.
`

func main() {
	parser := argo.NewParser()
	parser.Helptext = helptext
	parser.Version = version

	registerAddCmd(parser)
	registerCacheCmd(parser)
	registerConfigCmd(parser)
	registerDecryptCmd(parser)
	registerDumpCmd(parser)
	registerEditCmd(parser)
	registerEncryptCmd(parser)
	registerExportCmd(parser)
	registerGenCmd(parser)
	registerImportCmd(parser)
	registerInitCmd(parser)
	registerListCmd(parser)
	registerPassCmd(parser)
	registerPurgeCmd(parser)
	registerRestoreCmd(parser)
	registerRetireCmd(parser)
	registerSetMasterPassCmd(parser)
	registerSetCachePassCmd(parser)
	registerTagsCmd(parser)
	registerUrlCmd(parser)
	registerUserCmd(parser)
	registerShowCmd(parser)
	registerGoCmd(parser)

	if err := parser.ParseOsArgs(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	if parser.FoundCommandName == "" {
		fmt.Println(strings.TrimSpace(helptext))
	}
}
