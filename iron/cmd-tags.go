package main


import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/dmulholland/clio/go/clio"
    "github.com/dmulholland/ironclad/irondb"
    "sort"
)


// Help text for the 'tags' command.
var tagsHelptext = fmt.Sprintf(`
Usage: %s tags [FLAGS] [OPTIONS]

  List the tags in a database.

Options:
  -f, --file <str>          Database file.

Flags:
      --help                Print this command's help text and exit.
`, filepath.Base(os.Args[0]))


// Callback for the 'tags' command.
func tagsCallback(parser *clio.ArgParser) {

    var filename, password string
    var found bool

    // Determine the filename to use.
    filename = parser.GetStringOption("file")
    if filename == "" {
        if filename, found = fetchLastFilename(); !found {
            filename = input("Filename: ")
        }
    }

    // Determine the password to use.
    password = parser.GetStringOption("db-password")
    if password == "" {
        if password, found = fetchLastPassword(); !found {
            password = input("Password: ")
        }
    }

    // Load the database file.
    db, _, err := irondb.Load(password, filename)
    if err != nil {
        exit("Error:", err)
    }

    // Assemble a map of tags.
    tagmap := db.Tags()

    // Extract a sorted slice of tag strings.
    tags := make([]string, 0)
    for tag := range tagmap {
        tags = append(tags, tag)
    }
    sort.Strings(tags)

    // Print the tag list.
    line("-")
    fmt.Println("  Tags")
    line("-")
    for _, tag := range tags {
        fmt.Printf("  %s [%d]\n", tag, len(tagmap[tag]))
    }
    line("-")

    // Cache the password and filename.
    cacheLastPassword(password)
    cacheLastFilename(filename)
}
