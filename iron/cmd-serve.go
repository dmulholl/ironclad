package main


import (
    "os"
    "fmt"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/ironrpc"
)


// Help text for the 'serve' command.
var serveHelptext = fmt.Sprintf(`
Usage: %s serve [FLAGS]

  Run the cached-password server. This command should not be run manually.

Flags:
  --help                    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'serve' command.
func serveCallback(parser *clio.ArgParser) {
    err := ironrpc.Serve(ironaddress)
    if err != nil {
        exit("Error:", err)
    }
}
