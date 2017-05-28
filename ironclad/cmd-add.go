package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "strings"
    "os"
    "path/filepath"
)


import (
    "github.com/dmulholland/ironclad/irondb"
)



// Help text for the 'add' command.
var addHelp = fmt.Sprintf(`
Usage: %s add [FLAGS] [OPTIONS]

  Add a new entry to a database. You will be promped to supply values for
  the entry's fields - press return to leave an unwanted field blank.

Options:
  -f, --file <str>          Database file. Defaults to the most recent file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'add' command.
func addCallback(parser *clio.ArgParser) {

    // Load the database.
    filename, password, db := loadDB(parser)

    // Create a new Entry object to add to the database.
    entry := irondb.NewEntry()

    // Print header.
    line("─")
    fmt.Println("  Add Entry")
    line("─")

    // Fetch user input.
    entry.Title    = input("  Title:      ")
    entry.Url      = input("  URL:        ")
    entry.Username = input("  Username:   ")
    entry.Email    = input("  Email:      ")
    entry.Password = input("  Password:   ")

    // Split tags on commas.
    line("─")
    tagstring := input(
        "  Enter a comma-separated list of tags for this entry:\n> ")
    for _, tag := range strings.Split(tagstring, ",") {
        tag = strings.TrimSpace(tag)
        if tag != "" {
            entry.Tags = append(entry.Tags, tag)
        }
    }

    // Do we need to launch a text editor to add notes?
    line("─")
    answer := input("  Add a note to this entry? (y/n): ")
    if len(answer) > 0 && strings.ToLower(answer)[0] == 'y' {
        entry.Notes = inputViaEditor("add-note", "")
    } else {
        entry.Notes = ""
    }

    // Add the new entry to the database.
    db.Add(entry)

    // Save the updated database to disk.
    saveDB(filename, password, db)

    // Footer.
    line("─")
}
