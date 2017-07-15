package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


const version = "0.19.1"


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
    addCmd := parser.AddCmd("add", addHelp, addCallback)
    addCmd.AddStr("file f", "")

    // Register the 'cache' command.
    parser.AddCmd("cache", cacheHelp, cacheCallback)

    // Register the 'config' command.
    parser.AddCmd("config", configHelp, configCallback)

    // Register the 'decrypt' command.
    decryptCmd := parser.AddCmd("decrypt", decryptHelp, decryptCallback)
    decryptCmd.AddStr("out o", "")

    // Register the 'delete' command.
    deleteCmd := parser.AddCmd("delete", deleteHelp, deleteCallback)
    deleteCmd.AddStr("file f", "")

    // Register the 'dump' command.
    dumpCmd := parser.AddCmd("dump", dumpHelp, dumpCallback)
    dumpCmd.AddStr("file f", "")

    // Register the 'edit' command.
    editCmd := parser.AddCmd("edit", editHelp, editCallback)
    editCmd.AddStr("file f", "")
    editCmd.AddFlag("title t")
    editCmd.AddFlag("url l")
    editCmd.AddFlag("username u")
    editCmd.AddFlag("password p")
    editCmd.AddFlag("notes n")
    editCmd.AddFlag("tags s")
    editCmd.AddFlag("email e")

    // Register the 'encrypt' command.
    encryptCmd := parser.AddCmd("encrypt", encryptHelp, encryptCallback)
    encryptCmd.AddStr("out o", "")

    // Register the 'export' command.
    exportCmd := parser.AddCmd("export", exportHelp, exportCallback)
    exportCmd.AddStr("file f", "")
    exportCmd.AddStr("tag t", "")

    // Register the 'gen' command.
    genCmd := parser.AddCmd("gen", genHelp, genCallback)
    genCmd.AddStr("file f", "")
    genCmd.AddFlag("digits d")
    genCmd.AddFlag("exclude-similar x")
    genCmd.AddFlag("lowercase l")
    genCmd.AddFlag("symbols s")
    genCmd.AddFlag("uppercase u")
    genCmd.AddFlag("readable r")
    genCmd.AddFlag("print p")

    // Register the 'import' command.
    importCmd := parser.AddCmd("import", importHelp, importCallback)
    importCmd.AddStr("file f", "")

    // Register the 'init' command.
    parser.AddCmd("init", initHelp, initCallback)

    // Register the 'list' command.
    listCmd := parser.AddCmd("list show", listHelp, listCallback)
    listCmd.AddStr("file f", "")
    listCmd.AddStr("tag t", "")
    listCmd.AddFlag("verbose v")

    // Register the 'pass' command.
    passCmd := parser.AddCmd("pass", passHelp, passCallback)
    passCmd.AddStr("file f", "")
    passCmd.AddFlag("readable r")
    passCmd.AddFlag("print p")

    // Register the 'purge' command.
    purgeCmd := parser.AddCmd("purge", purgeHelp, purgeCallback)
    purgeCmd.AddStr("file f", "")

    // Register the 'tags' command.
    tagsCmd := parser.AddCmd("tags", tagsHelp, tagsCallback)
    tagsCmd.AddStr("file f", "")

    // Register the 'user' command.
    userCmd := parser.AddCmd("user", userHelp, userCallback)
    userCmd.AddStr("file f", "")
    userCmd.AddFlag("print p")

    // Parse the command line arguments. If no command is found, print the
    // help text and exit.
    parser.Parse()
    if !parser.HasCmd() {
        parser.Help()
    }
}
