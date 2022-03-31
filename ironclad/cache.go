package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/dmulholl/ironclad/ironconfig"
	"github.com/dmulholl/ironclad/ironrpc"
)

// Cache the current filename for the application's next run.
func setCachedFilename(filename string) {
	filename, err := filepath.Abs(filename)
	if err != nil {
		exit(err)
	}
	err = ironconfig.Set("file", filename)
	if err != nil {
		exit(err)
	}
}

// Fetch the cached filename (if it exists) from the application's last run.
func getCachedFilename() (filename string, found bool) {
	filename, found, err := ironconfig.Get("file")
	if err != nil {
		exit(err)
	}
	return filename, found
}

// Cache the database password for the application's next run.
func setCachedPassword(filename, masterpass, cachepass string) {

	// If the cache timeout has been set to 0, do nothing.
	timeout, found, _ := ironconfig.Get("cache-timeout-minutes")
	if found && timeout == "0" {
		return
	}

	// Attempt to connect to the cache server. If we can't make a connection,
	// launch a new server.
	client, err := ironrpc.NewClient()
	if err != nil {
		cmd := exec.Command(os.Args[0], "cache")
		cmd.Stderr = os.Stderr
		cmd.Start()

		// Give the new server time to warm up and try again.
		time.Sleep(time.Millisecond * 100)
		client, err = ironrpc.NewClient()
		if err != nil {
			return
		}
	}
	defer client.Close()

	client.SetPass(filename, masterpass, cachepass)
}

// Fetch the cached master password (if it exists) from the application's last run.
func getCachedPassword(filename string) (masterpass string, success bool) {

	// Attempt to make a connection to the cache server.
	client, err := ironrpc.NewClient()
	if err != nil {
		return "", false
	}
	defer client.Close()

	// Check if the server has a cache entry for the database file.
	if !client.IsCached(filename) {
		return "", false
	}

	// Attempt to retrieve the master password from the server.
	// Try first using an empty string as the cache password.
	masterpass, err = client.GetPass(filename, "")
	if err == nil {
		return masterpass, true
	}

	// The empty string didn't work. Ask the user for their cache password
	// and try again.
	cachepass := inputPass("Cache Password: ")
	masterpass, err = client.GetPass(filename, cachepass)
	if err != nil {
		println("Error: invalid cache password.")
		return "", false
	}

	return masterpass, true
}
