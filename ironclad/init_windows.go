package main

import (
	"golang.org/x/sys/windows"
	"os"
)

// Enable support for ANSI escape sequences in the Windows console.
// Ref: https://stackoverflow.com/questions/39627348
func init() {
	var originalMode uint32
	stdout := windows.Handle(os.Stdout.Fd())
	windows.GetConsoleMode(stdout, &originalMode)
	ansiMode := originalMode | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
	windows.SetConsoleMode(stdout, ansiMode)
}
