package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dmulholl/argo"
	"github.com/dmulholl/ironclad/ironio"
)

var encryptHelp = fmt.Sprintf(`
Usage: %s encrypt <file>

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
`, filepath.Base(os.Args[0]))

func registerEncryptCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("encrypt")
	cmdParser.Helptext = encryptHelp
	cmdParser.Callback = encryptCallback
	cmdParser.NewStringOption("out o", "")
}

func encryptCallback(cmdName string, cmdParser *argo.ArgParser) {
	if !cmdParser.HasArgs() {
		exit("missing filename")
	}

	inputfile := cmdParser.Arg(0)
	outputfile := cmdParser.StringValue("out")
	if outputfile == "" {
		outputfile = inputfile + ".encrypted"
	}

	password := inputPass("Password: ")
	content, err := ioutil.ReadFile(inputfile)
	if err != nil {
		exit(err)
	}

	err = ironio.Save(outputfile, password, content)
	if err != nil {
		exit(err)
	}
}
