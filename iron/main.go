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
const version = "0.6.1"


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
  gen               Generate a random password.
  import            Import data into a database.
  list              List database entries.
  new               Create a new database.
  pass              Print a password.
  purge             Purge deleted entries from a database.
  tags              List database tags.
  user              Print a username.

Command Help:
  help <command>    Print the specified command's help text and exit.
`, filepath.Base(os.Args[0]))


// Path to the application's configuration directory.
var ironpath = filepath.Join(os.Getenv("HOME"), ".config", "ironclad")


// Address for the cached-password server.
const ironaddress = "localhost:54512"


// Application entry point.
func main() {

    // Set the location of the application's configuration file.
    ironconfig.Configfile = filepath.Join(ironpath, "goconfig.toml")

    // Instantiate an argument parser.
    parser := clio.NewParser(helptext, version)

    // Register the 'add' command.
    addParser := parser.AddCmd("add", addCallback, addHelptext)
    addParser.AddStrOpt("file", "", 'f')
    addParser.AddStrOpt("db-password", "")

    // Register the 'cache' command.
    parser.AddCmd("cache", cacheCallback, cacheHelptext)

    // Register the 'config' command.
    parser.AddCmd("config", configCallback, configHelptext)

    // Register the 'delete' command.
    deleteParser := parser.AddCmd("delete", deleteCallback, deleteHelptext)
    deleteParser.AddStrOpt("file", "", 'f')
    deleteParser.AddStrOpt("db-password", "")

    // Register the 'dump' command.
    dumpParser := parser.AddCmd("dump", dumpCallback, dumpHelptext)
    dumpParser.AddStrOpt("file", "", 'f')
    dumpParser.AddStrOpt("db-password", "")

    // Register the 'edit' command.
    editParser := parser.AddCmd("edit", editCallback, editHelptext)
    editParser.AddStrOpt("file", "", 'f')
    editParser.AddStrOpt("db-password", "")
    editParser.AddFlag("title", 't')
    editParser.AddFlag("url", 'l')
    editParser.AddFlag("username", 'u')
    editParser.AddFlag("password", 'p')
    editParser.AddFlag("notes", 'n')
    editParser.AddFlag("tags", 's')

    // Register the 'export' command.
    exportParser := parser.AddCmd("export", exportCallback, exportHelptext)
    exportParser.AddStrOpt("file", "", 'f')
    exportParser.AddStrOpt("db-password", "")

    // Register the 'gen' command.
    genParser := parser.AddCmd("gen", genCallback, genHelptext)
    genParser.AddStrOpt("file", "", 'f')
    genParser.AddStrOpt("db-password", "")
    genParser.AddFlag("digits", 'd')
    genParser.AddFlag("exclude-similar", 'e')
    genParser.AddFlag("lowercase", 'l')
    genParser.AddFlag("symbols", 's')
    genParser.AddFlag("uppercase", 'u')
    genParser.AddFlag("readable", 'r')

    // Register the 'import' command.
    importParser := parser.AddCmd("import", importCallback, importHelptext)
    importParser.AddStrOpt("file", "", 'f')
    importParser.AddStrOpt("db-password", "")

    // Register the 'list' command.
    listParser := parser.AddCmd("list", listCallback, listHelptext)
    listParser.AddStrOpt("file", "", 'f')
    listParser.AddStrOpt("db-password", "")
    listParser.AddStrOpt("tag", "", 't')
    listParser.AddFlag("verbose", 'v')

    // Register the 'new' command.
    newParser := parser.AddCmd("new", newCallback, newHelptext)
    newParser.AddStrOpt("db-password", "")

    // Register the 'pass' command.
    passParser := parser.AddCmd("pass", passCallback, passHelptext)
    passParser.AddStrOpt("file", "", 'f')
    passParser.AddStrOpt("db-password", "")

    // Register the 'purge' command.
    purgeParser := parser.AddCmd("purge", purgeCallback, purgeHelptext)
    purgeParser.AddStrOpt("file", "", 'f')
    purgeParser.AddStrOpt("db-password", "")

    // Register the 'tags' command.
    tagsParser := parser.AddCmd("tags", tagsCallback, tagsHelptext)
    tagsParser.AddStrOpt("file", "", 'f')
    tagsParser.AddStrOpt("db-password", "")

    // Register the 'user' command.
    userParser := parser.AddCmd("user", userCallback, userHelptext)
    userParser.AddStrOpt("file", "", 'f')
    userParser.AddStrOpt("db-password", "")

    // Parse the application's command line arguments.
    // If a command is found, control will be passed to its
    // callback function.
    parser.Parse()

    // If no command has been found, print the help text and exit.
    if !parser.HasCmd() {
        parser.Help()
    }
}
