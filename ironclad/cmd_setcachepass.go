package main


import "github.com/dmulholl/janus/v2"


import (
    "fmt"
    "os"
    "path/filepath"
)


var cachepassHelp = fmt.Sprintf(`
Usage: %s setcachepass

  Change a database's cache password. This password is used to encrypt the
  master password while it's temporarily cached in memory.

  Note that if you set the cache password to an empty string you will not be
  prompted to enter it.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerSetCachePassCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("setcachepass", cachepassHelp, cachepassCallback)
    cmd.NewString("file f")
}


func cachepassCallback(parser *janus.ArgParser) {
    filename, masterpass, db := loadDB(parser)

    printLineOfChar("─")
    newCachePass        := inputPass("Enter new cache password: ")
    confirmNewCachePass := inputPass("     Re-enter to confirm: ")
    printLineOfChar("─")

    if newCachePass == confirmNewCachePass {
        db.CachePass = newCachePass
        saveDB(filename, masterpass, db)
        setCachedPassword(filename, masterpass, newCachePass)
    } else {
        exit("passwords do not match")
    }
}
