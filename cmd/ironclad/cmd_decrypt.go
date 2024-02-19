package main

import (
	"fmt"
	"os"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/fileio"
)

var decryptCmdHelptext = `
Usage: ironclad decrypt <file>

  Decrypts a file encrypted using the 'encrypt' command. (This command can also
  be used to directly decrypt a password database.)

Arguments:
  <file>                    File to decrypt.

Options:
  -o, --out                 Output filename. Defaults to adding '.decrypted'.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerDecryptCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("decrypt")
	cmdParser.Helptext = decryptCmdHelptext
	cmdParser.Callback = decryptCmdCallback
	cmdParser.NewStringOption("out o", "")
}

func decryptCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing filename argument")
	}

	inputfile := cmdParser.Args[0]
	outputfile := cmdParser.StringValue("out")
	if outputfile == "" {
		outputfile = inputfile + ".decrypted"
	}

	password, err := inputMasked("Password: ")
	if err != nil {
		return err
	}

	content, err := fileio.Load(inputfile, password)
	if err != nil {
		return fmt.Errorf("failed to load file: %w", err)
	}

	err = os.WriteFile(outputfile, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
