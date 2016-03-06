package main


import (
    "os"
    "fmt"
    "strconv"
    "time"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/ironrpc"
    "github.com/dmulholland/ironclad/ironconfig"
)


// Help text for the 'cache' command.
var cacheHelptext = fmt.Sprintf(`
Usage: %s cache [FLAGS]

  Run the cached-password server. This command should not be run manually.

Flags:
  --help                    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'serve' command.
func cacheCallback(parser *clio.ArgParser) {

    // Check if a cache timeout has been set in the config file.
    timeout, found, err := ironconfig.Get("timeout")
    if err != nil {
        exit("Error:", err)
    }
    if found {
        minutes, err := strconv.ParseInt(timeout, 10, 64)
        if err != nil {
            exit("Error:", err)
        }
        ironrpc.ServerTimeout = time.Duration(minutes) * time.Minute
    }

    // Run the cache server.
    err = ironrpc.Serve(ironaddress)
    if err != nil {
        exit("Error:", err)
    }
}
