package main


import "github.com/atotto/clipboard"


import (
    "os"
    "fmt"
    "time"
)


// Write a string to the system clipboard. Automatically overwrite the
// clipboard after a fifteen second delay.
func writeToClipboard(value string) {

    if clipboard.Unsupported {
        exit("clipboard functionality is not supported on this system")
    }

    err := clipboard.WriteAll(value)
    if err != nil {
        exit(err)
    }

    fmt.Fprint(os.Stderr, "Clipboard: ")

    ms_total := 15000
    intervals := terminalWidth() - 20
    ms_per_interval := ms_total / intervals
    strlen := intervals + 7

    for count := 0; count <= intervals; count++ {

        fmt.Fprintf(os.Stderr, "|")
        fmt.Fprintf(os.Stderr, charstr(count, '='))
        fmt.Fprintf(os.Stderr, charstr(intervals - count, ' '))
        fmt.Fprintf(os.Stderr, "|")

        ms_remaining := ms_total - count * ms_per_interval
        fmt.Fprintf(os.Stderr, " %02.fs ", float64(ms_remaining)/1000)

        time.Sleep(time.Duration(ms_per_interval) * time.Millisecond)

        fmt.Fprintf(os.Stderr, charstr(strlen, '\b'))
        fmt.Fprintf(os.Stderr, charstr(strlen, ' '))
        fmt.Fprintf(os.Stderr, charstr(strlen, '\b'))
    }

    clipContent, err := clipboard.ReadAll()
    if err != nil {
        exit(err)
    }

    if clipContent == value {
        err = clipboard.WriteAll("[Clipboard overwritten by Ironclad]")
        if err != nil {
            exit(err)
        }
        fmt.Fprintln(os.Stderr, "[CLEARED]")
    } else {
        fmt.Fprintln(os.Stderr, "[ALTERED]")
    }
}
