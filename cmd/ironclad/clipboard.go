package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/dmulholl/ironclad/internal/config"
	"github.com/dmulholl/ironclad/internal/textutils"
)

// Writes a string to the system clipboard.
func writeToClipboard(value string) error {
	if clipboard.Unsupported {
		return fmt.Errorf("clipboard functionality is not supported on this system")
	}

	err := clipboard.WriteAll(value)
	if err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	return nil
}

// Writes a string to the system clipboard. Automatically overwrites the clipboard after a
// customizable timeout read from the config file.
func writeToClipboardWithTimeout(value string) error {
	if clipboard.Unsupported {
		return fmt.Errorf("clipboard functionality is not supported on this system")
	}

	// Default timeout.
	timeoutMs := 15000

	// Check for a custom timeout.
	customTimeout, found, err := config.Get("clipboard-timeout-seconds")
	if err != nil {
		return fmt.Errorf("failed to read clipboard timeout from config file: %w", err)
	}

	if found {
		timeoutS, err := strconv.ParseInt(customTimeout, 10, 32)
		if err != nil {
			return fmt.Errorf("failed to parse config value for clipboard-timeout-seconds: %w", err)
		}
		timeoutMs = int(timeoutS) * 1000
	}

	err = clipboard.WriteAll(value)
	if err != nil {
		return fmt.Errorf("failed to write to clipboard: %w", err)
	}

	fmt.Fprint(os.Stderr, "Clipboard: ")

	intervals := terminalWidth() - 20
	ms_per_interval := timeoutMs / intervals
	strlen := intervals + 7

	for count := 0; count <= intervals; count++ {
		fmt.Fprint(os.Stderr, "▌")
		fmt.Fprint(os.Stderr, textutils.RuneString(count, '█'))
		fmt.Fprint(os.Stderr, textutils.RuneString(intervals-count, ' '))
		fmt.Fprint(os.Stderr, "▐")

		ms_remaining := timeoutMs - count*ms_per_interval
		fmt.Fprintf(os.Stderr, " %02.fs ", float64(ms_remaining)/1000)

		time.Sleep(time.Duration(ms_per_interval) * time.Millisecond)

		fmt.Fprint(os.Stderr, textutils.RuneString(strlen, '\b'))
		fmt.Fprint(os.Stderr, textutils.RuneString(strlen, ' '))
		fmt.Fprint(os.Stderr, textutils.RuneString(strlen, '\b'))
	}

	clipboardContent, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read from clipboard: %w", err)
	}

	if clipboardContent != value {
		fmt.Fprintln(os.Stderr, "[ALTERED]")
		return nil
	}

	err = clipboard.WriteAll("[CLIPBOARD OVERWRITTEN BY IRONCLAD]")
	if err != nil {
		return fmt.Errorf("failed to overwrite clipboard: %w", err)
	}

	fmt.Fprintln(os.Stderr, "[CLEARED]")
	return nil
}
