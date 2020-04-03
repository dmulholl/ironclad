package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
)


var importHelp = fmt.Sprintf(`
Usage: %s import [FLAGS] [OPTIONS] [ARGUMENT]

  Import a list of entries.

Arguments:
  <file>                    File to import.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerImportCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("import", importHelp, importCallback)
    cmd.NewString("file f")
}


func importCallback(parser *janus.ArgParser) {
    var input []byte
    var err error
    if parser.HasArgs() {
        input, err = ioutil.ReadFile(parser.GetArg(0))
        if err != nil {
            exit(err)
        }
    } else {
        exit("you must specify a file to import")
    }

    filename, masterpass, db := loadDB(parser)
    err = db.Import(input)
    if err != nil {
        exit(err)
    }
    saveDB(filename, masterpass, db)
}
