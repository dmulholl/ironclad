package main


import "github.com/tonnerre/golang-text"
import "github.com/mitchellh/go-wordwrap"


import (
    "strings"
)


import (
    "github.com/dmulholland/ironclad/irondb"
)


// Print a list of entries in compact format.
func printCompact(list irondb.EntryList, dbsize int) {

    // Bail if we have no entries to display.
    if len(list) == 0 {
        printLineOfChar("─")
        println("  No Entries")
        printLineOfChar("─")
        return
    }

    // Header.
    printLineOfChar("─")
    print("  ID")
    printGrey("  ·  ")
    print("TITLE\n")
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
func printVerbose(list irondb.EntryList, dbsize int, title string) {

    // Bail if we have no entries to display.
    if len(list) == 0 {
        printLineOfChar("─")
        println("  No Entries")
        printLineOfChar("─")
        return
    }

    // Header.
    printLineOfChar("─")
    println("  " + title)
    printLineOfChar("─")

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
            indented := text.Indent(wrapped, "  ")
            println(strings.Trim(indented, "\r\n"))
        }

        printLineOfChar("─")
    }

    // Footer.
    print("  %d/%d Entries\n", len(list), dbsize)
    printLineOfChar("─")
}
