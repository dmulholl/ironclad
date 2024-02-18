package main

import (
	"os"

	"golang.org/x/term"
)

// Returns true if stdout is connected to a terminal.
func stdoutIsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Returns the width of the terminal window. Defaults to 80 if the width
// cannot be determined.
func terminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		return width
	}
	return 80
}
