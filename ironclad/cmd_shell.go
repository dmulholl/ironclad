package main


import "github.com/dmulholl/janus-go/janus"


import (
    "fmt"
    "os"
    "path/filepath"
    "time"
    // "strings"
)

import (
    "github.com/dmulholl/ironclad/ironcrypt"
)


var shellHelp = fmt.Sprintf(`
Usage: %s shell [FLAGS] [OPTIONS]

  Run the Ironclad shell.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -h, --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerShellCmd(parser *janus.ArgParser) {
    cmd := parser.NewCmd("shell", shellHelp, shellCallback)
    cmd.NewString("file f")
}


func shellCallback(parser *janus.ArgParser) {
    // filename, masterpass, db := loadDB(parser)

    // for {
    //     cmdstr := input(">>> ")
    //     tokens := tokenize(cmdstr)
    //     if len(tokens) == 0 {
    //         continue
    //     }
    //     if tokens[0] == "exit" || tokens[0] == "quit" {
    //         break
    //     }
    // }

    //saveDB(filename, masterpass, db)


    // Generate a random salt.
    salt, err := ironcrypt.RandBytes(32)
    if err != nil {
        return
    }

    start := time.Now()

    // Use the password and salt to generate a file encryption key.
    ironcrypt.Key("this is a password", salt, 100000, 32)

    elapsed := time.Since(start)

    fmt.Println(elapsed)



}
