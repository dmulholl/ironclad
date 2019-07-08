package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
)


var restoreHelp = fmt.Sprintf(`
Usage: %s restore [FLAGS] [OPTIONS] ARGUMENTS

  Restore one or more deleted (i.e. inactive) entries. Entries to restore
  should be specified by ID.

Arguments:
  <entries>                 List of entry IDs.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerRestoreCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("restore", restoreHelp, restoreCallback)
    cmd.NewString("file f")
}


func restoreCallback(parser *janus.ArgParser) {

    // Check that at least one entry argument has been supplied.
    if !parser.HasArgs() {
        exit("you must specify at least one entry to restore")
    }

    // Load the database.
    filePath, password, db := loadDB(parser)
    fileName := filepath.Base(filePath)

    // Grab the entries to restore.
    list := db.Inactive().FilterByIDString(parser.GetArgs()...)
    if len(list) == 0 {
        exit("no matching entries")
    }

    // Print a listing and request confirmation.
    printCompact(list, len(db.Inactive()), fileName)
    answer := input("  Restore the entries listed above? (y/n): ")
    if strings.ToLower(answer) == "y" {
        for _, entry := range list {
            db.SetActive(entry.Id)
        }
        fmt.Println("  Entries restored.")
        printLineOfChar("─")
    } else {
        fmt.Println("  Restore aborted.")
        printLineOfChar("─")
    }

    // Save the updated database to disk.
    saveDB(filePath, password, db)
}
