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
)


// Cache the last-used filename for the application's next run.
func cacheLastFilename(filename string) {
    filename, err := filepath.Abs(filename)
    if err != nil {
        exit("Error:", err)
    }
    err = ironconfig.Set("file", filename)
    if err != nil {
        exit("Error:", err)
    }
}


// Fetch the last-used filename from the application's last run.
func fetchLastFilename() (string, bool) {
    filename, found, err := ironconfig.Get("file")
    if err != nil {
        exit("Error:", err)
    }
    return filename, found
}


// Temporarily cache the last-used password for the application's next run.
func cacheLastPassword(password string) {

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
            errmsg("Error connecting to password server: ", err)
            return
        }
    }

    // Generate a random 32-byte token.
    bytes, err := ironcrypt.RandBytes(32)
    if err != nil {
        exit("Error:", err)
    }
    token := base64.StdEncoding.EncodeToString(bytes)

    // Save the token in the application's config file.
    err = ironconfig.Set("token", token)
    if err != nil {
        exit("Error:", err)
    }

    // Cache the token and password pair.
    _, err = client.Set(token, password)
    if err != nil {
        errmsg("Error caching password: ", err)
    }
}


// Fetch the last-used password from the application's last run.
func fetchLastPassword() (password string, found bool) {

    // Fetch the password token from the application's config file.
    token, found, err := ironconfig.Get("token")
    if err != nil {
        exit("Error:", err)
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
        errmsg("Error fetching password: ", err)
        return "", false
    }

    return password, true
}
