package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
)


import (
    "github.com/dmulholland/ironclad/ironio"
)


// Help text for the 'decrypt' command.
var decryptHelp = fmt.Sprintf(`
Usage: %s decrypt [FLAGS] [OPTIONS] [ARGUMENTS]

  Decrypt a file.

Arguments:
  <file>                    File to decrypt.

Options:
  -o, --out                 Output filename. Defaults to adding '.decrypted'.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'export' command.
func decryptCallback(parser *clio.ArgParser) {

    if !parser.HasArgs() {
        exit("missing filename")
    }

    inputfile := parser.GetArg(0)

    outputfile := parser.GetStr("out")
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
