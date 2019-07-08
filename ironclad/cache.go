package main


import (
    "time"
    "os/exec"
    "os"
    "path/filepath"
)


import (
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
func setCachedPassword(filename, password string) {

    // If the cache timeout has been set to 0, do nothing.
    timeout, found, _ := ironconfig.Get("timeout")
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

    // Cache the token and password pair.
    client.SetPass(filename, password)
}


// Fetch the cached password (if it exists) from the application's last run.
func getCachedPassword(filename string) (password string, found bool) {

    // Fetch the authentication nonce.
    nonce, found, err := ironconfig.Get("nonce")
    if err != nil {
        exit(err)
    }
    if !found {
        return "", false
    }

    // Attempt to make a connection to the password server.
    client, err := ironrpc.NewClient()
    if err != nil {
        return "", false
    }
    defer client.Close()

    // Retrieve the password from the server.
    password, err = client.GetPass(filename, nonce)
    if err != nil {
        return "", false
    }

    return password, true
}
