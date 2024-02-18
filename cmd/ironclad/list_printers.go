package main

import (
	"fmt"
	"strings"

	"github.com/dmulholl/ironclad/internal/database"
	"github.com/dmulholl/ironclad/internal/textutils"
	"github.com/mitchellh/go-wordwrap"
)

// Print a list of entries in compact format.
func printCompactList(list database.EntryList, poolSize int, filename string) {
	if len(list) == 0 {
		printHeading("No Entries", filename)
		return
	}

	// Header.
	printLineOfChar("─")
	fmt.Printf("  ID")
	printGrey("  ·  ")
	fmt.Printf("TITLE")
	numSpaces := terminalWidth() - len(filename) - 16
	for i := 0; i < numSpaces; i += 1 {
		fmt.Printf(" ")
	}
	printlnGrey(filename)
	printLineOfChar("─")

	// Print the entry listing.
	for _, entry := range list {
		fmt.Printf("%4d", entry.Id)
		printGrey("  ·  ")
		fmt.Printf("%s\n", entry.Title)
	}

	// Footer.
	printLineOfChar("─")
	fmt.Printf("  %d/%d Entries\n", len(list), poolSize)
	printLineOfChar("─")
}

// Print a list of entries in verbose format.
func printVerboseList(list database.EntryList, poolSize int, showPassword, showNotes bool, title, filename string) {
	if len(list) == 0 {
		printHeading("No Entries", filename)
		return
	}

	// Header.
	printHeading(title, filename)

	// Print the entry listing.
	for _, entry := range list {
		fmt.Printf("  ID:       %d\n", entry.Id)
		fmt.Printf("  Title:    %s\n", entry.Title)

		if entry.Url != "" {
			fmt.Printf("  URL:      %s\n", entry.Url)
		}

		if entry.Username != "" {
			fmt.Printf("  Username: %s\n", entry.Username)
		}

		if entry.GetPassword() != "" {
			if showPassword {
				fmt.Printf("  Password: %s\n", entry.GetPassword())
			} else {
				fmt.Printf("  Password: %s\n", strings.Repeat("*", len(entry.GetPassword())))
			}
		}

		if entry.Email != "" {
			fmt.Printf("  Email:    %s\n", entry.Email)
		}

		if len(entry.Tags) > 0 {
			fmt.Printf("  Tags:     %s\n", strings.Join(entry.Tags, ", "))
		}

		if entry.Notes != "" && showNotes {
			printIndentedLineOfChar("·")
			wrapped := wordwrap.WrapString(entry.Notes, 76)
			indented := textutils.Indent(wrapped, "  ")
			fmt.Println(strings.Trim(indented, "\r\n"))
		}

		printLineOfChar("─")
	}

	// Footer.
	fmt.Printf("  %d/%d Entries\n", len(list), poolSize)
	printLineOfChar("─")
}
