package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
)


const version = "1.3.1"


var helptext = fmt.Sprintf(`
Usage: %s [FLAGS] [COMMAND]

  Ironclad is a command line password manager.

Flags:
  -h, --help        Print the application's help text and exit.
  -v, --version     Print the application's version number and exit.

Basic Commands:
  add               Add a new entry to a password database.
  delete            Delete one or more entries from a database.
  edit              Edit an existing database entry.
  gen               Generate a new random password.
  init              Initialize a new password database.
  list              List database entries.
  new               Alias for 'add'.
  pass              Copy a password to the clipboard.
  show              Alias for 'list --verbose'.
  url               Copy a url to the clipboard.
  user              Copy a username to the clipboard.

Additional Commands:
  config            Set or print a configuration option.
  decrypt           Decrypt a file.
  dump              Dump a database's internal JSON data store.
  encrypt           Encrypt a file.
  export            Export entries from a database.
  import            Import entries into a database.
  purge             Purge inactive (i.e. deleted) entries from a database.
  restore           Restore one or more previously deleted entries.
  setpass           Change a database's master password.
  tags              List database tags.

Command Help:
  help <command>    Print the specified command's help text and exit.
`, filepath.Base(os.Args[0]))


func main() {

    // Instantiate an argument parser.
    parser := args.NewParser()
    parser.Helptext = helptext
    parser.Version = version

    // Register commands.
    registerAdd(parser)
    registerCache(parser)
    registerConfig(parser)
    registerDecrypt(parser)
    registerDelete(parser)
    registerDump(parser)
    registerEdit(parser)
    registerEncrypt(parser)
    registerExport(parser)
    registerGen(parser)
    registerImport(parser)
    registerInit(parser)
    registerList(parser)
    registerPass(parser)
    registerPurge(parser)
    registerRestore(parser)
    registerSetpass(parser)
    registerTags(parser)
    registerUrl(parser)
    registerUser(parser)

    // Parse the command line arguments.
    parser.Parse()
    if !parser.HasCmd() {
        parser.ExitHelp()
    }
}
