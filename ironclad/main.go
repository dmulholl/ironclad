package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


const version = "0.21.2"


var helptext = fmt.Sprintf(`
Usage: %s [FLAGS] [COMMAND]

  Ironclad is a command line password manager.

Flags:
  --help            Print the application's help text and exit.
  --version         Print the application's version number and exit.

Commands:
  add               Add a new entry to a database.
  config            Set or print a configuration option.
  decrypt           Decrypt a file.
  delete            Delete entries from a database.
  dump              Dump a database's internal JSON data store.
  edit              Edit an existing database entry.
  encrypt           Encrypt a file.
  export            Export entries from a database.
  gen               Generate a random password.
  import            Import entries into a database.
  init              Initialize a new password database.
  list              List database entries.
  new               Add a new entry to a database.
  pass              Copy a password to the clipboard.
  purge             Purge deleted entries from a database.
  tags              List database tags.
  user              Copy a username to the clipboard.

Command Help:
  help <command>    Print the specified command's help text and exit.
`, filepath.Base(os.Args[0]))


func main() {

    // Instantiate an argument parser.
    parser := clio.NewParser(helptext, version)

    // Register the 'add' command.
    addCmd := parser.NewCmd("add new", addHelp, addCallback)
    addCmd.NewStr("file f", "")
    addCmd.NewFlag("no-editor")

    // Register the 'cache' command.
    parser.NewCmd("cache", cacheHelp, cacheCallback)

    // Register the 'config' command.
    parser.NewCmd("config", configHelp, configCallback)

    // Register the 'decrypt' command.
    decryptCmd := parser.NewCmd("decrypt", decryptHelp, decryptCallback)
    decryptCmd.NewStr("out o", "")

    // Register the 'delete' command.
    deleteCmd := parser.NewCmd("delete", deleteHelp, deleteCallback)
    deleteCmd.NewStr("file f", "")

    // Register the 'dump' command.
    dumpCmd := parser.NewCmd("dump", dumpHelp, dumpCallback)
    dumpCmd.NewStr("file f", "")

    // Register the 'edit' command.
    editCmd := parser.NewCmd("edit", editHelp, editCallback)
    editCmd.NewStr("file f", "")
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
    encryptCmd.NewStr("out o", "")

    // Register the 'export' command.
    exportCmd := parser.NewCmd("export", exportHelp, exportCallback)
    exportCmd.NewStr("file f", "")
    exportCmd.NewStr("tag t", "")

    // Register the 'gen' command.
    genCmd := parser.NewCmd("gen", genHelp, genCallback)
    genCmd.NewStr("file f", "")
    genCmd.NewFlag("digits d")
    genCmd.NewFlag("exclude-similar x")
    genCmd.NewFlag("lowercase l")
    genCmd.NewFlag("symbols s")
    genCmd.NewFlag("uppercase u")
    genCmd.NewFlag("readable r")
    genCmd.NewFlag("print p")

    // Register the 'import' command.
    importCmd := parser.NewCmd("import", importHelp, importCallback)
    importCmd.NewStr("file f", "")

    // Register the 'init' command.
    parser.NewCmd("init", initHelp, initCallback)

    // Register the 'list' command.
    listCmd := parser.NewCmd("list show", listHelp, listCallback)
    listCmd.NewStr("file f", "")
    listCmd.NewStr("tag t", "")
    listCmd.NewFlag("verbose v")

    // Register the 'pass' command.
    passCmd := parser.NewCmd("pass", passHelp, passCallback)
    passCmd.NewStr("file f", "")
    passCmd.NewFlag("readable r")
    passCmd.NewFlag("print p")

    // Register the 'purge' command.
    purgeCmd := parser.NewCmd("purge", purgeHelp, purgeCallback)
    purgeCmd.NewStr("file f", "")

    // Register the 'tags' command.
    tagsCmd := parser.NewCmd("tags", tagsHelp, tagsCallback)
    tagsCmd.NewStr("file f", "")

    // Register the 'user' command.
    userCmd := parser.NewCmd("user", userHelp, userCallback)
    userCmd.NewStr("file f", "")
    userCmd.NewFlag("print p")

    // Parse the command line arguments. If no command is found, print the
    // help text and exit.
    parser.Parse()
    if !parser.HasCmd() {
        parser.ExitHelp()
    }
}
