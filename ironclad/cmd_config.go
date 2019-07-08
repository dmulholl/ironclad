package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
    "strings"
)


import (
    "github.com/dmulholl/ironclad/ironconfig"
)


var configHelp = fmt.Sprintf(`
Usage: %s config [FLAGS] [ARGUMENTS]

  Set or display a configuration value.

  A single argument will be treated as a key and the associated value
  displayed. Two arguments will be treated as a key-value pair to be set.
  If no arguments are supplied, the config file itself will be printed.

  The following options are supported:

  clip-timeout              Clipboard timeout in seconds.
  timeout                   Password timeout in minutes.

Arguments:
  [key]                     Key to set or print.
  [value]                   Value to set.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerConfigCmd(parser *janus.ArgParser) {
    parser.NewCmd("config", configHelp, configCallback)
}


func configCallback(parser *janus.ArgParser) {
    if !parser.HasArgs() {
        if !ironconfig.FileExists() {
            exit("no config file exists")
        }
        content, err := ioutil.ReadFile(ironconfig.ConfigFile)
        if err != nil {
            exit(err)
        }
        fmt.Println(strings.TrimSpace(string(content)))
    } else if parser.NumArgs() == 1 {
        value, found, err := ironconfig.Get(parser.GetArg(0))
        if err != nil {
            exit(err)
        }
        if !found {
            exit("key not found")
        }
        fmt.Println(value)
    } else if parser.NumArgs() == 2 {
        err := ironconfig.Set(parser.GetArg(0), parser.GetArg(1))
        if err != nil {
            exit(err)
        }
    } else {
        exit("too many arguments")
    }
}
