package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/ironclad/ironconfig"
    "github.com/dmulholland/clio/go/clio"
)


// Help text for the 'config' command.
var configHelptext = fmt.Sprintf(`
Usage: %s config [FLAGS] ARGUMENTS

  Set or print a configuration option.

Arguments:
  <key>                     Key to set or print.
  [value]                   Value to set.

Flags:
  --help                    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'config' command.
func configCallback(parser *clio.ArgParser) {

    if !parser.HasArgs() {
        exit("you must supply at least one argument")
    }

    if parser.LenArgs() == 1 {
        value, found, err := ironconfig.Get(parser.GetArg(0))
        if err != nil {
            exit(err)
        }
        if !found {
            exit("key not found")
        }
        fmt.Println(value)
    } else {
        err := ironconfig.Set(parser.GetArg(0), parser.GetArg(1))
        if err != nil {
            exit(err)
        }
    }
}
