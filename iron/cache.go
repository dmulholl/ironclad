package main


import (
    "github.com/dmulholland/ironclad/ironconfig"
    "github.com/dmulholland/ironclad/ironcrypt"
    "github.com/dmulholland/ironclad/ironrpc"
    "encoding/base64"
    "time"
    "os/exec"
    "path/filepath"
    "os"
    "fmt"
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
func setCachedPassword(password string) {

    // Attempt to make a connection to the password server.
    // If we can't make a connection, launch a new server.
    client, err := ironrpc.NewClient(ironaddress)
    if err != nil {
        cmd := exec.Command(os.Args[0], "cache")
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Start()

        // Give the new sever time to warm up.
        time.Sleep(time.Millisecond * 100)

        // Try making a connection again.
        client, err = ironrpc.NewClient(ironaddress)
        if err != nil {
            fmt.Fprintf(os.Stderr,
                "Error connecting to password server: %v\n", err)
            return
        }
    }

    // Generate a random 32-byte token.
    bytes, err := ironcrypt.RandBytes(32)
    if err != nil {
        exit(err)
    }
    token := base64.StdEncoding.EncodeToString(bytes)

    // Save the token in the application's config file.
    err = ironconfig.Set("token", token)
    if err != nil {
        exit(err)
    }

    // Cache the token and password pair.
    _, err = client.Set(token, password)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error caching password: %v\n", err)
    }
}


// Fetch the cached password (if it exists) from the application's last run.
func getCachedPassword() (password string, found bool) {

    // Fetch the password token from the application's config file.
    token, found, err := ironconfig.Get("token")
    if err != nil {
        exit(err)
    }
    if !found {
        return "", false
    }

    // Attempt to make a connection to the password server.
    client, err := ironrpc.NewClient(ironaddress)
    if err != nil {
        return "", false
    }

    // Retrieve the password from the server.
    password, err = client.Get(token)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error fetching password: %v\n", err)
        return "", false
    }

    return password, true
}
