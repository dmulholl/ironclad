package main


import "github.com/dmulholland/clio/go/clio"


import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)


var editHelp = fmt.Sprintf(`
Usage: %s edit [FLAGS] [OPTIONS] ARGUMENTS

  Edit an existing database entry.

  You can specify the fields to edit using the flags listed below. If no
  flags are specified you will be prompted to edit each field in turn.
  Enter 'y' to edit a field or 'n' (or simply hit return) to leave the
  field unchanged.

Arguments:
  <entry>                   Entry to edit by ID or title.

Options:
  -f, --file <str>          Database file. Defaults to the last used file.

Flags:
  -e, --email               Edit the entry's email address.
      --help                Print this command's help text and exit.
      --no-editor           Do not launch an external editor to edit notes.
  -n, --notes               Edit the entry's notes.
  -p, --password            Edit the entry's password.
      --tags                Edit the entry's tags.
  -t, --title               Edit the entry's title.
      --url                 Edit the entry's url.
  -u, --username            Edit the entry's username.
`, filepath.Base(os.Args[0]))


func editCallback(parser *clio.ArgParser) {

    // Make sure an argument has been specified.
    if !parser.HasArgs() {
        exit("missing entry argument")
    }

    // Load the database.
    filename, password, db := loadDB(parser)

    // Search for an entry corresponding to the specified argument.
    list := db.Active().FilterProgressive(parser.GetArgs()[0])
    if len(list) == 0 {
        exit("no matching entry")
    } else if len(list) > 1 {
        exit("query matches multiple entries")
    }
    entry := list[0]

    // Default to editing all fields if no flags are present.
    allFields := false
    if !parser.GetFlag("title") && !parser.GetFlag("url") &&
        !parser.GetFlag("username") && !parser.GetFlag("password") &&
        !parser.GetFlag("tags") && !parser.GetFlag("notes") &&
        !parser.GetFlag("email") {
        allFields = true
    }

    // Header.
    line("─")
    fmt.Println("  Editing Entry: " + entry.Title)
    line("─")

    if parser.GetFlag("title") || (allFields && editField("title")) {
        fmt.Println("  TITLE")
        fmt.Println("  Old value: " + entry.Title)
        entry.Title = input("  New value: ")
        line("·")
    }

    if parser.GetFlag("url") || (allFields && editField("url")) {
        fmt.Println("  URL")
        fmt.Println("  Old value: " + entry.Url)
        entry.Url = input("  New value: ")
        line("·")
    }

    if parser.GetFlag("username") || (allFields && editField("username")) {
        fmt.Println("  USERNAME")
        fmt.Println("  Old value: " + entry.Username)
        entry.Username = input("  New value: ")
        line("·")
    }

    if parser.GetFlag("password") || (allFields && editField("password")) {
        fmt.Println("  PASSWORD")
        fmt.Println("  Old value: " + entry.GetPassword())
        entry.SetPassword(input("  New value: "))
        line("·")
    }

    if parser.GetFlag("email") || (allFields && editField("email")) {
        fmt.Println("  EMAIL")
        fmt.Println("  Old value: " + entry.Email)
        entry.Email = input("  New value: ")
        line("·")
    }

    if parser.GetFlag("tags") || (allFields && editField("tags")) {
        fmt.Println("  TAGS")
        fmt.Println("  Old value: " + strings.Join(entry.Tags, ", "))
        tagstring := input("  New value: ")
        tagslice := strings.Split(tagstring, ",")
        entry.Tags = make([]string, 0)
        for _, tag := range tagslice {
            tag = strings.TrimSpace(tag)
            if tag != "" {
                entry.Tags = append(entry.Tags, tag)
            }
        }
        line("·")
    }

    if parser.GetFlag("notes") || (allFields && editField("notes")) {
        if parser.GetFlag("no-editor") {
            oldnotes := strings.Trim(entry.Notes, "\r\n")
            if oldnotes != "" {
                fmt.Println(oldnotes)
                line("·")
            }
            entry.Notes = inputViaStdin()
            line("·")
        } else {
            entry.Notes = inputViaEditor("edit-note", entry.Notes)
        }
    }

    // Save the updated database to disk.
    saveDB(filename, password, db)

    // Footer.
    fmt.Println("  Entry updated.")
    line("─")
}


// Ask the user whether they want to edit the specified field.
func editField(field string) bool {
    answer := input("  Edit " + field + "? (y/n) ")
    line("·")
    return strings.ToLower(answer) == "y"
}
