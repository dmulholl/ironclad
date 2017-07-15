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


var encryptHelp = fmt.Sprintf(`
Usage: %s encrypt [FLAGS] [OPTIONS] [ARGUMENTS]

  Encrypt a file using 256-bit AES encryption.

Arguments:
  <file>                    File to encrypt.

Options:
  -o, --out                 Output filename. Defaults to adding '.encrypted'.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func encryptCallback(parser *clio.ArgParser) {

    if !parser.HasArgs() {
        exit("missing filename")
    }

    inputfile := parser.GetArg(0)

    outputfile := parser.GetStr("out")
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
