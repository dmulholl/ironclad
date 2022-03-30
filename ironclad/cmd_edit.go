package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmulholl/argo"
)

var editHelp = fmt.Sprintf(`
Usage: %s edit <entry>

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
      --no-editor           Do not launch an external editor to edit notes.
  -n, --notes               Edit the entry's notes.
  -p, --password            Edit the entry's password.
      --tags                Edit the entry's tags.
  -t, --title               Edit the entry's title.
      --url                 Edit the entry's url.
  -u, --username            Edit the entry's username.
`, filepath.Base(os.Args[0]))

func registerEditCmd(parser *argo.ArgParser) {
	cmdParser := parser.NewCommand("edit")
	cmdParser.Helptext = editHelp
	cmdParser.Callback = editCallback
	cmdParser.NewStringOption("file f", "")
	cmdParser.NewFlag("title t")
	cmdParser.NewFlag("url l")
	cmdParser.NewFlag("username u")
	cmdParser.NewFlag("password p")
	cmdParser.NewFlag("notes n")
	cmdParser.NewFlag("tags s")
	cmdParser.NewFlag("email e")
	cmdParser.NewFlag("no-editor")
}

func editCallback(cmdName string, cmdParser *argo.ArgParser) {
	if !cmdParser.HasArgs() {
		exit("missing entry argument")
	}
	filename, masterpass, db := loadDB(cmdParser)

	// Search for an entry corresponding to the supplied arguments.
	list := db.Active().FilterByAll(cmdParser.Args...)
	if len(list) == 0 {
		exit("no matching entry")
	} else if len(list) > 1 {
		println("Error: the query string matches multiple entries.")
		printCompact(list, len(db.Active()), filepath.Base(filename))
		os.Exit(1)
	}
	entry := list[0]

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
		entry.Title = input("  New title: ")
		printLineOfChar("·")
	}

	if cmdParser.Found("url") || (allFields && editField("url")) {
		fmt.Println("  Old URL: " + entry.Url)
		entry.Url = input("  New URL: ")
		printLineOfChar("·")
	}

	if cmdParser.Found("username") || (allFields && editField("username")) {
		fmt.Println("  Old username: " + entry.Username)
		entry.Username = input("  New username: ")
		printLineOfChar("·")
	}

	if cmdParser.Found("password") || (allFields && editField("password")) {
		fmt.Println("  Old password: " + entry.GetPassword())
		entry.SetPassword(input("  New password: "))
		printLineOfChar("·")
	}

	if cmdParser.Found("email") || (allFields && editField("email")) {
		fmt.Println("  Old email: " + entry.Email)
		entry.Email = input("  New email: ")
		printLineOfChar("·")
	}

	if cmdParser.Found("tags") || (allFields && editField("tags")) {
		fmt.Println("  Old tags: " + strings.Join(entry.Tags, ", "))
		tagstring := input("  New tags: ")
		tagslice := strings.Split(tagstring, ",")
		entry.Tags = make([]string, 0)
		for _, tag := range tagslice {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				entry.Tags = append(entry.Tags, tag)
			}
		}
		printLineOfChar("·")
	}

	if cmdParser.Found("notes") || (allFields && editField("notes")) {
		if cmdParser.Found("no-editor") {
			oldnotes := strings.Trim(entry.Notes, "\r\n")
			if oldnotes != "" {
				fmt.Println(oldnotes)
				printLineOfChar("·")
			}
			entry.Notes = inputViaStdin()
			printLineOfChar("·")
		} else {
			entry.Notes = inputViaEditor("edit-note", entry.Notes)
		}
	}

	saveDB(filename, masterpass, db)
	fmt.Println("  Entry updated.")
	printLineOfChar("─")
}

// Ask the user whether they want to edit the specified field.
func editField(field string) bool {
	answer := input("  Edit " + field + "? (y/n) ")
	printLineOfChar("·")
	return strings.ToLower(answer) == "y"
}
