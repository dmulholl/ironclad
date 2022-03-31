package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo"
	"github.com/dmulholl/ironclad/ironconfig"
)

var configHelp = fmt.Sprintf(`
Usage: %s config [key] [value]

  Sets or displays a configuration value.

  A single argument will be treated as a key and the associated value
  displayed. Two arguments will be treated as a key-value pair to be set.
  If no arguments are supplied, the config file itself will be printed.

  The following options are supported:

  cache-timeout-minutes         Master-password cache timeout in minutes.
  clipboard-timeout-seconds     Clipboard timeout in seconds.

Arguments:
  [key]                         Key to set or print.
  [value]                       Value to set.

Flags:
  -h, --help                    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerConfigCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("config")
	cmdParser.Helptext = configHelp
	cmdParser.Callback = configCallback
}

func configCallback(cmdName string, cmdParser *argo.ArgParser) {
	if !cmdParser.HasArgs() {
		if !ironconfig.FileExists() {
			exit("no config file exists")
		}
		content, err := ioutil.ReadFile(ironconfig.ConfigFile)
		if err != nil {
			exit(err)
		}
		fmt.Println(strings.TrimSpace(string(content)))
	} else if cmdParser.CountArgs() == 1 {
		value, found, err := ironconfig.Get(cmdParser.Args[0])
		if err != nil {
			exit(err)
		}
		if !found {
			exit("key not found")
		}
		fmt.Println(value)
	} else if cmdParser.CountArgs() == 2 {
		err := ironconfig.Set(cmdParser.Args[0], cmdParser.Args[1])
		if err != nil {
			exit(err)
		}
	} else {
		exit("too many arguments")
	}
}
