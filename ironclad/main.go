package main


import "github.com/dmulholl/janus/v2"


import (
    "fmt"
    "os"
    "path/filepath"
)


const version = "2.2.2"


var helptext = fmt.Sprintf(`
Usage: %s [FLAGS] [COMMAND]

  A utility for creating and managing encrypted password databases.

Flags:
  -h, --help        Print the application's help text.
  -v, --version     Print the application's version number.

Basic Commands:
  add               Add a new entry to a password database.
  edit              Edit an existing database entry.
  gen               Generate a new random password.
  init              Initialize a new password database.
  list              List database entries.
  pass              Copy a password to the clipboard.
  retire            Mark one or more entries as inactive.
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
  new               Alias for 'add'.
  show              Alias for 'list --verbose'.

Command Help:
  help <command>    Print the specified command's help text and exit.
`, filepath.Base(os.Args[0]))


func main() {

    // Instantiate an argument parser.
    parser := janus.NewParser()
    parser.Helptext = helptext
    parser.Version = version

    // Register commands.
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

    // Parse the command line arguments.
    parser.Parse()
    if !parser.HasCmd() {
        parser.ExitHelp()
    }
}
