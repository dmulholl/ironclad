package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/ioutils"
)

var encryptCmdHelptext = `
Usage: ironclad encrypt

  Encrypts data using 256-bit AES encryption.

  If no -i/--input file is specified, input will be read from the standard
  input stream. In this case, the password must be supplied via the
  -p/--password option or an $IRONCLAD_PASSWORD environment variable.

  if no -o/--output file is specified, output will be written to the standard
  output stream.

  Encrypted files can be decrypted using the 'decrypt' command.

Options:
  -i, --input <file>        Input filename.
  -o, --output <file>       Output filename.
  -p, --password <str>      Password for encrypting the file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerEncryptCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("encrypt")
	cmdParser.Helptext = encryptCmdHelptext
	cmdParser.Callback = encryptCmdCallback
	cmdParser.NewStringOption("input i", "")
	cmdParser.NewStringOption("output o", "")
	cmdParser.NewStringOption("password p", "")
}

func encryptCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	var password string

	if cmdParser.Found("password") {
		password = cmdParser.StringValue("password")
	} else if value, found := os.LookupEnv("IRONCLAD_PASSWORD"); found {
		password = value
	} else if cmdParser.Found("input") {
		value, err := inputMasked("Password: ")
		if err != nil {
			return err
		}
		password = value
	} else {
		return fmt.Errorf("missing password")
	}

	var plaintext []byte

	if cmdParser.Found("input") {
		data, err := os.ReadFile(cmdParser.StringValue("input"))
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}
		plaintext = data
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading stdin: %w", err)
		}
		plaintext = data
	}

	if cmdParser.Found("output") {
		err := ioutils.Save(cmdParser.StringValue("output"), password, plaintext)
		if err != nil {
			return fmt.Errorf("error saving file: %w", err)
		}
		return nil
	}

	encrypted, err := ioutils.Encrypt(password, plaintext)
	if err != nil {
		return err
	}

	fmt.Print(string(encrypted))
	return nil
}
