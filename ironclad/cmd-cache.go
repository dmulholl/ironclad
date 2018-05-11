package main


import "github.com/dmulholland/args"


import (
    "os"
    "os/signal"
    "fmt"
    "strconv"
    "time"
    "path/filepath"
)


import (
    "github.com/dmulholland/ironclad/ironrpc"
    "github.com/dmulholland/ironclad/ironconfig"
)


var cacheHelp = fmt.Sprintf(`
Usage: %s cache [FLAGS]

  Run the cached-password server. This comand is run automatically when
  required; it should not be run manually.

Flags:
  -h, --help    Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


func registerCache(parser *args.ArgParser) {
    parser.NewCmd("cache", cacheHelp, cacheCallback)
}


func cacheCallback(parser *args.ArgParser) {

    // Set up a handler to intercept SIGINT interrupts. This fixes an annoying
    // bug where hitting Ctrl-C to short-circuit a clipboard countdown could
    // kill the cache server, forcing the user to re-enter their password.
    //
    // > What's happening is that if you send a process SIGINT (as e.g.
    // > os.Interrupt does), all proceses in the same process group will also
    // > get that signal (which includes child processes) - SIGINT will by
    // > default terminate a process.
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    go func() {
        for _ = range c {
            // Do nothing.
        }
    }()

    // Check if a cache timeout has been set in the config file.
    timeout, found, err := ironconfig.Get("timeout")
    if err != nil {
        exit("cacheCallback:", err)
    }
    if found {
        minutes, err := strconv.ParseInt(timeout, 10, 64)
        if err != nil {
            exit("cacheCallback:", err)
        }
        if minutes == 0 {
            os.Exit(0)
        }
        ironrpc.CacheTimeout = time.Duration(minutes) * time.Minute
    }

    // Run the cache server.
    err = ironrpc.Serve()
    if err != nil {
        exit("cacheCallback:", err)
    }
}
