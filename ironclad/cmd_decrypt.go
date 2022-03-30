package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dmulholl/ironclad/ironio"
	"github.com/dmulholl/janus/v2"
)

var decryptHelp = fmt.Sprintf(`
Usage: %s decrypt <file>

  Decrypt a file encrypted using the 'encrypt' command. (This command can also
  be used to directly decrypt a password database.)

Arguments:
  <file>                    File to decrypt.

Options:
  -o, --out                 Output filename. Defaults to adding '.decrypted'.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))

func registerDecryptCmd(parser *janus.ArgParser) {
	cmd := parser.NewCmd("decrypt", decryptHelp, decryptCallback)
	cmd.NewString("out o")
}

func decryptCallback(parser *janus.ArgParser) {
	if !parser.HasArgs() {
		exit("missing filename")
	}

	inputfile := parser.GetArg(0)
	outputfile := parser.GetString("out")
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
