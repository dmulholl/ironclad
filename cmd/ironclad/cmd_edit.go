package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo/v4"
)

var editCmdHelptext = `
Usage: ironclad edit <entry>

  Edits an existing database entry.

  You can specify the fields to edit using the flags listed below. If no flags
  are specified you will be prompted to edit each field in turn. Enter 'y' to
  edit a field or 'n' (or simply hit return) to leave the field unchanged.

  The entry can be specified by its ID or by any unique set of case-insensitive
  substrings of its title.

Arguments:
  <entry>                   Entry to edit by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -e, --email               Edit the entry's email address.
  -h, --help                Print this command's help text and exit.
  -n, --notes               Edit the entry's notes.
  -p, --password            Edit the entry's password.
      --tags                Edit the entry's tags.
  -t, --title               Edit the entry's title.
      --url                 Edit the entry's url.
  -u, --username            Edit the entry's username.
`

func registerEditCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("edit")
	cmdParser.Helptext = editCmdHelptext
	cmdParser.Callback = editCmdCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("title t")
	cmdParser.NewFlag("url l")
	cmdParser.NewFlag("username u")
	cmdParser.NewFlag("password p")
	cmdParser.NewFlag("notes n")
	cmdParser.NewFlag("tags s")
	cmdParser.NewFlag("email e")
}

func editCmdCallback(cmdName string, cmdParser *argo.ArgParser) error {
	if len(cmdParser.Args) == 0 {
		return fmt.Errorf("missing entry argument")
	}

	filename, err := getDatabaseFilename(cmdParser)
	if err != nil {
		return err
	}

	masterpass, db, err := loadDB(filename)
	if err != nil {
		return err
	}

	matchingEntries := db.Active().FilterByAll(cmdParser.Args...)
	if len(matchingEntries) == 0 {
		return fmt.Errorf("no matching entry")
	}

	if len(matchingEntries) > 1 {
		fmt.Println("The query string matches multiple entries:")
		printCompactList(matchingEntries, len(db.Active()), filepath.Base(filename))
		return nil
	}

	entry := matchingEntries[0]

	// Default to editing all fields if no flags are present.
	allFields := false
	if !cmdParser.Found("title") && !cmdParser.Found("url") &&
		!cmdParser.Found("username") && !cmdParser.Found("password") &&
		!cmdParser.Found("tags") && !cmdParser.Found("notes") &&
		!cmdParser.Found("email") {
		allFields = true
	}

	printHeading("Editing Entry: "+entry.Title, filepath.Base(filename))

	if cmdParser.Found("title") || (allFields && editField("title")) {
		fmt.Println("  Old title: " + entry.Title)
		entry.Title, err = input("  New title: ")
		if err != nil {
			return err
		}
		printLineOfChar("·")
	}

	if cmdParser.Found("url") || (allFields && editField("url")) {
		fmt.Println("  Old URL: " + entry.Url)
		entry.Url, err = input("  New URL: ")
		if err != nil {
			return err
		}
		printLineOfChar("·")
	}

	if cmdParser.Found("username") || (allFields && editField("username")) {
		fmt.Println("  Old username: " + entry.Username)
		entry.Username, err = input("  New username: ")
		if err != nil {
			return err
		}
		printLineOfChar("·")
	}

	if cmdParser.Found("password") || (allFields && editField("password")) {
		fmt.Println("  Old password: " + entry.GetPassword())
		newPassword, err := input("  New password: ")
		if err != nil {
			return err
		}
		entry.SetPassword(newPassword)
		printLineOfChar("·")
	}

	if cmdParser.Found("email") || (allFields && editField("email")) {
		fmt.Println("  Old email: " + entry.Email)
		entry.Email, err = input("  New email: ")
		if err != nil {
			return err
		}
		printLineOfChar("·")
	}

	if cmdParser.Found("tags") || (allFields && editField("tags")) {
		fmt.Println("  Old tags: " + strings.Join(entry.Tags, ", "))
		tagstring, err := input("  New tags: ")
		if err != nil {
			return err
		}

		tagslice := strings.Split(tagstring, ",")
		entry.Tags = []string{}

		for _, tag := range tagslice {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				entry.Tags = append(entry.Tags, tag)
			}
		}

		printLineOfChar("·")
	}

	if cmdParser.Found("notes") || (allFields && editField("notes")) {
		entry.Notes, err = inputViaEditor(entry.Notes)
		if err != nil {
			return err
		}
	}

	if err := saveDB(filename, masterpass, db); err != nil {
		return err
	}

	fmt.Println("  Entry updated.")
	printLineOfChar("─")

	return nil
}

// Ask the user whether they want to edit the specified field.
func editField(fieldname string) bool {
	answer, err := input("  Edit " + fieldname + "? (y/n) ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}

	printLineOfChar("·")
	return strings.ToLower(answer) == "y"
}
