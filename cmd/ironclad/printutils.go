package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// Returns true if stdout is connected to a terminal.
func stdoutIsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// Returns the width of the terminal window. Defaults to 80 if the width cannot be determined.
func terminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		return width
	}
	return 80
}

// Print a line of characters.
func printLineOfChar(char string) {
	length := terminalWidth()
	fmt.Print("\u001B[90m")
	for i := 0; i < length; i++ {
		fmt.Print(char)
	}
	fmt.Println("\u001B[0m")
}

// Print an indented line of characters.
func printIndentedLineOfChar(char string) {
	length := terminalWidth() - 4
	fmt.Print("\u001B[90m")
	fmt.Print("  ")
	for i := 0; i < length; i++ {
		fmt.Print(char)
	}
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
	for i := 0; i < numSpaces; i += 1 {
		fmt.Print(" ")
	}
	printlnGrey(meta)
	printLineOfChar("─")
}
