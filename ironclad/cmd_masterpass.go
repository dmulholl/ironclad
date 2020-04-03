package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
)


var masterpassHelp = fmt.Sprintf(`
Usage: %s masterpass [FLAGS] [OPTIONS]

  Change a database's master password.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerMasterpassCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd(
        "masterpass setpass", masterpassHelp, masterpassCallback)
    cmd.NewString("file f")
}


func masterpassCallback(parser *janus.ArgParser) {
    filename, _, db := loadDB(parser)

    printLineOfChar("─")
    newMasterPass        := inputPass("Enter new master password:   ")
    confirmNewMasterPass := inputPass("Confirm new master password: ")
    printLineOfChar("─")

    if newMasterPass == confirmNewMasterPass {
        saveDB(filename, newMasterPass, db)
        setCachedPassword(filename, newMasterPass, db.CachePass)
    } else {
        exit("passwords do not match")
    }
}
