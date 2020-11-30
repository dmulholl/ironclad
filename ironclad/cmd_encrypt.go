package main


import "github.com/dmulholl/janus/v2"


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
)


import (
    "github.com/dmulholl/ironclad/ironio"
)


var encryptHelp = fmt.Sprintf(`
Usage: %s encrypt <file>

  Encrypt a file using 256-bit AES encryption.

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


func registerEncryptCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("encrypt", encryptHelp, encryptCallback)
    cmd.NewString("out o")
}


func encryptCallback(parser *janus.ArgParser) {
    if !parser.HasArgs() {
        exit("missing filename")
    }

    inputfile := parser.GetArg(0)
    outputfile := parser.GetString("out")
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
