package main


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
    "strings"
    "github.com/dmulholland/ironclad/ironconfig"
    "github.com/dmulholland/clio/go/clio"
)


// Help text for the 'config' command.
var configHelptext = fmt.Sprintf(`
Usage: %s config [FLAGS] ARGUMENTS

  Print or set the value of a configuration option.

  If a single argument is supplied, the value of that key will be printed. If
  two arguments are supplied, the first will be treated as a key and the
  second as a value to be set. If no arguments are supplied, the content of the
  configuration file itself will be printed.

Arguments:
  <key>                     Key to set or print.
  [value]                   Value to set.

Flags:
  --help                    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'config' command.
func configCallback(parser *clio.ArgParser) {

    if !parser.HasArgs() {
        content, err := ioutil.ReadFile(configfile)
        if err != nil {
            exit(err)
        }
        fmt.Println(strings.TrimSpace(string(content)))
    } else if parser.LenArgs() == 1 {
        value, found, err := ironconfig.Get(parser.GetArg(0))
        if err != nil {
            exit(err)
        }
        if !found {
            exit("key not found")
        }
        fmt.Println(value)
    } else if parser.LenArgs() == 2 {
        err := ironconfig.Set(parser.GetArg(0), parser.GetArg(1))
        if err != nil {
            exit(err)
        }
    } else {
        exit("too many arguments")
    }
}
