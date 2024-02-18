package main

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// Inserts the prefix string at the beginning of each non-empty line.
func indent(text, prefix string) string {
	var builder strings.Builder
	isBeginningOfLine := true

	for _, r := range text {
		if isBeginningOfLine && r != '\n' {
			builder.WriteString(prefix)
		}
		builder.WriteRune(r)
		isBeginningOfLine = (r == '\n')
	}

	return builder.String()
}

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
