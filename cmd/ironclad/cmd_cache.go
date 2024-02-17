package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/ironconfig"
	"github.com/dmulholl/ironclad/internal/ironrpc"
)

var cacheCmdHelptext = `
Usage: ironclad cache

  Runs the cached-password server. This comand is run automatically when
  required; it should not be run manually.

Flags:
  -h, --help    Print this command's help text and exit.
`

func registerCacheCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("cache")
	cmdParser.Helptext = cacheCmdHelptext
	cmdParser.Callback = cacheCmdCallback
}

func cacheCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
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
	timeout, found, err := ironconfig.Get("cache-timeout-minutes")
	if err != nil {
		return fmt.Errorf("ironclad cache server: failed to get timeout: %w", err)
	}

	if found {
		numMinutes, err := strconv.ParseInt(timeout, 10, 64)
		if err != nil {
			return fmt.Errorf("ironclad cache server: failed to parse timeout: %w", err)
		}
		if numMinutes == 0 {
			return nil
		}
		ironrpc.CacheTimeout = time.Duration(numMinutes) * time.Minute
	}

	// Run the cache server.
	err = ironrpc.Serve()
	if err != nil {
		return fmt.Errorf("ironclad cache server: %w", err)
	}

	return nil
}
