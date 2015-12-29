package main


import (
    "github.com/dmulholland/ironclad/ironconfig"
)


// Cache the last-used filename for the application's next run.
func cacheLastFilename(filename string) {
    err := ironconfig.Set("file", filename)
    if err != nil {
        exit("Error:", err)
    }
}


// Fetch the last-used filename from the application's last run.
func fetchLastFilename() (string, bool) {
    found, filename, err := ironconfig.Get("file")
    if err != nil {
        exit("Error:", err)
    }
    return filename, found
}


// Temporarily cache the last-used password for the application's next run.
func cacheLastPassword(password string) {}


// Fetch the last-used password from the application's last run.
func fetchLastPassword() (string, bool) {
    return "", false
}
