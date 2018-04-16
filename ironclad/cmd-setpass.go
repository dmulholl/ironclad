package main


import "github.com/dmulholland/args"


import (
    "fmt"
    "os"
    "path/filepath"
)


var setpassHelp = fmt.Sprintf(`
Usage: %s setpass [FLAGS] [OPTIONS]

  Change a database's master password.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func setpassCallback(parser *args.ArgParser) {
    filename, _, db := loadDB(parser)

    line("─")
    password := inputPass("Enter new password:   ")
    confirm  := inputPass("Confirm new password: ")
    line("─")

    if password == confirm {
        saveDB(filename, password, db)
    } else {
        exit("passwords do not match")
    }
}
