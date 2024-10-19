package main

import (
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo/v4"
	"github.com/dmulholl/ironclad/internal/database"
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
`

func registerAddCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("add new")
	cmdParser.Helptext = addCmdHelptext
	cmdParser.Callback = addCmdCallback
	cmdParser.NewStringOption("file f", "")
}

func addCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	masterpass, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	// Create a new Entry object to add to the database.
	entry := database.NewEntry()

	// Fetch user input.
	printHeading("Add Entry", filepath.Base(filename))

	entry.Title, err = input("  Title:      ")
	if err != nil {
		return err
	}

	entry.Url, err = input("  URL:        ")
	if err != nil {
		return err
	}

	entry.Username, err = input("  Username:   ")
	if err != nil {
		return err
	}

	entry.Email, err = input("  Email:      ")
	if err != nil {
		return err
	}

	// Get or autogenerate a password.
	printLineOfChar("─")
	prompt := "  Enter a password or press return to automatically generate one:\n\u001B[90m  >>\u001B[0m  "

	password, err := input(prompt)
	if err != nil {
		return err
	}

	password = strings.TrimSpace(password)
	if password == "" {
		generatedPassword, err := genPassword(DefaultLength, true, true, false, true, false)
		if err != nil {
			return err
		}
		password = generatedPassword
	}

	entry.SetPassword(password)

	// Split tags on commas.
	printLineOfChar("─")
	prompt = "  Enter a comma-separated list of tags for this entry:\n\u001B[90m  >>\u001B[0m  "

	tagstring, err := input(prompt)
	if err != nil {
		return err
	}

	for _, tag := range strings.Split(tagstring, ",") {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			entry.Tags = append(entry.Tags, tag)
		}
	}

	// Add a note?
	printLineOfChar("─")

	answer, err := input("  Add a note to this entry? (y/n): ")
	if err != nil {
		return err
	}

	if strings.ToLower(answer) == "y" {
		entry.Notes, err = inputViaEditor("")
		if err != nil {
			return err
		}
	}

	printLineOfChar("─")

	db.Add(entry)

	if err := saveDB(filename, masterpass, db); err != nil {
		return err
	}

	return nil
}
