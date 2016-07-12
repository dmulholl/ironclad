package main


import (
    "os"
    "fmt"
    "time"
    "github.com/atotto/clipboard"
)


// Write a string to the system clipboard. Automatically overwrite after ten
// seconds.
func writeToClipboard(value string) {

    if clipboard.Unsupported {
        exit("Error: clipboard functionality is not supported on this system.")
    }

    err := clipboard.WriteAll(value)
    if err != nil {
        exit("Error:", err)
    }

    fmt.Fprint(os.Stderr, "Copied to clipboard. Clearing in:   ")
    for i := 10; i > 0; i-- {
        fmt.Fprintf(os.Stderr, "\b\b%2v", i)
        time.Sleep(time.Second)
    }

    clipContent, err := clipboard.ReadAll()
    if err != nil {
        exit("Error:", err)
    }

    if clipContent == value {
        err = clipboard.WriteAll("[Clipboard overwritten by Ironclad]")
        if err != nil {
            exit("Error:", err)
        }
        fmt.Fprintln(os.Stderr, "\b\b[CLEARED]")
    } else {
        fmt.Fprintln(os.Stderr, "\b\b[ALTERED]")
    }
}
