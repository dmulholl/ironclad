package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
)


import (
    "github.com/dmulholl/ironclad/irondb"
)


var initHelp = fmt.Sprintf(`
Usage: %s init [FLAGS] ARGUMENTS

  Create a new encrypted password database. You will be prompted to supply
  a master password.

Arguments:
  <file>                    Filename for the new database.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerInitCmd(parser *janus.ArgParser) {
    parser.NewCmd("init", initHelp, initCallback)
}


func initCallback(parser *janus.ArgParser) {
    if !parser.HasArgs() {
        exit("you must supply a filename for the database")
    }
    filename := parser.GetArgs()[0]

    masterpass1 := inputPass("Enter the master password for the new database: ")
    masterpass2 := inputPass("Please confirm the master password:             ")
    if masterpass1 != masterpass2 {
        exit("the passwords do not match")
    }

    cachepass1 := inputPass("Enter the cache password for the new database:  ")
    cachepass2 := inputPass("Please confirm the cache password:              ")
    if cachepass1 != cachepass2 {
        exit("the passwords do not match")
    }

    db := irondb.New(cachepass1)
    setCachedFilename(filename)
    saveDB(filename, masterpass1, db)
}
