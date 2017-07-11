package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
)


import (
    "github.com/dmulholland/ironclad/ironconfig"
)


// Application version number.
const version = "0.17.1"


// Application help text.
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
  list              List database entries.
  new               Create a new password database.
  pass              Copy a password to the clipboard.
  purge             Purge deleted entries from a database.
  tags              List database tags.
  user              Copy a username to the clipboard.

Command Help:
  help <command>    Print the specified command's help text and exit.
`, filepath.Base(os.Args[0]))


// Path to the application's configuration directory.
var configdir = filepath.Join(os.Getenv("HOME"), ".config", "ironclad")


// Path to the application's configuration file.
var configfile = filepath.Join(configdir, "goconfig.toml")


// Default port for the cached-password server.
const defaultport = "54313"


// Application entry point.
func main() {

    // Set the location of the application's configuration file.
    ironconfig.ConfigFile = configfile

    // Instantiate an argument parser.
    parser := clio.NewParser(helptext, version)

    // Register the 'add' command.
    addParser := parser.AddCmd("add", addHelp, addCallback)
    addParser.AddStr("file f", "")
    addParser.AddStr("masterpass", "")

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
    deleteParser.AddStr("masterpass", "")

    // Register the 'dump' command.
    dumpParser := parser.AddCmd("dump", dumpHelp, dumpCallback)
    dumpParser.AddStr("file f", "")
    dumpParser.AddStr("masterpass", "")

    // Register the 'edit' command.
    editParser := parser.AddCmd("edit", editHelp, editCallback)
    editParser.AddStr("file f", "")
    editParser.AddStr("masterpass", "")
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
    exportParser.AddStr("masterpass", "")
    exportParser.AddStr("tag t", "")

    // Register the 'gen' command.
    genParser := parser.AddCmd("gen", genHelp, genCallback)
    genParser.AddStr("file f", "")
    genParser.AddStr("masterpass", "")
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
    importParser.AddStr("masterpass", "")

    // Register the 'list' command.
    listParser := parser.AddCmd("list show", listHelp, listCallback)
    listParser.AddStr("file f", "")
    listParser.AddStr("masterpass", "")
    listParser.AddStr("tag t", "")
    listParser.AddFlag("verbose v")

    // Register the 'new' command.
    newParser := parser.AddCmd("new", newHelp, newCallback)
    newParser.AddStr("masterpass", "")

    // Register the 'pass' command.
    passParser := parser.AddCmd("pass", passHelp, passCallback)
    passParser.AddStr("file f", "")
    passParser.AddStr("masterpass", "")
    passParser.AddFlag("readable r")
    passParser.AddFlag("print p")

    // Register the 'purge' command.
    purgeParser := parser.AddCmd("purge", purgeHelp, purgeCallback)
    purgeParser.AddStr("file f", "")
    purgeParser.AddStr("masterpass", "")

    // Register the 'tags' command.
    tagsParser := parser.AddCmd("tags", tagsHelp, tagsCallback)
    tagsParser.AddStr("file f", "")
    tagsParser.AddStr("masterpass", "")

    // Register the 'user' command.
    userParser := parser.AddCmd("user", userHelp, userCallback)
    userParser.AddStr("file f", "")
    userParser.AddStr("masterpass", "")
    userParser.AddFlag("print p")

    // Parse the application's command line arguments.  If a command is found,
    // control will be passed to its callback function.
    parser.Parse()

    // If no command has been found, print the help text and exit.
    if !parser.HasCmd() {
        parser.Help()
    }
}
