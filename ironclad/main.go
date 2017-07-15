package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


const version = "0.19.0.dev"


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
    addParser := parser.AddCmd("add", addHelp, addCallback)
    addParser.AddStr("file f", "")

    // Register the 'cache' command.
    parser.AddCmd("cache", cacheHelp, cacheCallback)

    // Register the 'config' command.
    parser.AddCmd("config", configHelp, configCallback)

    // Register the 'decrypt' command.
    decryptParser := parser.AddCmd("decrypt", decryptHelp, decryptCallback)
    decryptParser.AddStr("out o", "")

    // Register the 'delete' command.
    deleteParser := parser.AddCmd("delete", deleteHelp, deleteCallback)
    deleteParser.AddStr("file f", "")

    // Register the 'dump' command.
    dumpParser := parser.AddCmd("dump", dumpHelp, dumpCallback)
    dumpParser.AddStr("file f", "")

    // Register the 'edit' command.
    editParser := parser.AddCmd("edit", editHelp, editCallback)
    editParser.AddStr("file f", "")
    editParser.AddFlag("title t")
    editParser.AddFlag("url l")
    editParser.AddFlag("username u")
    editParser.AddFlag("password p")
    editParser.AddFlag("notes n")
    editParser.AddFlag("tags s")
    editParser.AddFlag("email e")

    // Register the 'encrypt' command.
    encryptParser := parser.AddCmd("encrypt", encryptHelp, encryptCallback)
    encryptParser.AddStr("out o", "")

    // Register the 'export' command.
    exportParser := parser.AddCmd("export", exportHelp, exportCallback)
    exportParser.AddStr("file f", "")
    exportParser.AddStr("tag t", "")

    // Register the 'gen' command.
    genParser := parser.AddCmd("gen", genHelp, genCallback)
    genParser.AddStr("file f", "")
    genParser.AddFlag("digits d")
    genParser.AddFlag("exclude-similar x")
    genParser.AddFlag("lowercase l")
    genParser.AddFlag("symbols s")
    genParser.AddFlag("uppercase u")
    genParser.AddFlag("readable r")
    genParser.AddFlag("print p")

    // Register the 'import' command.
    importParser := parser.AddCmd("import", importHelp, importCallback)
    importParser.AddStr("file f", "")

    // Register the 'init' command.
    parser.AddCmd("init", initHelp, initCallback)

    // Register the 'list' command.
    listParser := parser.AddCmd("list show", listHelp, listCallback)
    listParser.AddStr("file f", "")
    listParser.AddStr("tag t", "")
    listParser.AddFlag("verbose v")

    // Register the 'pass' command.
    passParser := parser.AddCmd("pass", passHelp, passCallback)
    passParser.AddStr("file f", "")
    passParser.AddFlag("readable r")
    passParser.AddFlag("print p")

    // Register the 'purge' command.
    purgeParser := parser.AddCmd("purge", purgeHelp, purgeCallback)
    purgeParser.AddStr("file f", "")

    // Register the 'tags' command.
    tagsParser := parser.AddCmd("tags", tagsHelp, tagsCallback)
    tagsParser.AddStr("file f", "")

    // Register the 'user' command.
    userParser := parser.AddCmd("user", userHelp, userCallback)
    userParser.AddStr("file f", "")
    userParser.AddFlag("print p")

    // Parse the command line arguments. If no command is found, print the
    // help text and exit.
    parser.Parse()
    if !parser.HasCmd() {
        parser.Help()
    }
}
