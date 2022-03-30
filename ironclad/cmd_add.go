package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/ironclad/irondb"
	"github.com/dmulholl/janus/v2"
)

var addHelp = fmt.Sprintf(`
Usage: %s add

  Add a new entry to a database. You will be prompted to supply values for
  the entry's fields - press return to leave an unwanted field blank.

  This command has an alias, 'new', which works identically.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
      --no-editor           Do not launch an external editor to add notes.
`, filepath.Base(os.Args[0]))

func registerAddCmd(parser *janus.ArgParser) {
	cmd := parser.NewCmd("add new", addHelp, addCallback)
	cmd.NewString("file f")
	cmd.NewFlag("no-editor")
}

func addCallback(parser *janus.ArgParser) {
	filename, masterpass, db := loadDB(parser)

	// Create a new Entry object to add to the database.
	entry := irondb.NewEntry()

	// Fetch user input.
	printHeading("Add Entry", filepath.Base(filename))
	entry.Title = input("  Title:      ")
	entry.Url = input("  URL:        ")
	entry.Username = input("  Username:   ")
	entry.Email = input("  Email:      ")

	// Get or autogenerate a password.
	printLineOfChar("─")
	prompt := "  Enter a password or press return to automatically generate one"
	prompt += ":\n\u001B[90m  >>\u001B[0m  "
	entrypass := input(prompt)
	entrypass = strings.TrimSpace(entrypass)
	if entrypass == "" {
		entrypass = genPassword(DefaultLength, true, true, true, true, false)
	}
	entry.SetPassword(entrypass)

	// Split tags on commas.
	printLineOfChar("─")
	prompt = "  Enter a comma-separated list of tags for this entry"
	prompt += ":\n\u001B[90m  >>\u001B[0m  "
	tagstring := input(prompt)
	for _, tag := range strings.Split(tagstring, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			entry.Tags = append(entry.Tags, tag)
		}
	}

	// Add a note?
	printLineOfChar("─")
	answer := input("  Add a note to this entry? (y/n): ")
	if strings.ToLower(answer) == "y" {
		if parser.GetFlag("no-editor") {
			printLineOfChar("─")
			entry.Notes = inputViaStdin()
		} else {
			entry.Notes = inputViaEditor("add-note", "")
		}
	} else {
		entry.Notes = ""
	}
	printLineOfChar("─")

	db.Add(entry)
	saveDB(filename, masterpass, db)
}
