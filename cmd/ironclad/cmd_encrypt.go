package main

import (
	"fmt"
	"os"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/ioutils"
)

var encryptCmdHelptext = `
Usage: ironclad encrypt <file>

  Encrypts a file using 256-bit AES encryption.

  This command encrypts an arbitrary file using the same 256-bit AES
  encryption that Ironclad uses for password databases. Note that the file
  is read into memory, encrypted, and written out to disk in a single
  operation, so the command is only suitable for encrypting files which fit
  comfortably into your system's RAM.

  Encrypted files can be decrypted using the 'decrypt' command.

Arguments:
  <file>                    File to encrypt.

Options:
  -o, --out                 Output filename. Defaults to adding '.encrypted'.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerEncryptCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("encrypt")
	cmdParser.Helptext = encryptCmdHelptext
	cmdParser.Callback = encryptCmdCallback
	cmdParser.NewStringOption("out o", "")
}

func encryptCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing filename argument")
	}

	inputfile := cmdParser.Args[0]
	outputfile := cmdParser.StringValue("out")
	if outputfile == "" {
		outputfile = inputfile + ".encrypted"
	}

	password, err := inputMasked("Password: ")
	if err != nil {
		return err
	}

	content, err := os.ReadFile(inputfile)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	err = ioutils.Save(outputfile, password, content)
	if err != nil {
		return fmt.Errorf("error saving file: %w", err)
	}

	return nil
}
