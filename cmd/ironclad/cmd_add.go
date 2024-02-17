package main

import (
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/irondb"
)

var addCmdHelptext = `
Usage: ironclad add

  Adds a new entry to a database. You will be prompted to supply values for
  the entry's fields - press return to leave an unwanted field blank.

  This command has an alias, 'new', which works identically.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
  -h, --help                Print this command's help text and exit.
      --no-editor           Do not launch an external editor to add notes.
`

func registerAddCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("add new")
	cmdParser.Helptext = addCmdHelptext
	cmdParser.Callback = addCmdCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("no-editor")
}

func addCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	filename, masterpass, db := loadDB(cmdParser)

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
	prompt := "  Enter a password or press return to automatically generate one:\n\u001B[90m  >>\u001B[0m  "
	password := input(prompt)
	password = strings.TrimSpace(password)
	if password == "" {
		generatedPassword, err := genPassword(DefaultLength, true, true, true, true, false)
		if err != nil {
			return err
		}
		password = generatedPassword
	}
	entry.SetPassword(password)

	// Split tags on commas.
	printLineOfChar("─")
	prompt = "  Enter a comma-separated list of tags for this entry:\n\u001B[90m  >>\u001B[0m  "
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
		if cmdParser.Found("no-editor") {
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

	return nil
}
