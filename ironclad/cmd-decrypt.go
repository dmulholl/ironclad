package main


import "github.com/dmulholland/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
    "io/ioutil"
)


import (
    "github.com/dmulholland/ironclad/ironio"
)


var decryptHelp = fmt.Sprintf(`
Usage: %s decrypt [FLAGS] [OPTIONS] ARGUMENT

  Decrypt a file encrypted using the 'encrypt' command. (This command can also
  be used to directly decrypt a password database.)

Arguments:
  <file>                    File to decrypt.

Options:
  -o, --out                 Output filename. Defaults to adding '.decrypted'.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerDecrypt(parser *janus.ArgParser) {
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
