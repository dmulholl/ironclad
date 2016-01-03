/*
    Command line tool for managing Ironclad password databases.

    Author: Darren Mulholland <darren@mulholland.xyz>
    License: MIT
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
const version = "0.2.2"


// Application help text.
var helptext = fmt.Sprintf(`
Usage: %s [FLAGS] [COMMAND]

  Command line utility for managing Ironclad password databases.

Flags:
  --help            Print the application's help text and exit.
  --version         Print the application's version number and exit.

Commands:
  add               Add a new entry to a database.
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
  user              Print a username.

Command Help:
  help <command>    Print the specified command's help text and exit.
`, filepath.Base(os.Args[0]))


// Path to the application's home directory.
var ironpath = filepath.Join(os.Getenv("HOME"), ".ironclad")


// Address for the cached-password server.
const ironaddress = "localhost:54512"


// Application entry point.
func main() {

    // Set the location of the application's configuration file.
    ironconfig.Configfile = filepath.Join(ironpath, "goconfig.toml")

    // Instantiate an argument parser.
    parser := clio.NewParser(helptext, version)

    // Register the 'add' command.
    addParser := parser.AddCommand("add", addCallback, addHelptext)
    addParser.AddStringOption("file", "", 'f')
    addParser.AddStringOption("db-password", "")

    // Register the 'delete' command.
    deleteParser := parser.AddCommand("delete", deleteCallback, deleteHelptext)
    deleteParser.AddStringOption("file", "", 'f')
    deleteParser.AddStringOption("db-password", "")

    // Register the 'dump' command.
    dumpParser := parser.AddCommand("dump", dumpCallback, dumpHelptext)
    dumpParser.AddStringOption("file", "", 'f')
    dumpParser.AddStringOption("db-password", "")

    // Register the 'edit' command.
    editParser := parser.AddCommand("edit", editCallback, editHelptext)
    editParser.AddStringOption("file", "", 'f')
    editParser.AddStringOption("db-password", "")
    editParser.AddFlag("title")
    editParser.AddFlag("url")
    editParser.AddFlag("username")
    editParser.AddFlag("password")
    editParser.AddFlag("notes")
    editParser.AddFlag("tags")

    // Register the 'export' command.
    exportParser := parser.AddCommand("export", exportCallback, exportHelptext)
    exportParser.AddStringOption("file", "", 'f')
    exportParser.AddStringOption("db-password", "")

    // Register the 'gen' command.
    genParser := parser.AddCommand("gen", genCallback, genHelptext)
    genParser.AddStringOption("file", "", 'f')
    genParser.AddStringOption("db-password", "")
    genParser.AddFlag("digits", 'd')
    genParser.AddFlag("exclude", 'e')
    genParser.AddFlag("lowercase", 'l')
    genParser.AddFlag("symbols", 's')
    genParser.AddFlag("uppercase", 'u')

    // Register the 'import' command.
    importParser := parser.AddCommand("import", importCallback, importHelptext)
    importParser.AddStringOption("file", "", 'f')
    importParser.AddStringOption("db-password", "")

    // Register the 'list' command.
    listParser := parser.AddCommand("list", listCallback, listHelptext)
    listParser.AddStringOption("file", "", 'f')
    listParser.AddStringOption("db-password", "")
    listParser.AddStringOption("tag", "", 't')
    listParser.AddFlag("verbose", 'v')

    // Register the 'new' command.
    newParser := parser.AddCommand("new", newCallback, newHelptext)
    newParser.AddStringOption("db-password", "")

    // Register the 'pass' command.
    passParser := parser.AddCommand("pass", passCallback, passHelptext)
    passParser.AddStringOption("file", "", 'f')
    passParser.AddStringOption("db-password", "")

    // Register the 'purge' command.
    purgeParser := parser.AddCommand("purge", purgeCallback, purgeHelptext)
    purgeParser.AddStringOption("file", "", 'f')
    purgeParser.AddStringOption("db-password", "")

    // Register the 'serve' command.
    parser.AddCommand("serve", serveCallback, serveHelptext)

    // Register the 'user' command.
    userParser := parser.AddCommand("user", userCallback, userHelptext)
    userParser.AddStringOption("file", "", 'f')
    userParser.AddStringOption("db-password", "")

    // Parse the application's command line arguments.
    // If a command is found, control will be passed to its
    // callback function.
    parser.Parse()

    // If no command has been found, print the help text and exit.
    if !parser.HasCommand() {
        parser.Help()
    }
}
