package main


import "github.com/dmulholland/janus-go/janus"


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
)


var deleteHelp = fmt.Sprintf(`
Usage: %s delete [FLAGS] [OPTIONS] ARGUMENTS

  Delete one or more entries from a database. Entries to delete should be
  specified by ID.

  Deleted entries are marked as inactive and do not appear in normal listings
  but their data remains in the database. Inactive entries can be stripped
  from the database using the 'purge' command.

  You can view inactive entries using the 'list --deleted' command.

Arguments:
  <entries>                 List of entry IDs.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerDeleteCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("delete", deleteHelp, deleteCallback)
    cmd.NewString("file f")
}


func deleteCallback(parser *janus.ArgParser) {

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
            db.SetInactive(entry.Id)
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
