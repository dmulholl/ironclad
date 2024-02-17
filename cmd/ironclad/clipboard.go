package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/dmulholl/ironclad/internal/config"
)

// Write a string to the system clipboard. Automatically overwrite the
// clipboard after a customizable delay.
func writeToClipboard(value string) {
	if clipboard.Unsupported {
		exit("clipboard functionality is not supported on this system")
	}

	// Default timeout.
	milliseconds := 15000

	// Check for a custom timeout.
	strval, found, err := config.Get("clipboard-timeout-seconds")
	if err != nil {
		exit(err)
	}
	if found {
		seconds, err := strconv.ParseInt(strval, 10, 32)
		if err != nil {
			exit(err)
		}
		milliseconds = int(seconds) * 1000
	}

	err = clipboard.WriteAll(value)
	if err != nil {
		exit(err)
	}

	fmt.Fprint(os.Stderr, "Clipboard: ")

	intervals := terminalWidth() - 20
	ms_per_interval := milliseconds / intervals
	strlen := intervals + 7

	for count := 0; count <= intervals; count++ {
		fmt.Fprintf(os.Stderr, "▌")
		fmt.Fprintf(os.Stderr, charstr(count, '█'))
		fmt.Fprintf(os.Stderr, charstr(intervals-count, ' '))
		fmt.Fprintf(os.Stderr, "▐")

		ms_remaining := milliseconds - count*ms_per_interval
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
