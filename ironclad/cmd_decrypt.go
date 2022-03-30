package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
	"github.com/dmulholl/ironclad/ironio"
)

var decryptHelp = fmt.Sprintf(`
Usage: %s decrypt <file>

  Decrypts a file encrypted using the 'encrypt' command. (This command can also
  be used to directly decrypt a password database.)

Arguments:
  <file>                    File to decrypt.

Options:
  -o, --out                 Output filename. Defaults to adding '.decrypted'.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerDecryptCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("decrypt")
	cmdParser.Helptext = decryptHelp
	cmdParser.Callback = decryptCallback
	cmdParser.NewStringOption("out o", "")
}

func decryptCallback(cmdName string, cmdParser *argo.ArgParser) {
	if !cmdParser.HasArgs() {
		exit("missing filename")
	}

	inputfile := cmdParser.Arg(0)
	outputfile := cmdParser.StringValue("out")
	if outputfile == "" {
		outputfile = inputfile + ".decrypted"
	}

	password := inputPass("Password: ")
	content, err := ironio.Load(inputfile, password)
	if err != nil {
		exit(err)
	}

	err = ioutil.WriteFile(outputfile, content, 0644)
	if err != nil {
		exit(err)
	}
}
