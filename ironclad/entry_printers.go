package main

import (
	"strings"

	"github.com/dmulholl/ironclad/irondb"
	"github.com/mitchellh/go-wordwrap"
)

// Print a list of entries in compact format.
func printCompact(list irondb.EntryList, dbsize int, filename string) {

	// Bail if we have no entries to display.
	if len(list) == 0 {
		printHeading("No Entries", filename)
		return
	}

	// Header.
	printLineOfChar("─")
	print("  ID")
	printGrey("  ·  ")
	print("TITLE")
	numSpaces := terminalWidth() - len(filename) - 16
	for i := 0; i < numSpaces; i += 1 {
		print(" ")
	}
	printlnGrey(filename)
	printLineOfChar("─")

	// Print the entry listing.
	for _, entry := range list {
		print("%4d", entry.Id)
		printGrey("  ·  ")
		print("%s\n", entry.Title)
	}

	// Footer.
	printLineOfChar("─")
	print("  %d/%d Entries\n", len(list), dbsize)
	printLineOfChar("─")
}

// Print a list of entries in verbose format.
func printVerbose(list irondb.EntryList, dbsize int, title, filename string) {

	// Bail if we have no entries to display.
	if len(list) == 0 {
		printHeading("No Entries", filename)
		return
	}

	// Header.
	printHeading(title, filename)

	// Print the entry listing.
	for _, entry := range list {
		print("  ID:       %d\n", entry.Id)
		print("  Title:    %s\n", entry.Title)

		if entry.Url != "" {
			print("  URL:      %s\n", entry.Url)
		}

		if entry.Username != "" {
			print("  Username: %s\n", entry.Username)
		}

		if entry.GetPassword() != "" {
			print("  Password: %s\n", entry.GetPassword())
		}

		if entry.Email != "" {
			print("  Email:    %s\n", entry.Email)
		}

		if len(entry.Tags) > 0 {
			print("  Tags:     %s\n", strings.Join(entry.Tags, ", "))
		}

		if entry.Notes != "" {
			printIndentedLineOfChar("·")
			wrapped := wordwrap.WrapString(entry.Notes, 76)
			indented := indent(wrapped, "  ")
			println(strings.Trim(indented, "\r\n"))
		}

		printLineOfChar("─")
	}

	// Footer.
	print("  %d/%d Entries\n", len(list), dbsize)
	printLineOfChar("─")
}
