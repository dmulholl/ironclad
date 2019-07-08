package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "strings"
    "path/filepath"
)


var purgeHelp = fmt.Sprintf(`
Usage: %s purge [FLAGS] [OPTIONS]

  Purge inactive entries from a database.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerPurgeCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("purge", purgeHelp, purgeCallback)
    cmd.NewString("file f")
}


func purgeCallback(parser *janus.ArgParser) {
    filePath, password, db := loadDB(parser)
    fileName := filepath.Base(filePath)

    list := db.Inactive()
    if len(list) == 0 {
        exit("no inactive entries to purge")
    }

    printCompact(list, len(list), fileName)
    answer := input("  Purge the entries listed above? (y/n): ")
    if strings.ToLower(answer) == "y" {
        db.PurgeInactive()
        fmt.Println("  Entries purged.")
        printLineOfChar("─")
    } else {
        fmt.Println("  Purge aborted.")
        printLineOfChar("─")
    }

    saveDB(filePath, password, db)
}
