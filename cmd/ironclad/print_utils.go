package main

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

// Returns true if stdout is connected to a terminal.
func stdoutIsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Returns the width of the terminal window. Defaults to 80 if the width cannot be determined.
func terminalWidth() int {
	if !stdoutIsTerminal() {
		return 80
	}

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}

	return width
}

// Print a line of characters.
func printLineOfChar(char string) {
	length := terminalWidth()
	fmt.Print("\u001B[90m")
	fmt.Print(strings.Repeat(char, length))
	fmt.Println("\u001B[0m")
}

// Print an indented line of characters.
func printIndentedLineOfChar(char string) {
	length := terminalWidth() - 4
	fmt.Print("\u001B[90m")
	fmt.Print("  ")
	fmt.Print(strings.Repeat(char, length))
	fmt.Println("\u001B[0m")
}

// Print in grey.
func printGrey(format string, args ...interface{}) {
	fmt.Print("\u001B[90m")
	fmt.Printf(format, args...)
	fmt.Print("\u001B[0m")
}

// Print in grey with a newline.
func printlnGrey(format string, args ...interface{}) {
	fmt.Print("\u001B[90m")
	fmt.Printf(format, args...)
	fmt.Print("\u001B[0m\n")
}

// Print a heading.
func printHeading(text, meta string) {
	printLineOfChar("─")
	fmt.Print("  ")
	fmt.Print(text)
	numSpaces := terminalWidth() - len(text) - len(meta) - 4
	fmt.Print(strings.Repeat(" ", numSpaces))
	printlnGrey(meta)
	printLineOfChar("─")
}
