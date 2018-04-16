package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
)


const version = "1.1.0-dev"


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
  user              Copy a username to the clipboard.

Additional Commands:
  config            Set or print a configuration option.
  decrypt           Decrypt a file.
  dump              Dump a database's internal JSON data store.
  encrypt           Encrypt a file.
  export            Export entries from a database.
  import            Import entries into a database.
  purge             Purge deleted entries from a database.
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

    // Register the 'add' command.
    addCmd := parser.NewCmd("add new", addHelp, addCallback)
    addCmd.NewString("file f")
    addCmd.NewFlag("no-editor")

    // Register the 'cache' command.
    parser.NewCmd("cache", cacheHelp, cacheCallback)

    // Register the 'config' command.
    parser.NewCmd("config", configHelp, configCallback)

    // Register the 'decrypt' command.
    decryptCmd := parser.NewCmd("decrypt", decryptHelp, decryptCallback)
    decryptCmd.NewString("out o")

    // Register the 'delete' command.
    deleteCmd := parser.NewCmd("delete", deleteHelp, deleteCallback)
    deleteCmd.NewString("file f")

    // Register the 'dump' command.
    dumpCmd := parser.NewCmd("dump", dumpHelp, dumpCallback)
    dumpCmd.NewString("file f")

    // Register the 'edit' command.
    editCmd := parser.NewCmd("edit", editHelp, editCallback)
    editCmd.NewString("file f")
    editCmd.NewFlag("title t")
    editCmd.NewFlag("url l")
    editCmd.NewFlag("username u")
    editCmd.NewFlag("password p")
    editCmd.NewFlag("notes n")
    editCmd.NewFlag("tags s")
    editCmd.NewFlag("email e")
    editCmd.NewFlag("no-editor")

    // Register the 'encrypt' command.
    encryptCmd := parser.NewCmd("encrypt", encryptHelp, encryptCallback)
    encryptCmd.NewString("out o")

    // Register the 'export' command.
    exportCmd := parser.NewCmd("export", exportHelp, exportCallback)
    exportCmd.NewString("file f")
    exportCmd.NewString("tag t")

    // Register the 'gen' command.
    genCmd := parser.NewCmd("gen", genHelp, genCallback)
    genCmd.NewString("file f")
    genCmd.NewFlag("digits d")
    genCmd.NewFlag("exclude-similar x")
    genCmd.NewFlag("lowercase l")
    genCmd.NewFlag("symbols s")
    genCmd.NewFlag("uppercase u")
    genCmd.NewFlag("readable r")
    genCmd.NewFlag("print p")

    // Register the 'import' command.
    importCmd := parser.NewCmd("import", importHelp, importCallback)
    importCmd.NewString("file f")

    // Register the 'init' command.
    parser.NewCmd("init", initHelp, initCallback)

    // Register the 'list' command.
    listCmd := parser.NewCmd("list show", listHelp, listCallback)
    listCmd.NewString("file f")
    listCmd.NewString("tag t")
    listCmd.NewFlag("verbose v")

    // Register the 'pass' command.
    passCmd := parser.NewCmd("pass", passHelp, passCallback)
    passCmd.NewString("file f")
    passCmd.NewFlag("readable r")
    passCmd.NewFlag("print p")

    // Register the 'purge' command.
    purgeCmd := parser.NewCmd("purge", purgeHelp, purgeCallback)
    purgeCmd.NewString("file f")

    // Register the 'tags' command.
    tagsCmd := parser.NewCmd("tags", tagsHelp, tagsCallback)
    tagsCmd.NewString("file f")

    // Register the 'user' command.
    userCmd := parser.NewCmd("user", userHelp, userCallback)
    userCmd.NewString("file f")
    userCmd.NewFlag("print p")

    // Register the 'setpass' command.
    setpassCmd := parser.NewCmd("setpass", setpassHelp, setpassCallback)
    setpassCmd.NewString("file f")

    // Parse the command line arguments. If no command is found, print the
    // help text and exit.
    parser.Parse()
    if !parser.HasCmd() {
        parser.ExitHelp()
    }
}
