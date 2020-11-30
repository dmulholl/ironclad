package main


import "github.com/dmulholl/janus/v2"


import (
    "fmt"
    "os"
    "path/filepath"
)


var masterpassHelp = fmt.Sprintf(`
Usage: %s setmasterpass

  Change a database's master password.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerSetMasterPassCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("setmasterpass", masterpassHelp, masterpassCallback)
    cmd.NewString("file f")
}


func masterpassCallback(parser *janus.ArgParser) {
    filename, _, db := loadDB(parser)

    printLineOfChar("─")
    newMasterPass        := inputPass("Enter new master password: ")
    confirmNewMasterPass := inputPass("      Re-enter to confirm: ")
    printLineOfChar("─")

    if newMasterPass == confirmNewMasterPass {
        saveDB(filename, newMasterPass, db)
        setCachedPassword(filename, newMasterPass, db.CachePass)
    } else {
        exit("passwords do not match")
    }
}
