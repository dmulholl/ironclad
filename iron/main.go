/*
    Command line tool for managing Ironclad password databases.
*/
package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/ironconfig"
)


// Application version number.
const version = "0.13.0"


// Application help text.
var helptext = fmt.Sprintf(`
Usage: %s [FLAGS] [COMMAND]

  Command line utility for managing Ironclad password databases.

Flags:
  --help            Print the application's help text and exit.
  --version         Print the application's version number and exit.

Commands:
  add               Add a new entry to a database.
  config            Set or print a configuration option.
  delete            Delete an entry from a database.
  dump              Dump a database's internal JSON data store.
  edit              Edit an existing database entry.
  export            Export data from a database.
  gen               Generate a new random password.
  import            Import data into a database.
  list              List database entries.
  new               Create a new database.
  pass              Copy a password to the clipboard.
  purge             Purge deleted entries from a database.
  tags              List database tags.
  user              Copy a username to the clipboard.

Command Help:
  help <command>    Print the specified command's help text and exit.
`, filepath.Base(os.Args[0]))


// Path to the application's configuration directory.
var ironpath = filepath.Join(os.Getenv("HOME"), ".config", "ironclad")


// Path to the application's configuration file.
var configfile = filepath.Join(ironpath, "goconfig.toml")


// Address for the cached-password server.
const ironaddress = "localhost:54512"


// Application entry point.
func main() {

    // Set the location of the application's configuration file.
    ironconfig.ConfigFile = configfile

    // Instantiate an argument parser.
    parser := clio.NewParser(helptext, version)

    // Register the 'add' command.
    addParser := parser.AddCmd("add", addHelptext, addCallback)
    addParser.AddStr("file f", "")
    addParser.AddStr("db-password", "")

    // Register the 'cache' command.
    parser.AddCmd("cache", cacheHelptext, cacheCallback)

    // Register the 'config' command.
    parser.AddCmd("config", configHelptext, configCallback)

    // Register the 'delete' command.
    deleteParser := parser.AddCmd("delete", deleteHelptext, deleteCallback)
    deleteParser.AddStr("file f", "")
    deleteParser.AddStr("db-password", "")

    // Register the 'dump' command.
    dumpParser := parser.AddCmd("dump", dumpHelptext, dumpCallback)
    dumpParser.AddStr("file f", "")
    dumpParser.AddStr("db-password", "")

    // Register the 'edit' command.
    editParser := parser.AddCmd("edit", editHelptext, editCallback)
    editParser.AddStr("file f", "")
    editParser.AddStr("db-password", "")
    editParser.AddFlag("title t")
    editParser.AddFlag("url l")
    editParser.AddFlag("username u")
    editParser.AddFlag("password p")
    editParser.AddFlag("notes n")
    editParser.AddFlag("tags s")
    editParser.AddFlag("email e")

    // Register the 'export' command.
    exportParser := parser.AddCmd("export", exportHelptext, exportCallback)
    exportParser.AddStr("file f", "")
    exportParser.AddStr("db-password", "")
    exportParser.AddStr("tag t", "")

    // Register the 'gen' command.
    genParser := parser.AddCmd("gen", genHelptext, genCallback)
    genParser.AddStr("file f", "")
    genParser.AddStr("db-password", "")
    genParser.AddFlag("digits d")
    genParser.AddFlag("exclude-similar d")
    genParser.AddFlag("lowercase l")
    genParser.AddFlag("symbols s")
    genParser.AddFlag("uppercase u")
    genParser.AddFlag("readable r")
    genParser.AddFlag("print p")

    // Register the 'import' command.
    importParser := parser.AddCmd("import", importHelptext, importCallback)
    importParser.AddStr("file f", "")
    importParser.AddStr("db-password", "")

    // Register the 'list' command.
    listParser := parser.AddCmd("list show", listHelptext, listCallback)
    listParser.AddStr("file f", "")
    listParser.AddStr("db-password", "")
    listParser.AddStr("tag t", "")
    listParser.AddFlag("verbose v")
    listParser.AddFlag("cleartext c")

    // Register the 'new' command.
    newParser := parser.AddCmd("new", newHelptext, newCallback)
    newParser.AddStr("db-password", "")

    // Register the 'pass' command.
    passParser := parser.AddCmd("pass", passHelptext, passCallback)
    passParser.AddStr("file f", "")
    passParser.AddStr("db-password", "")
    passParser.AddFlag("readable r")
    passParser.AddFlag("print p")

    // Register the 'purge' command.
    purgeParser := parser.AddCmd("purge", purgeHelptext, purgeCallback)
    purgeParser.AddStr("file f", "")
    purgeParser.AddStr("db-password", "")

    // Register the 'tags' command.
    tagsParser := parser.AddCmd("tags", tagsHelptext, tagsCallback)
    tagsParser.AddStr("file f", "")
    tagsParser.AddStr("db-password", "")

    // Register the 'user' command.
    userParser := parser.AddCmd("user", userHelptext, userCallback)
    userParser.AddStr("file f", "")
    userParser.AddStr("db-password", "")
    userParser.AddFlag("print p")

    // Parse the application's command line arguments.  If a command is found,
    // control will be passed to its callback function.
    parser.Parse()

    // If no command has been found, print the help text and exit.
    if !parser.HasCmd() {
        parser.Help()
    }
}
