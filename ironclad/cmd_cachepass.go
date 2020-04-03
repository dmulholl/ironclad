package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
)


var cachepassHelp = fmt.Sprintf(`
Usage: %s cachepass [FLAGS] [OPTIONS]

  Change a database's cache password.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerCachepassCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd(
        "cachepass", cachepassHelp, cachepassCallback)
    cmd.NewString("file f")
}


func cachepassCallback(parser *janus.ArgParser) {
    filename, masterpass, db := loadDB(parser)

    printLineOfChar("─")
    newCachePass        := inputPass("Enter new cache password:   ")
    confirmNewCachePass := inputPass("Confirm new cache password: ")
    printLineOfChar("─")

    if newCachePass == confirmNewCachePass {
        db.CachePass = newCachePass
        saveDB(filename, masterpass, db)
        setCachedPassword(filename, masterpass, newCachePass)
    } else {
        exit("passwords do not match")
    }
}
