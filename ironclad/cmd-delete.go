package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
)


var deleteHelp = fmt.Sprintf(`
Usage: %s delete [FLAGS] [OPTIONS] ARGUMENTS

  Delete entries from a database. Entries should be specified by ID.

Arguments:
  <entries>                 List of entry IDs.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func deleteCallback(parser *args.ArgParser) {

    // Check that at least one entry argument has been supplied.
    if !parser.HasArgs() {
        exit("you must specify at least one entry to delete")
    }

    // Load the database.
    filename, password, db := loadDB(parser)

    // Grab the entries to delete.
    list := db.Active().FilterByIDString(parser.GetArgs()...)
    if len(list) == 0 {
        exit("no matching entries")
    }

    // Print a listing and request confirmation.
    printCompact(list, db.Size())
    answer := input("  Delete the entries listed above? (y/n): ")
    if strings.ToLower(answer) == "y" {
        for _, entry := range list {
            db.Delete(entry.Id)
        }
        fmt.Println("  Entries deleted.")
        line("─")
    } else {
        fmt.Println("  Deletion aborted.")
        line("─")
    }

    // Save the updated database to disk.
    saveDB(filename, password, db)
}
