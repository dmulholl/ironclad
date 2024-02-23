package main

import (
	"fmt"
	"io"
	"os"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/ioutils"
)

var decryptCmdHelptext = `
Usage: ironclad decrypt

  Decrypts data encrypted using the 'encrypt' command.

  If no -i/--input file is specified, input will be read from the standard
  input stream. In this case, the password must be supplied via the
  -p/--password option or an $IRONCLAD_PASSWORD environment variable.

  if no -o/--output file is specified, output will be written to the standard
  output stream.

  You can use this command to decrypt a password database. The output will be
  the database's raw JSON data store.

Options:
  -i, --input <file>        Input filename.
  -o, --output <file>       Output filename.
  -p, --password <str>      Password for decrypting the file.

Flags:
  -h, --help                Print this command's help text and exit.
`

func registerDecryptCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("decrypt")
	cmdParser.Helptext = decryptCmdHelptext
	cmdParser.Callback = decryptCmdCallback
	cmdParser.NewStringOption("input i", "")
	cmdParser.NewStringOption("output o", "")
	cmdParser.NewStringOption("password p", "")
}

func decryptCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
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
		decrypted, err := ioutils.Load(cmdParser.StringValue("input"), password)
		if err != nil {
			return fmt.Errorf("error loading file: %w", err)
		}

		plaintext = decrypted
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading stdin: %w", err)
		}

		decrypted, err := ioutils.Decrypt(password, data)
		if err != nil {
			return err
		}

		plaintext = decrypted
	}

	if cmdParser.Found("output") {
		err := os.WriteFile(cmdParser.StringValue("output"), plaintext, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		return nil
	}

	fmt.Print(string(plaintext))
	return nil
}
