package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/ironconfig"
)

var configCmdHelptext = `
Usage: ironclad config
       ironclad config <key>
       ironclad config <key> <value>

  Sets or displays a configuration value.

  - If a single argument is supplied, it will be treated as a key and the
    associated value will be printed.
  - If two arguments are supplied, they will be treated as a key-value pair
    to be set.
  - If no arguments are supplied, the config file itself will be printed.

  The following configuration options are supported:

  cache-timeout-minutes         Master-password cache timeout in minutes.
  clipboard-timeout-seconds     Clipboard timeout in seconds.

Arguments:
  [key]                         Key to set or print.
  [value]                       Value to set.

Flags:
  -h, --help                    Print this command's help text and exit.
`

func registerConfigCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("config")
	cmdParser.Helptext = configCmdHelptext
	cmdParser.Callback = configCmdCallback
}

func configCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		if !ironconfig.FileExists() {
			return fmt.Errorf("config file not found")
		}

		content, err := os.ReadFile(ironconfig.ConfigFile)
		if err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}

		fmt.Println(strings.TrimSpace(string(content)))
		return nil
	}

	if len(cmdParser.Args) == 1 {
		value, found, err := ironconfig.Get(cmdParser.Args[0])
		if err != nil {
			return fmt.Errorf("error reading config file: %w", err)
		}

		if !found {
			return fmt.Errorf("key not found")
		}

		fmt.Println(value)
		return nil
	}

	if len(cmdParser.Args) == 2 {
		err := ironconfig.Set(cmdParser.Args[0], cmdParser.Args[1])
		if err != nil {
			return fmt.Errorf("error setting config value: %w", err)
		}

		return nil
	}

	return fmt.Errorf("too many arguments")
}
