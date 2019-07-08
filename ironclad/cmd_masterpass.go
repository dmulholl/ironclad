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
    password := inputPass("Enter new password:   ")
    confirm  := inputPass("Confirm new password: ")
    printLineOfChar("─")

    if password == confirm {
        saveDB(filename, password, db)
    } else {
        exit("passwords do not match")
    }
}
