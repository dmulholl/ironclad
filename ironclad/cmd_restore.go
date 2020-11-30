package main


import "github.com/dmulholl/janus/v2"


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
)


var restoreHelp = fmt.Sprintf(`
Usage: %s restore <entries>

  Restore one or more inactive entries to active status. Entries to restore
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
    if !parser.HasArgs() {
        exit("you must specify at least one entry to restore")
    }
    filename, masterpass, db := loadDB(parser)

    // Grab the entries to restore.
    list := db.Inactive().FilterByIDString(parser.GetArgs()...)
    if len(list) == 0 {
        exit("no matching entries")
    }

    // Print a listing and request confirmation.
    printCompact(list, len(db.Inactive()), filepath.Base(filename))
    answer := input("  Restore the entries listed above? (y/n): ")
    if strings.ToLower(answer) == "y" {
        for _, entry := range list {
            db.SetActive(entry.Id)
        }
        saveDB(filename, masterpass, db)
        fmt.Println("  Entries restored.")
    } else {
        fmt.Println("  Restore aborted.")
    }
    printLineOfChar("─")
}
